---
title: "Upgrade"
date:
draft: false
weight: 33
---

# Overview

Upgrading to a new version of IVYO is typically as simple as following the various installation
guides defined within the IVYO documentation:

- [IVYO Kustomize Install]({{< relref "./kustomize.md" >}})
- [IVYO Helm Install]({{< relref "./helm.md" >}})

However, when upgrading to or from certain versions of IVYO, extra steps may be required in order
to ensure a clean and successful upgrade.

This section provides detailed instructions for upgrading IVYO 5.x using Kustomize or Helm, along with information for upgrading from IVYO v4 to IVYO v5.

{{% notice info %}}
Depending on version updates, upgrading IVYO may automatically rollout changes to managed Ivory clusters. This could result in downtime--we cannot guarantee no interruption of service, though IVYO attempts graceful incremental rollouts of affected pods, with the goal of zero downtime.
{{% /notice %}}

## Upgrading IVYO 5.x

- [IVYO Kustomize Upgrade]({{< relref "./kustomize.md" >}})
- [IVYO Helm Upgrade]({{< relref "./helm.md" >}})

## Upgrading from IVYO v4 to IVYO v5

- [V4 to V5 Upgrade Methods]({{< relref "./v4tov5" >}})
