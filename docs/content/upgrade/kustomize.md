---
title: "Upgrading IVYO v5 Using Kustomize"
date:
draft: false
weight: 50
---

## Upgrading to v5.4.0 from v5.3.x

Apply the new version of the Kubernetes installer:

```bash
kubectl apply --server-side -k kustomize/install/default
```

IVYO versions from 5.1.x through 5.3.x include a ivyo-upgrade deployment, which
is no longer needed after upgrading to v5.4.x. Delete the deployment:

```bash
kubectl delete deployment ivyo-upgrade
```

## Upgrading from IVYO v5.0.0 Using Kustomize

Starting with IVYO v5.0.1, both the Deployment and ServiceAccount created when installing IVYO via
the installers in the
[Ivory Operator examples repository](https://github.com/ivorysql/ivory-operator-examples)
have been renamed from `ivory-operator` to `ivyo`.  As a result of this change, if using
Kustomize to install IVYO and upgrading from IVYO v5.0.0, the following step must be completed prior
to upgrading.  This will ensure multiple versions of IVYO are not installed and running concurrently
within your Kubernetes environment.

Prior to upgrading IVYO, first manually delete the IVYO v5.0.0 `ivory-operator` Deployment and
ServiceAccount:

```bash
kubectl -n ivory-operator delete deployment,serviceaccount ivory-operator
```

Then, once both the Deployment and ServiceAccount have been deleted, proceed with upgrading IVYO
by applying the new version of the Kustomize installer:

```bash
kubectl apply --server-side -k kustomize/install/default
```

## Upgrading from IVYO v5.0.2 and Below

As a result of changes to pgBackRest dedicated repository host deployments in IVYO v5.0.3
(please see the [IVYO v5.0.3 release notes]({{< relref "../releases/5.0.3.md" >}}) for more details),
reconciliation of a pgBackRest dedicated repository host might become stuck with the following
error (as shown in the IVYO logs) following an upgrade from IVYO versions v5.0.0 through v5.0.2:

```bash
StatefulSet.apps \"hippo-repo-host\" is invalid: spec: Forbidden: updates to statefulset spec for fields other than 'replicas', 'template', 'updateStrategy' and 'minReadySeconds' are forbidden
```

If this is the case, proceed with deleting the pgBackRest dedicated repository host StatefulSet,
and IVYO will then proceed with recreating and reconciling the dedicated repository host normally:

```bash
kubectl delete sts hippo-repo-host
```

Additionally, please be sure to update and apply all PostgresCluster custom resources in accordance
with any applicable spec changes described in the
[IVYO v5.0.3 release notes]({{< relref "../releases/5.0.3.md" >}}).

## Upgrading from IVYO v5.0.5 and Below

Starting in IVYO v5.1, new pgBackRest features available in version 2.38 are used
that impact both the `highgo-ivory` and `highgo-pgbackrest` images. For any
clusters created before v5.0.6, you will need to update these image values
BEFORE upgrading to IVYO {{< param operatorVersion >}}. These changes will need
to be made in one of two places, depending on your desired configuration.

If you are setting the image values on your `PostgresCluster` manifest,
you would update the images value as shown (updating the `image` values as
appropriate for your environment):

```yaml
apiVersion: ivory-operator.crunchydata.com/v1beta1
kind: PostgresCluster
metadata:
  name: hippo
spec:
  image: {{< param imageCrunchyPostgres >}}
  postgresVersion: {{< param postgresVersion >}}
...
  backups:
    pgbackrest:
      image: {{< param imageCrunchyPGBackrest >}}
...
```

After updating these values, you will apply these changes to your PostgresCluster
custom resources. After these changes are completed and the new images are in place,
you may update IVYO to {{< param operatorVersion >}}.

Relatedly, if you are instead using the `RELATED_IMAGE` environment variables to
set the image values, you would instead check and update these as needed before
redeploying IVYO.

For Kustomize installations, these can be found in the `manager` directory and
`manager.yaml` file. Here you will note various key/value pairs, these will need
to be updated before deploying IVYO {{< param operatorVersion >}}. Besides updating the
`RELATED_IMAGE_PGBACKREST` value, you will also need to update the relevant
Ivory image for your environment. For example, if you are using IvorySQL 14,
you would update the value for `RELATED_IMAGE_POSTGRES_14`. If instead you are
using the PostGIS 3.1 enabled IvorySQL 13 image, you would update the value
for `RELATED_IMAGE_POSTGRES_13_GIS_3.1`.

For Helm deployments, you would instead need to similarly update your `values.yaml`
file, found in the `install` directory. There you will note a `relatedImages`
section, followed by similar values as mentioned above. Again, be sure to update
`pgbackrest` as well as the appropriate `postgres` value for your clusters.

Once there values have been properly verified, you may deploy IVYO {{< param operatorVersion >}}.