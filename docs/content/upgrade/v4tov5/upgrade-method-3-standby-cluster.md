---
title: "Upgrade Method #3: Standby Cluster"
date:
draft: false
weight: 30
---

{{% notice info %}}
Before attempting to upgrade from v4.x to v5, please familiarize yourself with the [prerequisites]({{< relref "upgrade/v4tov5/_index.md" >}}) applicable for all v4.x to v5 upgrade methods.
{{% /notice %}}

This upgrade method allows you to migrate from IVYO v4 to IVYO v5 by creating a new IVYO v5 Ivory cluster in a "standby" mode, allowing it to mirror the IVYO v4 cluster and continue to receive data updates in real time. This has the advantage of being able to fully inspect your IVYO v5 Ivory cluster while leaving your IVYO v4 cluster up and running, thus minimizing downtime when you cut over. The tradeoff is that you will temporarily use more resources while this migration is occurring.

This method only works if your IVYO v4 cluster uses S3 or an S3-compatible storage system, or GCS. For more information on standby clusters, please refer to the [tutorial]({{< relref "tutorial/disaster-recovery.md" >}}#standby-cluster).

### Step 1: Migrate to IVYO v5

Create a [`PostgresCluster`]({{< relref "references/crd.md" >}}) custom resource. This migration method does not carry over any specific configurations or customizations from IVYO v4: you will need to create the specific `PostgresCluster` configuration that you need.

To complete the upgrade process, your `PostgresCluster` custom resource **MUST** include the following:

1\. Configure your pgBackRest to use an object storage system such as S3/GCS. You will need to configure your `spec.backups.pgbackrest.repos.volume` to point to the backup storage system. For instance, if AWS S3 storage is being utilized, the repo would be defined similar to the following:

```
spec:
  backups:
    pgbackrest:
      repos:
        - name: repo1
          s3:
            bucket: hippo
            endpoint: s3.amazonaws.com
            region: us-east-1
```

Any required secrets or desired custom pgBackRest configuration should be created and configured as described in the [backup tutorial]({{< relref "tutorial/backups.md" >}}).

You will also need to ensure that the “pgbackrest-repo-path” configured for the repository matches the path used by the IVYO v4 cluster. The default repository path follows the pattern `/backrestrepo/<clusterName>-backrest-shared-repo`. Note that the path name here is different than migrating from a PVC-based repository.

Using the `hippo` Ivory cluster as an example, you would set the following in the `spec.backups.pgbackrest.global` section:

```
spec:
  backups:
    pgbackrest:
      global:
        repo1-path: /backrestrepo/hippo-backrest-shared-repo
```

2\. A `spec.standby` cluster configuration within the spec that is populated according to the name of pgBackRest repo configured in the spec. For example:

```
spec:
  standby:
    enabled: true
    repoName: repo1
```

3\. If you customized other Ivory parameters, you will need to ensure they match in the IVYO v5 cluster. For more information, please review the tutorial on [customizing a Ivory cluster]({{< relref "tutorial/customize-cluster.md" >}}).

4\. Once the `PostgresCluster` spec is populated according to these guidelines, you can create the `PostgresCluster` custom resource.  For example, if the `PostgresCluster` you're creating is a modified version of the [`postgres` example](https://github.com/Highgo/ivory-operator-examples/tree/main/kustomize/postgres) in the [IVYO examples repo](https://github.com/Highgo/ivory-operator-examples), you can run the following command:

```
kubectl apply -k examples/postgrescluster
```

5\. Once the standby cluster is up and running and you are satisfied with your set up, you can promote it.

First, you will need to shut down your IVYO v4 cluster. You can do so with the following command, e.g.:

```
ivyo update cluster hippo --shutdown
```

You can then update your IVYO v5 cluster spec to promote your standby cluster:

```
spec:
  standby:
    enabled: false
```

Note: When the v5 cluster is running in non-standby mode, you will not be able to restart the v4 cluster, as that data is now being managed by the v5 cluster.

Once the v5 cluster is up and running, you will need to run the following SQL commands as a IvorySQL superuser. For example, you can login as the `postgres` user, or exec into the Pod and use `psql`:

```sql
-- add the managed replication user
CREATE ROLE _crunchyrepl WITH LOGIN REPLICATION;

-- allow for the replication user to execute the functions required as part of "rewinding"
GRANT EXECUTE ON function pg_catalog.pg_ls_dir(text, boolean, boolean) TO _crunchyrepl;
GRANT EXECUTE ON function pg_catalog.pg_stat_file(text, boolean) TO _crunchyrepl;
GRANT EXECUTE ON function pg_catalog.pg_read_binary_file(text) TO _crunchyrepl;
GRANT EXECUTE ON function pg_catalog.pg_read_binary_file(text, bigint, bigint, boolean) TO _crunchyrepl;
```

The above step will be automated in an upcoming release.

Your upgrade is now complete! Once you verify that the IVYO v5 cluster is running and you have recorded the user credentials from the v4 cluster, you can remove the old cluster:

```
ivyo delete cluster hippo
```

For more information on how to use IVYO v5, we recommend reading through the [IVYO v5 tutorial]({{< relref "tutorial/_index.md" >}}).
