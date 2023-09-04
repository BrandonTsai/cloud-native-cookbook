---
title: "#4 OpenShift Monitoring - Prometheus"
author: Brandon Tsai
---

Cluster Monitoring Stack Overview
==================================

OpenShift Container Platform includes a pre-configured, pre-installed, and self-updating monitoring stack that is based on the Prometheus open source project and its wider eco-system. It provides monitoring of cluster components and includes a set of alerts to immediately notify the cluster administrator about any occurring problems and a set of Grafana dashboards. The cluster monitoring stack is only supported for monitoring OpenShift Container Platform clusters.


Monitoring Stack components 
----------------------------

The monitoring stack includes these components:

**Cluster Monitoring Operator**

The OpenShift Container Platform Cluster Monitoring Operator (CMO) is the central component of the stack. It controls the deployed monitoring components and resources and ensures that they are always up to date. For more details, Please refer [here](https://github.com/openshift/cluster-monitoring-operator)

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



Monitor Application on OpenShift
================


### Step 1. Enabling monitoring of your own services


Monitoring your own services is a Technology Preview feature only. Technology Preview features are not supported with Red Hat production service level agreements (SLAs) and might not be functionally complete.
You can enable monitoring your own services by setting the techPreviewUserWorkload/enabled flag in the cluster monitoring ConfigMap.

```
# cluster-monitoring-config.yaml 
apiVersion: v1
kind: ConfigMap
metadata:
  name: cluster-monitoring-config
  namespace: openshift-monitoring
data:
  config.yaml: |
    techPreviewUserWorkload:
      enabled: true
```

Check that the prometheus-user-workload pods were created:

```
$ oc -n openshift-user-workload-monitoring get pod
NAME                                   READY   STATUS    RESTARTS   AGE
prometheus-operator-5857b6db84-mc5zx   1/1     Running   0          13m
prometheus-user-workload-0             5/5     Running   1          13m
prometheus-user-workload-1             5/5     Running   1          13m

```


### Step 2. Deploying a sample service


```
apiVersion: v1
kind: Namespace
metadata:
  name: brandon
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: prometheus-example-app
  name: prometheus-example-app
  namespace: brandon
spec:
  replicas: 1
  selector:
    matchLabels:
      app: prometheus-example-app
  template:
    metadata:
      labels:
        app: prometheus-example-app
    spec:
      containers:
      - image: quay.io/brancz/prometheus-example-app:v0.2.0
        imagePullPolicy: IfNotPresent
        name: prometheus-example-app
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: prometheus-example-app
  name: prometheus-example-app
  namespace: brandon
spec:
  ports:
  - port: 8080
    protocol: TCP
    targetPort: 8080
    name: web
  selector:
    app: prometheus-example-app
  type: ClusterIP

```

### Step 3. Setting up metrics collection

```
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  labels:
    k8s-app: prometheus-example-monitor
  name: prometheus-example-monitor
  namespace: brandon
spec:
  endpoints:
  - interval: 30s
    port: web
    scheme: http
  selector:
    matchLabels:
      app: prometheus-example-app
```


### Step 4. Create Alerts Rule

```
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: example-alert
  namespace: ns1
spec:
  groups:
  - name: example
    rules:
    - alert: VersionAlert
      expr: version{job="prometheus-example-app"} == 0	  
```  


Monitor external node
==============


Create Endpoint

```
apiVersion: v1
kind: Endpoints
metadata:
  name: quay-metrics-ports
  labels:
    k8s-app: quay-metrics-ports
    prometheus: kube-prometheus
  namespace: brandon
subsets:
- addresses:
  - ip: {{ quay_ip }}
    targetRef:
      kind: Node
      name: quay-uat
  ports:
  - name: quay-metrics
    port: 9092
    protocol: TCP
  - name: node-metrics
    port: 9100
    protocol: TCP
 ```
 
 Create "ExternalName" Service
 
 ```
kind: "Service"
apiVersion: "v1"
metadata:
  name: quay-metrics-ports
  labels:
    k8s-app: quay-metrics-ports
    prometheus: kube-prometheus
  namespace: openshift-monitoring
spec:
  type: ExternalName
  externalName: {{ quay_ip }}
  ports:
    - name: quay-metrics
      port: 9092
      protocol: TCP
      targetPort: 9092
    - name: node-metrics
      port: 9100
      protocol: TCP
      targetPort: 9100
selector: {}
```

Create ServiceMonitor


```
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  generation: 1
  name: quay
  labels:
    k8s-app: quay-metrics
    prometheus: kube-prometheus
  namespace: openshift-monitoring
spec:
  jobLabel: k8s-app
  selector:
    matchLabels:
      k8s-app: quay-{{ cluster_region }}-{{ cluster_type }}-metrics-ports
  namespaceSelector:
    matchNames:
      - openshift-monitoring
  endpoints:
  - port: quay-metrics
    interval: 30s
    path: /metrics
    scheme: http
    honorLabels: true
  - port: node-metrics
    interval: 30s
    path: /metrics
    scheme: http
    honorLabels: true
```


Create alerts

```
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  labels:
    prometheus: k8s
    role: alert-rules
  name: quay-rules
  namespace: openshift-monitoring
spec:
  groups:
  - name: quay.rules
    rules:
    - expr: quay_repository_count
      record: quay_repository_count
    - expr: |
        100 - 100 * (node_filesystem_avail_bytes{mountpoint="/opt"} / node_filesystem_size_bytes{mountpoint="/opt"})
      record: quay_node_data_disk_usage
    - expr: |
        100 - 100 * (node_filesystem_avail_bytes{device="rootfs",mountpoint="/",job="quay"} / node_filesystem_size_bytes{device="rootfs",mountpoint="/",job="quay"})
      record: quay_node_root_disk_usage
  - name: quay-alert.rules
    rules:
    - alert: QuayDataDiskRunningFull
      annotations:
        message: 'Quay data volume usage on target {{ $labels.instance }} at {{ $value }}%'
      expr: |
        100 - 100 * (node_filesystem_avail_bytes{mountpoint="/opt"} / node_filesystem_size_bytes{mountpoint="/opt"}) > 85
      for: 15m
      labels:
        severity: warning
    - alert: QuayRootDiskRunningFull
      annotations:
        message: 'Quay root volume usage on target {{ $labels.instance }} at {{ $value }}%'
      expr: |
        100 - 100 * (node_filesystem_avail_bytes{device="rootfs",mountpoint="/",job="quay"} / node_filesystem_size_bytes{device="rootfs",mountpoint="/",job="quay"}) > 85
      for: 15m
      labels:
        severity: warning
    - alert: HighSystemLoad
      annotations:
        message: 'Quay has high system load on target {{ $labels.instance }} at {{ $value }} for past 30 minutes'
      expr: |
        node_load15{service=~"quay"} > 2
      for: 30m
      labels:
        severity: warning
```




Notes
-----

- You can not change the default Grafana dashboard and alert rules in openshift-monitoring namespaces, it is managed by Cluster Monitoring Operator. These setting will be overwrite by Cluster Monitoring Operator. 

- The Grafana instance shipped within OpenShift Container Platform Monitoring is read-only and displays only infrastructure-related dashboards.

To solve this issue we could use the Grafana Operator from OperatorHub or Grafana Helm Chart.
