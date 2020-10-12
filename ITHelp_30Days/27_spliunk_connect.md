Monitor Kubernetes apps with Splunk - Part 1 : Introduction and Splunk Setting
============================================================================

Some companies use Splunk as the logging platform to store and to aggregate the logs for all their environments.
This post explains how to integrate Splunk with Kubernetes using the [Splunk-connect-for-kubernetes](https://github.com/splunk/splunk-connect-for-kubernetes) helm charts.


Architecture
-------------

Splunk-connect-for-kubernetes contains 3 componenets:

| Component | Usage |
| --------- | ------ |
| logging   | To collect container logs. |
| metric    | To collect metrics, such as cpu/memory usage. |
| objects   | To collect kubernetes resource status by calling the Kubernetes API. |


Prerequisites
--------------

- Allow connection from kubernetes to Splunk on HEC port 8088
- Access to Dockerhub to pull images or import these images to your private registry.
- Helm 3
- splunk-connect-for-k8s v1.3.0

Splunk Setting
--------------

### Create Splunk App (Optional)

Create new app "kubernetes" or using exist app

> Note: The app must be 'visible' to be able to add indexes to it.

### Create Splunk Indexes


Create following indexes for the default indexes of http event collector

| index name | type | app |
| ---------- | ---- | --- |
| k8s_\<ENV>_logging | Events | kubernetes |
| k8s_\<ENV>_metrics | Metrics | kubernetes |
| k8s_\<ENV>_objects | Events | kubernetes |


Note:
- The size of each index should be tailored for the amount of data expected and not left as default
- \<ENV> is the environment name for the kubernetes cluster (such as UAT, PROD) if you using different kubernetes cluster for different environment.


### Push App setting to indexer cluster


**Step 1.** Copy app folder in splunk master instance

For app: "kubernetes", copy that app folder from /opt/splunk/etc/apps to /opt/splunk/etc/master-apps/


**Step 2.** On Splunk UI,

- Click Settings > Indexer Clustering.
- Click Edit > Configuration Bundle Actions.
- Click Validate and Check Restart to check the bundle is valid
- Click Push if the "Validate and Check Restart" result is fine.

**Step 3.** Check the splunk data

Check that any indexer replication issue has resolved and that Splunk is showing green. This will cause the Indexer servers to restart on initial push

### Create Splunk HTTP Event Collector Token
Navigate to Settings > Data Inputs > HTTP Event Collector


> **Notice:** Do NOT enable indexer acknowledgement when creating following tokens

We need to create 3 HEC token for logging, metrics and object

| HEC Token name     | App Context | Select Allowed Indexes | Default index |
| ------------------ | ----------- | ---------------------- | ------------- |
| k8s-\<ENV>-logging | kubernetes  | k8s_\<ENV>\_logging    | k8s_\<ENV>\_logging |
| k8s-\<ENV>-metrics | kubernetes  | k8s_\<ENV>\_metrics    | k8s_\<ENV>\_metrics |
| k8s-\<ENV>-objects | kubernetes  | k8s_\<ENV>\_objects    | k8s_\<ENV>\_objects |


Monitor Kubernetes apps with Splunk - Part 2 : logging
=======================================================

Prerequisites
--------------

Download the lastest Helm package from [Splunk-connect-for-kubernetes](https://github.com/splunk/splunk-connect-for-kubernetes).



Update values.yaml for logging
-----------------------------

The minimal value example:

```YAML
splunk:
  hec:
    host: < splunk_host >
    port: 8088
    token: < splunk_hec_logging_token >
    indexName: < splunk_logging_indexname >
```

**Optional:** Customize buffer setting

``` YAML
buffer:
  "@type": memory
  total_limit_size: 8000m
  chunk_limit_size: 8m
  chunk_limit_records: 10000000000
  flush_at_shutdown: true
  flush_interval: 3s
  flush_thread_count: 20
  flush_thread_interval: 0.1
  flush_thread_burst_interval: 0.01
  overflow_action: block
  retry_forever: true
  retry_wait: 30
  compress: gzip
```

**Optional:** Customize filter setting
```YAML
customFilters:
  SetNamespaceFilter:
    tag: "**"
    type: grep
    body: |
        <exclude>
                  key namespace
                  pattern /(kube-system)/
                </exclude>
                <exclude>
                  key sourcetype
                  pattern /(fluentd:monitor-agent|kube:container:calico-node)/
                </exclude>
```

Deplopy to Kubernetes Cluster
--------------------------------

You can deploy to kubernetes cluster via helm directly.
Or generate kubernetes yaml files via helm template command and then deploy via kubectl.

```bash
helm template --name-template=k8s --namespace splunk-connect --output-dir ${output_folder} splunk-kubernetes-logging/

kubectl apply -f ${output_folder}/splunk-kubernetes-logging/templates/
```


Verify on Splunk
----------------

```
index="k8s_<ENV>_logging"
```
Monitor Kubernetes apps with Splunk - Part 3 : Metrics
========================================================


Prerequisites
--------------

Download the lastest Helm package from [Splunk-connect-for-kubernetes](https://github.com/splunk/splunk-connect-for-kubernetes).



Update values.yaml for metrics
-----------------------------

The minimal value example:

```YAML
splunk:
  hec:
    host: < splunk_host >
    port: 8088
    token: < splunk_hec_metrics_token >
    indexName: < splunk_metrics_indexname >
```

**Optional:** Customize filter setting

Please refer [metrics-information](https://github.com/splunk/fluent-plugin-kubernetes-metrics/blob/develop/metrics-information.md) for all supported metrics

It is recommand to customize the fluentd setting to collect minimal metrics that are required for monitoing.


```YAML
customFilters:
  SetContainerFilter:
    tag: kube.container.**
    type: grep
    body: |
        <regexp>
                key metric_name
                pattern /(cpu.usage_rate|cpu.limit|memory.usage|memory.limit)/
              </regexp>
  SetPodFilter:
    tag: kube.pod.**
    type: grep
    body: |
        <regexp>
                key metric_name
                pattern /(network.rx_bytes|network.tx_bytes|network.rx_errors|network.tx_errors|cpu.load.average.10s|cpu.usage_rate|cpu.limit|memory.usage|memory.limit|memory.available_bytes|volume.available_bytes|volume.used_bytes)/
              </regexp>
  SetNamespaceFilter:
    tag: kube.namespace.**
    type: grep
    body: |
        <regexp>
                key metric_name
                pattern /(usage|limit)/
              </regexp>
  SetNodeFilter:
    tag: kube.node.**
    type: grep
    body: |
        <regexp>
                key metric_name
                pattern /(network.rx_bytes|network.tx_bytes|network.rx_errors|network.tx_errors|cpu.usage_rate|memory.usage|memory.capacity|memory.available_bytes)/
              </regexp>

```

Deplopy to Kubernetes Cluster
--------------------------------

You can deploy to kubernetes cluster via helm directly.
Or generate kubernetes yaml files via helm template command and then deploy via kubectl.

```bash
helm template --name-template=k8s --namespace splunk-connect --output-dir ${output_folder} splunk-kubernetes-metrics/

kubectl apply -f ${output_folder}/splunk-kubernetes-metrics/templates/
```

Verify on Splunk
----------------

Following splunk search can be used to check the supported dimensions of a metric:

```
| mcatalog values(_dims) WHERE index="*_metrics" AND metric_name="kube.pod.cpu.load.average.10s"
```

Monitor Kubernetes apps with Splunk - Part 4 : Objects
=========================================================

Splunk collects the resource information by calling the Kubernetes API. It help user/operator to set up splunk alerts when pod is in Error status.

Prerequisites
--------------

Download the lastest Helm package from [Splunk-connect-for-kubernetes](https://github.com/splunk/splunk-connect-for-kubernetes).



Update values.yaml for objects
-------------------------------

The minimal value example:

```YAML
splunk:
  hec:
    host: < splunk_host >
    port: 8088
    token: < splunk_hec_objects_token >
    indexName: < splunk_objects_indexname >
```

Deplopy to Kubernetes Cluster
--------------------------------

You can deploy to kubernetes cluster via helm directly.
Or generate kubernetes yaml files via helm template command and then deploy via kubectl.

```bash
helm template --name-template=k8s --namespace splunk-connect --output-dir ${output_folder} splunk-kubernetes-objects/

kubectl apply -f ${output_folder}/splunk-kubernetes-objects/templates/
```


Verify on Splunk
----------------

Query the number of running pods in splunk-connect namespace

```
index="k8s_<ENV>_objects" metadata.namespace="splunk-connect" status.phase="Running" | stats distinct_count(metadata.uid)
```
