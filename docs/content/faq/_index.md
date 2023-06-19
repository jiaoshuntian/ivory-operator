---
title: "FAQ"
date:
draft: false
weight: 105

aliases:
 - /contributing
---

## Project FAQ

### What is The IVYO Project?

The IVYO Project is the open source project associated with the development of [IVYO](https://github.com/ivorysql/ivory-operator), the [Ivory Operator](https://github.com/ivorysql/ivory-operator) for Kubernetes from [Highgo](https://www.crunchydata.com).

IVYO is a [Kubernetes Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/), providing a declarative solution for managing your IvorySQL clusters.  Within a few moments, you can have a Ivory cluster complete with high availability, disaster recovery, and monitoring, all over secure TLS communications.

IVYO is the upstream project from which [Highgo IvorySQL for Kubernetes](https://www.crunchydata.com/products/highgo-ivorysql-for-kubernetes/) is derived. You can find more information on Highgo IvorySQL for Kubernetes [here](https://www.crunchydata.com/products/highgo-ivorysql-for-kubernetes/).

### What’s the difference between IVYO and Highgo IvorySQL for Kubernetes?

IVYO is the Ivory Operator from Highgo. It developed pursuant to the IVYO Project and is designed to be a frequently released, fast-moving project where all new development happens.

[Highgo IvorySQL for Kubernetes](https://www.crunchydata.com/products/highgo-ivorysql-for-kubernetes/) is produced by taking selected releases of IVYO, combining them with Highgo Certified IvorySQL and IvorySQL containers certified by Highgo, maintained for commercial support, and made available to customers as the Highgo IvorySQL for Kubernetes offering.

### Where can I find support for IVYO?

The community can help answer questions about IVYO via the [IVYO mailing list](https://groups.google.com/a/crunchydata.com/forum/#!forum/ivory-operator/join).

Information regarding support for IVYO is available in the [Support]({{< relref "support/_index.md" >}}) section of the IVYO documentation, which you can find [here]({{< relref "support/_index.md" >}}).

For additional information regarding commercial support and Highgo IvorySQL for Kubernetes, you can [contact Highgo](https://www.crunchydata.com/contact/).

### Under which open source license is IVYO source code available?

The IVYO source code is available under the [Apache License 2.0](https://github.com/ivorysql/ivory-operator/blob/master/LICENSE.md).

### Where are the release tags for IVYO v5?

With IVYO v5, we've made some changes to our overall process. Instead of providing quarterly release
tags as we did with IVYO v4, we're focused on ongoing active development in the v5 primary
development branch (`master`, which will become `main`).  Consistent with our practices in v4,
previews of stable releases with the release tags are made available in the
[Highgo Developer Portal](https://www.crunchydata.com/developers).

These changes allow for more rapid feature development and releases in the upstream IVYO project,
while providing
[Highgo Ivory for Kubernetes](https://www.crunchydata.com/products/highgo-ivorysql-for-kubernetes/)
users with stable releases for production use.

To the extent you have constraints specific to your use, please feel free to reach out on
[info@crunchydata.com](mailto:info@crunchydata.com) to discuss how we can address those
specifically.

### How can I get involved with the IVYO Project?

IVYO is developed by the IVYO Project. The IVYO Project that welcomes community engagement and contribution.

The IVYO source code and community issue trackers are hosted at [GitHub](https://github.com/ivorysql/ivory-operator).

For community questions and support, please sign up for the [IVYO mailing list](https://groups.google.com/a/crunchydata.com/forum/#!forum/ivory-operator/join).

For information regarding contribution, please review the contributor guide [here](https://github.com/ivorysql/ivory-operator/blob/master/CONTRIBUTING.md).

Please register for the [Highgo Developer Portal mailing list](https://www.crunchydata.com/developers/newsletter) to receive updates regarding Highgo IvorySQL for Kubernetes releases and the [Highgo newsletter](https://www.crunchydata.com/newsletter/) for general updates from Highgo.

### Where do I report a IVYO bug?

The IVYO Project uses GitHub for its [issue tracking](https://github.com/ivorysql/ivory-operator/issues/new/choose). You can file your issue [here](https://github.com/ivorysql/ivory-operator/issues/new/choose).

### How often is IVYO released?

The IVYO team currently plans to release new builds approximately every few weeks. The IVYO team will flag certain builds as “stable” at their discretion. Note that the term “stable” does not imply fitness for production usage or any kind of warranty whatsoever.
