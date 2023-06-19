---
title: "Private Registries"
date:
draft: false
weight: 200
---

IVYO, the open source Ivory Operator, can use containers that are stored in private registries.
There are a variety of techniques that are used to load containers from private registries,
including [image pull secrets](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/).
This guide will demonstrate how to install IVYO and deploy a Ivory cluster using the
[Highgo Customer Portal](https://access.crunchydata.com/) registry as an example.

## Create an Image Pull Secret

The Kubernetes documentation provides several methods for creating
[image pull secrets](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/).
You can choose the method that is most appropriate for your installation. You will need to create
image pull secrets in the namespace that IVYO is deployed and in each namespace where you plan to
deploy Ivory clusters.

For example, to create an image pull secret for accessing the Highgo Customer Portal image
registry in the `ivory-operator` namespace, you can execute the following commands:

```shell
kubectl create ns ivory-operator

kubectl create secret docker-registry highgo-regcred -n ivory-operator \
  --docker-server=registry.crunchydata.com \
  --docker-username=<YOUR USERNAME> \
  --docker-email=<YOUR EMAIL> \
  --docker-password=<YOUR PASSWORD>
```

This creates an image pull secret named `highgo-regcred` in the `ivory-operator` namespace.

## Install IVYO from a Private Registry

To [install IVYO]({{< relref "installation/_index.md" >}}) from a private registry, you will need to
set an image pull secret on the installation manifest.

For example, to set up an image pull secret using the [Kustomize install method]({{< relref "installation/_index.md" >}})
to install IVYO from the [Highgo Customer Portal](https://access.crunchydata.com/), you can set
the following in the `kustomize/install/default/kustomization.yaml` manifest:

```yaml
images:
- name: ivory-operator
  newName: {{< param operatorRepositoryPrivate >}}
  newTag: {{< param postgresOperatorTag >}}

patchesJson6902:
  - target:
      group: apps
      version: v1
      kind: Deployment
      name: ivyo
    patch: |-
      - op: remove
        path: /spec/selector/matchLabels/app.kubernetes.io~1name
      - op: remove
        path: /spec/selector/matchLabels/app.kubernetes.io~1version
      - op: add
        path: /spec/template/spec/imagePullSecrets
        value:
          - name: highgo-regcred
```

If you are using a version of `kubectl` prior to `v1.21.0`, you will have to create an explicit
patch file named `install-ops.yaml`:

```yaml
- op: remove
  path: /spec/selector/matchLabels/app.kubernetes.io~1name
- op: remove
  path: /spec/selector/matchLabels/app.kubernetes.io~1version
- op: add
  path: /spec/template/spec/imagePullSecrets
  value:
    - name: highgo-regcred
```

and modify the manifest to be the following:

```yaml
images:
- name: ivory-operator
  newName: {{< param operatorRepositoryPrivate >}}
  newTag: {{< param postgresOperatorTag >}}

patchesJson6902:
  - target:
      group: apps
      version: v1
      kind: Deployment
      name: ivyo
    path: install-ops.yaml
```

You can then install IVYO from the private registry using the standard installation procedure, e.g.:

```shell
kubectl apply --server-side -k kustomize/install/default
```

## Deploy a Ivory cluster from a Private Registry

To deploy a Ivory cluster using images from a private registry, you will need to set the value of
`spec.imagePullSecrets` on a `PostgresCluster` custom resource.

For example, to deploy a Ivory cluster using images from the [Highgo Customer Portal](https://access.crunchydata.com/)
with an image pull secret in the `ivory-operator` namespace, you can use the following manifest:

```yaml
apiVersion: ivory-operator.crunchydata.com/v1beta1
kind: PostgresCluster
metadata:
  name: hippo
spec:
  imagePullSecrets:
    - name: highgo-regcred
  image: {{< param imageCrunchyPostgresPrivate >}}
  postgresVersion: {{< param postgresVersion >}}
  instances:
    - name: instance1
      dataVolumeClaimSpec:
        accessModes:
        - "ReadWriteOnce"
        resources:
          requests:
            storage: 1Gi
  backups:
    pgbackrest:
      image: {{< param imageCrunchyPGBackrestPrivate >}}
      repos:
      - name: repo1
        volume:
          volumeClaimSpec:
            accessModes:
            - "ReadWriteOnce"
            resources:
              requests:
                storage: 1Gi
```
