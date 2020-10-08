

Cluster Monitoring Stack Overview
------------------

OpenShift Container Platform includes a pre-configured, pre-installed, and self-updating monitoring stack that is based on the Prometheus open source project and its wider eco-system. It provides monitoring of cluster components and includes a set of alerts to immediately notify the cluster administrator about any occurring problems and a set of Grafana dashboards. The cluster monitoring stack is only supported for monitoring OpenShift Container Platform clusters.


Monitoring Stack components 
----------------------------

The monitoring stack includes these components:


**[Cluster Monitoring Operator](https://github.com/openshift/cluster-monitoring-operator)**:

The OpenShift Container Platform Cluster Monitoring Operator (CMO) is the central component of the stack. It controls the deployed monitoring components and resources and ensures that they are always up to date.

**Prometheus Operator**

The Prometheus Operator (PO) creates, configures, and manages Prometheus and Alertmanager instances. It also automatically generates monitoring target configurations based on familiar Kubernetes label queries.

**Prometheus**

The Prometheus is the systems and service monitoring system, around which the monitoring stack is based.

**Prometheus Adapter**

The Prometheus Adapter exposes cluster resource metrics API for horizontal pod autoscaling. Resource metrics are CPU and memory utilization.

**Alertmanager**

The Alertmanager service handles alerts sent by Prometheus.

**kube-state-metrics**

The kube-state-metrics exporter agent converts Kubernetes objects to metrics that Prometheus can use.

**openshift-state-metrics**

The openshift-state-metrics exporter expands upon kube-state-metrics by adding metrics for OpenShift Container Platform-specific resources.

**node-exporter**

node-exporter is an agent deployed on every node to collect metrics about it.

**Thanos Querier**

The Thanos Querier enables aggregating and, optionally, deduplicating cluster and user workload metrics under a single, multi-tenant interface.

**Grafana**

The Grafana analytics platform provides dashboards for analyzing and visualizing the metrics. The Grafana instance that is provided with the monitoring stack, along with its dashboards, is read-only.





You can check these components in "openshift-monitoring" namespaces.


```
$ oc project openshift-monitoring
$ oc get deploy
NAME                          READY   UP-TO-DATE   AVAILABLE   AGE
cluster-monitoring-operator   1/1     1            1           110d
grafana                       1/1     1            1           110d
kube-state-metrics            1/1     1            1           110d
openshift-state-metrics       1/1     1            1           110d
prometheus-adapter            2/2     2            2           110d
prometheus-operator           1/1     1            1           110d
thanos-querier                2/2     2            2           110d
```

Monitoring Walkthrough
-------------------------

With the integration of Prometheus and Grafana, OpenShift allow system admin from the OpenShift Web UI directly.

1. Check alerts in directly in Monitoring > Metrics page

2. Do Metrics Query directly in Monitoring > Metrics page


3. Check Dashboard from in Monitoring > Dashboards page





Besides, you still can go to the original Prometheus and Grafana Web UI and login to OpenShift Account.

1. 

2. 










Notes
-----

- You can not change the default Grafana dashboard and alert rules in openshift-monitoring namespaces, it is managed by Cluster Monitoring Operator. These setting will be overwrite by Cluster Monitoring Operator. 

- The Grafana instance shipped within OpenShift Container Platform Monitoring is read-only and displays only infrastructure-related dashboards.

