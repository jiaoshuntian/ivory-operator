/*
 Copyright 2021 - 2023 Highgo Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package pgbouncer

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/ivorysql/ivory-operator/internal/testing/require"
	"github.com/ivorysql/ivory-operator/pkg/apis/ivory-operator.ivorysql.org/v1beta1"
)

func TestPrettyYAML(t *testing.T) {
	b, err := yaml.Marshal(iniValueSet{
		"x": "y",
		"z": "",
	}.String())
	assert.NilError(t, err)
	assert.Assert(t, strings.HasPrefix(string(b), `|`),
		"expected literal block scalar, got:\n%s", b)
}

func TestAuthFileContents(t *testing.T) {
	t.Parallel()

	password := `very"random`
	data := authFileContents(password)
	assert.Equal(t, string(data), `"_highgopgbouncer" "very""random"`+"\n")
}

func TestClusterINI(t *testing.T) {
	t.Parallel()

	cluster := new(v1beta1.IvoryCluster)
	cluster.Default()

	cluster.Name = "foo-baz"
	*cluster.Spec.Port = 9999

	cluster.Spec.Proxy = new(v1beta1.IvoryProxySpec)
	cluster.Spec.Proxy.PGBouncer = new(v1beta1.PGBouncerPodSpec)
	cluster.Spec.Proxy.PGBouncer.Port = new(int32)
	*cluster.Spec.Proxy.PGBouncer.Port = 8888

	t.Run("Default", func(t *testing.T) {
		assert.Equal(t, clusterINI(cluster), strings.Trim(`
# Generated by ivory-operator. DO NOT EDIT.
# Your changes will not be saved.

[pgbouncer]
%include /etc/pgbouncer/pgbouncer.ini

[pgbouncer]
auth_file = /etc/pgbouncer/~ivory-operator/users.txt
auth_query = SELECT username, password from pgbouncer.get_auth($1)
auth_user = _highgopgbouncer
client_tls_ca_file = /etc/pgbouncer/~ivory-operator/frontend-ca.crt
client_tls_cert_file = /etc/pgbouncer/~ivory-operator/frontend-tls.crt
client_tls_key_file = /etc/pgbouncer/~ivory-operator/frontend-tls.key
client_tls_sslmode = require
conffile = /etc/pgbouncer/~ivory-operator.ini
ignore_startup_parameters = extra_float_digits
listen_addr = *
listen_port = 8888
server_tls_ca_file = /etc/pgbouncer/~ivory-operator/backend-ca.crt
server_tls_sslmode = verify-full
unix_socket_dir =

[databases]
* = host=foo-baz-primary port=9999
		`, "\t\n")+"\n")
	})

	t.Run("CustomSettings", func(t *testing.T) {
		cluster.Spec.Proxy.PGBouncer.Config.Global = map[string]string{
			"ignore_startup_parameters": "custom",
			"verbose":                   "whomp",
		}
		cluster.Spec.Proxy.PGBouncer.Config.Databases = map[string]string{
			"appdb": "conn=str",
		}
		cluster.Spec.Proxy.PGBouncer.Config.Users = map[string]string{
			"app": "mode=rad",
		}

		assert.Equal(t, clusterINI(cluster), strings.Trim(`
# Generated by ivory-operator. DO NOT EDIT.
# Your changes will not be saved.

[pgbouncer]
%include /etc/pgbouncer/pgbouncer.ini

[pgbouncer]
auth_file = /etc/pgbouncer/~ivory-operator/users.txt
auth_query = SELECT username, password from pgbouncer.get_auth($1)
auth_user = _highgopgbouncer
client_tls_ca_file = /etc/pgbouncer/~ivory-operator/frontend-ca.crt
client_tls_cert_file = /etc/pgbouncer/~ivory-operator/frontend-tls.crt
client_tls_key_file = /etc/pgbouncer/~ivory-operator/frontend-tls.key
client_tls_sslmode = require
conffile = /etc/pgbouncer/~ivory-operator.ini
ignore_startup_parameters = custom
listen_addr = *
listen_port = 8888
server_tls_ca_file = /etc/pgbouncer/~ivory-operator/backend-ca.crt
server_tls_sslmode = verify-full
unix_socket_dir =
verbose = whomp

[databases]
appdb = conn=str

[users]
app = mode=rad
		`, "\t\n")+"\n")

		// The "conffile" setting cannot be changed.
		cluster.Spec.Proxy.PGBouncer.Config.Global["conffile"] = "too-far"
		assert.Assert(t, !strings.Contains(clusterINI(cluster), "too-far"))
	})
}

func TestPodConfigFiles(t *testing.T) {
	t.Parallel()

	config := v1beta1.PGBouncerConfiguration{}
	configmap := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "some-cm"}}
	secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "some-shh"}}

	t.Run("Default", func(t *testing.T) {
		projections := podConfigFiles(config, configmap, secret)
		assert.Assert(t, marshalMatches(projections, `
- configMap:
    items:
    - key: pgbouncer-empty
      path: pgbouncer.ini
    name: some-cm
- configMap:
    items:
    - key: pgbouncer.ini
      path: ~ivory-operator.ini
    name: some-cm
- secret:
    items:
    - key: pgbouncer-users.txt
      path: ~ivory-operator/users.txt
    name: some-shh
		`))
	})

	t.Run("CustomFiles", func(t *testing.T) {
		config.Files = []corev1.VolumeProjection{
			{Secret: &corev1.SecretProjection{
				LocalObjectReference: corev1.LocalObjectReference{Name: "my-thing"},
			}},
			{Secret: &corev1.SecretProjection{
				LocalObjectReference: corev1.LocalObjectReference{Name: "also"},
				Items: []corev1.KeyToPath{
					{Key: "specific", Path: "files"},
				},
			}},
		}

		projections := podConfigFiles(config, configmap, secret)
		assert.Assert(t, marshalMatches(projections, `
- configMap:
    items:
    - key: pgbouncer-empty
      path: pgbouncer.ini
    name: some-cm
- secret:
    name: my-thing
- secret:
    items:
    - key: specific
      path: files
    name: also
- configMap:
    items:
    - key: pgbouncer.ini
      path: ~ivory-operator.ini
    name: some-cm
- secret:
    items:
    - key: pgbouncer-users.txt
      path: ~ivory-operator/users.txt
    name: some-shh
		`))
	})
}

func TestReloadCommand(t *testing.T) {
	shellcheck := require.ShellCheck(t)
	command := reloadCommand("some-name")

	// Expect a bash command with an inline script.
	assert.DeepEqual(t, command[:3], []string{"bash", "-ceu", "--"})
	assert.Assert(t, len(command) > 3)

	// Write out that inline script.
	dir := t.TempDir()
	file := filepath.Join(dir, "script.bash")
	assert.NilError(t, os.WriteFile(file, []byte(command[3]), 0o600))

	// Expect shellcheck to be happy.
	cmd := exec.Command(shellcheck, "--enable=all", file)
	output, err := cmd.CombinedOutput()
	assert.NilError(t, err, "%q\n%s", cmd.Args, output)
}
