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

Note: The app must be 'visible' to be able to add indexes to it.

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
