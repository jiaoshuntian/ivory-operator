// Copyright 2021 - 2023 Crunchy Data Solutions, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ivyupgrade

import (
	"github.com/ivorysql/ivory-operator/pkg/apis/ivory-operator.ivorysql.org/v1beta1"
)

const (
	// ConditionIvyUpgradeProgressing is the type used in a condition to indicate that
	// an Ivory major upgrade is in progress.
	ConditionIvyUpgradeProgressing = "Progressing"

	// ConditionIvyUpgradeSucceeded is the type used in a condition to indicate the
	// status of a Ivory major upgrade.
	ConditionIvyUpgradeSucceeded = "Succeeded"

	labelPrefix           = "ivory-operator.ivorysql.org/"
	LabelIvyUpgrade       = labelPrefix + "ivyupgrade"
	LabelCluster          = labelPrefix + "cluster"
	LabelRole             = labelPrefix + "role"
	LabelVersion          = labelPrefix + "version"
	LabelPatroni          = labelPrefix + "patroni"
	LabelPGBackRestBackup = labelPrefix + "pgbackrest-backup"
	LabelInstance         = labelPrefix + "instance"

	ReplicaCreate     = "replica-create"
	ContainerDatabase = "database"

	pgUpgrade  = "ivyupgrade"
	removeData = "removedata"
)

func commonLabels(role string, upgrade *v1beta1.IvyUpgrade) map[string]string {
	return map[string]string{
		LabelIvyUpgrade: upgrade.Name,
		LabelCluster:    upgrade.Spec.IvoryClusterName,
		LabelRole:       role,
	}
}
