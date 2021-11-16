Monitor Application on OpenShift
================


Enabling monitoring of your own services
-----------

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

check that the prometheus-user-workload pods were created:

```
$ oc -n openshift-user-workload-monitoring get pod
NAME                                   READY   STATUS    RESTARTS   AGE
prometheus-operator-5857b6db84-mc5zx   1/1     Running   0          13m
prometheus-user-workload-0             5/5     Running   1          13m
prometheus-user-workload-1             5/5     Running   1          13m

```


Deploying a sample service
-----------

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
	  
	  
Setting up metrics collection
-----------

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

Make sure you can Query metric from Web UI
-------------




Create Alerts Rule
----------

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
========


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

create ServiceMonitor


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

