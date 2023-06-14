---
title: "Upgrade Method #1: Data Volumes"
date:
draft: false
weight: 10
---

{{% notice info %}}
Before attempting to upgrade from v4.x to v5, please familiarize yourself with the [prerequisites]({{< relref "upgrade/v4tov5/_index.md" >}}) applicable for all v4.x to v5 upgrade methods.
{{% /notice %}}

This upgrade method allows you to migrate from IVYO v4 to IVYO v5 using the existing data volumes that were created in IVYO v4. Note that this is an "in place" migration method: this will immediately move your Ivory clusters from being managed by IVYO v4 and IVYO v5. If you wish to have some failsafes in place, please use one of the other migration methods. Please also note that you will need to perform the cluster upgrade in the same namespace as the original cluster in order for your v5 cluster to access the existing PVCs.

### Step 1: Prepare the IVYO v4 Cluster for Migration

You will need to set up your IVYO v4 Ivory cluster so that it can be migrated to a IVYO v5 cluster. The following describes how to set up a IVYO v4 cluster for using this migration method.

1. Scale down any existing replicas within the cluster.  This will ensure that the primary PVC does not change again prior to the upgrade.

You can get a list of replicas using the `ivyo scaledown --query` command, e.g.:
```
ivyo scaledown hippo --query
```

If there are any replicas, you will see something similar to:

```
Cluster: hippo
REPLICA                 STATUS          NODE ...
hippo                   running         node01 ...
```

Scaledown any replicas that are running in this cluser, e.g.:

```
ivyo scaledown hippo --target=hippo
```

2\. Once all replicas are removed and only the primary remains, proceed with deleting the cluster while retaining the data and backups. You can do this `--keep-data` and  `--keep-backups` flags:

**You MUST run this command with the `--keep-data` and `--keep-backups` flag otherwise you risk deleting ALL of your data.**

```
ivyo delete cluster hippo --keep-data --keep-backups
```

3\. The PVC for the primary Ivory instance and the pgBackRest repository should still remain. You can verify this with the command below:

```
kubectl get pvc --selector=pg-cluster=hippo
```

This should yield something similar to:

```
NAME              STATUS   VOLUME ...
hippo-jgut        Bound    pvc-a0b89bdb- ...
hippo-pgbr-repo   Bound    pvc-25501671- …
```

A third PVC used to store write-ahead logs (WAL) may also be present if external WAL volumes were enabled for the cluster.

### Step 2: Migrate to IVYO v5

With the IVYO v4 cluster's volumes prepared for the move to IVYO v5, you can now create a [`PostgresCluster`]({{< relref "references/crd.md" >}}) custom resource using these volumes. This migration method does not carry over any specific configurations or customizations from IVYO v4: you will need to create the specific `PostgresCluster` configuration that you need.

{{% notice warning %}}

Additional steps are required to set proper file permissions when using certain storage options,
such as NFS and HostPath storage, due to a known issue with how fsGroups are applied. When
migrating from IVYO v4, this will require the user to manually set the group value of the pgBackRest
repo directory, and all subdirectories, to `26` to match the `postgres` group used in IVYO v5.
Please see [here](https://github.com/kubernetes/examples/issues/260) for more information.

{{% /notice %}}

To complete the upgrade process, your `PostgresCluster` custom resource **MUST** include the following:

1\. A `volumes` data source that points to the IvorySQL data, IvorySQL WAL (if applicable) and pgBackRest repository PVCs identified in the `spec.dataSource.volumes` section.

For example, using the `hippo` cluster:

```
spec:
  dataSource:
    volumes:
      pgDataVolume:
        pvcName: hippo-jgut
        directory: "hippo-jgut"
      pgBackRestVolume:
        pvcName: hippo-pgbr-repo
        directory: "hippo-backrest-shared-repo"
      # Only specify external WAL PVC if enabled in IVYO v4 cluster. If enabled
      # in v4, a WAL volume must be defined for the v5 cluster as well.
      # pgWALVolume:
      #  pvcName: hippo-jgut-wal
```

Please see the [Data Migration]({{< relref "guides/data-migration.md" >}}) section of the [tutorial]({{< relref "tutorial/_index.md" >}}) for more details on how to properly populate this section of the spec when migrating from a IVYO v4 cluster.

{{% notice info %}}
Note that when migrating data volumes from v4 to v5, IVYO relabels all volumes for IVYO v5, but **will not remove existing IVYO v4 labels**. This results in PVCs that are labeled for both IVYO v4 and v5, which can lead to unintended behavior.
<br><br>
To avoid that behavior, follow the instructions in the section on [removing IVYO v4 labels]({{< ref "guides/data-migration.md#removing-ivyo-v4-labels" >}}).
{{% /notice %}}

2\. If you customized Ivory parameters, you will need to ensure they match in the IVYO v5 cluster. For more information, please review the tutorial on [customizing a Ivory cluster]({{< relref "tutorial/customize-cluster.md" >}}).

3\. Once the `PostgresCluster` spec is populated according to these guidelines, you can create the `PostgresCluster` custom resource.  For example, if the `PostgresCluster` you're creating is a modified version of the [`postgres` example](https://github.com/Highgo/ivory-operator-examples/tree/main/kustomize/postgres) in the [IVYO examples repo](https://github.com/Highgo/ivory-operator-examples), you can run the following command:

```
kubectl apply -k examples/postgrescluster
```

Your upgrade is now complete! You should now remove the `spec.dataSource.volumes` section from your `PostgresCluster`. For more information on how to use IVYO v5, we recommend reading through the [IVYO v5 tutorial]({{< relref "tutorial/_index.md" >}}).
