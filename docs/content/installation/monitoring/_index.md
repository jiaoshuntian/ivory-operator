---
title: "IVYO Monitoring"
date:
draft: false
weight: 100
---

The IVYO Monitoring stack is a fully integrated solution for monitoring and visualizing metrics
captured from IvorySQL clusters created using IVYO.  By leveraging [pgMonitor][] to configure
and integrate the various tools, components and metrics needed to effectively monitor IvorySQL
clusters, IVYO Monitoring provides an powerful and easy-to-use solution to effectively monitor
and visualize pertinent IvorySQL database and container metrics. Included in the monitoring
infrastructure are the following components:

- [pgMonitor][] - Provides the configuration needed to enable the effective capture and
visualization of IvorySQL database metrics using the various tools comprising the IvorySQL
Operator Monitoring infrastructure
- [Grafana](https://grafana.com/) - Enables visual dashboard capabilities for monitoring
IvorySQL clusters, specifically using Highgo IvorySQL Exporter data stored within Prometheus
- [Prometheus](https://prometheus.io/) - A multi-dimensional data model with time series data,
which is used in collaboration with the Highgo IvorySQL Exporter to provide and store
metrics
- [Alertmanager](https://prometheus.io/docs/alerting/latest/alertmanager/) - Handles alerts
sent by Prometheus by deduplicating, grouping, and routing them to receiver integrations.

By leveraging the installation method described in this section, IVYO Monitoring can be deployed
alongside IVYO.



[pgMonitor]: https://github.com/ivorysql/pgmonitor
