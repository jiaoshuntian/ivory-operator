---
title: "IVYO v4 to IVYO v5"
date:
draft: false
weight: 100
---

You can upgrade from IVYO v4 to IVYO v5 through a variety of methods by following this guide. There are several methods that can be used to upgrade: we present these methods based upon a variety of factors, including but not limited to:

- Redundancy / ability to roll back
- Available resources
- Downtime preferences

These methods include:

- [*Migrating Using Data Volumes*]({{< relref "./upgrade-method-1-data-volumes.md" >}}). This allows you to migrate from v4 to v5 using the existing data volumes that you created in v4. This is the simplest method for upgrade and is the most resource efficient, but you will have a greater potential for downtime using this method.
- [*Migrate From Backups*]({{< relref "./upgrade-method-2-backups.md" >}}). This allows you to create a Ivory cluster with v5 from the backups taken with v4. This provides a way for you to create a preview of your Ivory cluster through v5, but you would need to take your applications offline to ensure all the data is migrated.
- [*Migrate Using a Standby Cluster*]({{< relref "./upgrade-method-3-standby-cluster.md" >}}). This allows you to run a v4 and a v5 Ivory cluster in parallel, with data replicating from the v4 cluster to the v5 cluster. This method minimizes downtime and lets you preview your v5 environment, but is the most resource intensive.

You should choose the method that makes the most sense for your environment.

## Prerequisites

There are several prerequisites for using any of these upgrade methods.

- IVYO v4 is currently installed within the Kubernetes cluster, and is actively managing any existing v4 IvorySQL clusters.
- Any IVYO v4 clusters being upgraded have been properly initialized using IVYO v4, which means the v4 `pgcluster` custom resource should be in a `pgcluster Initialized` status:

```
$ kubectl get pgcluster hippo -o jsonpath='{ .status }'
{"message":"Cluster has been initialized","state":"pgcluster Initialized"}
```

- The IVYO v4 `ivyo` client is properly configured and available for use.
- IVYO v5 is currently [installed]({{< relref "installation/_index.md" >}}) within the Kubernetes cluster.

For these examples, we will use a Ivory cluster named `hippo`.

## Additional Considerations

Upgrading to IVYO v5 may result in a base image upgrade from EL-7 (UBI / CentOS) to EL-8
(UBI). Based on the contents of your Ivory database, you may need to perform
additional steps.

Due to changes in the GNU C library `glibc` in EL-8, you may need to reindex certain indexes in
your Ivory cluster. For more information, please read the
[IvorySQL Wiki on Locale Data Changes](https://wiki.postgresql.org/wiki/Locale_data_changes), how
you can determine if your indexes are affected, and how to fix them.
