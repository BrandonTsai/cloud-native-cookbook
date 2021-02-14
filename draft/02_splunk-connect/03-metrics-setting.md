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
