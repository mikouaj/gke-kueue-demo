# Extra assets

## Grafana Kueue dashboard

![Kueue Dashboard](./kueue-dashboard.png)

The dashboard uses Prometheus data source and metrics from the following sources:
Requirements:

* [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics) for cluster node insights
* [Kueue metrics](https://kueue.sigs.k8s.io/docs/reference/metrics/) for Kueue insights
