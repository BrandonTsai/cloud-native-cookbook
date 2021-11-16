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
