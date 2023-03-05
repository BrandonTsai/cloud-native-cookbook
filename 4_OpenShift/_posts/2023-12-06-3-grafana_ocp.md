---
title: "#6-2 Monitoring with Grafana"
author: Brandon Tsai
---


I followed the steps below to deploy a community-powered Grafana operator 3.5.0 from OperatorHub on a running OpenShift 4.5 cluster.  This allowed me to write custom queries against the built-in Prometheus to extract metrics relevant to me, and in turn I’m able to create custom dashboards to visualize those metrics.


Deploying Custom Grafana
----

From Installed Operators, select the Grafana Operator.  For the Grafana resource, press Create Instance to create a new Grafana instance.


![](go01.PNG)

In the Grafana instance YAML, make a note of the default username and password to log in, and press Create.

![](go02.PNG)


Checking all resource are running

![](go03.PNG)


Connecting Prometheus to our Custom Grafana
----

The grafana-serviceaccount service account was created alongside the Grafana instance.  We will grant it the cluster-monitoring-view cluster role.

```
$ oc project brandon
$ oc adm policy add-cluster-role-to-user cluster-monitoring-view -z grafana-serviceaccount
```


The bearer token for this service account is used to authenticate access to Prometheus in the openshift-monitoring namespace.  The following command will display this token.

```
$ oc serviceaccounts get-token grafana-serviceaccount -n brandon
```


From the Grafana Data Source resource, press Create Instance, and navigate to the YAML view.  In the below YAML, substitute ${BEARER_TOKEN} with the output of the command above, copy the YAML, and press Create.


```
apiVersion: integreatly.org/v1alpha1
kind: GrafanaDataSource
metadata:
  name: prometheus-grafanadatasource
spec:
  datasources:
    - access: proxy
      editable: true
      isDefault: true
      jsonData:
        httpHeaderName1: 'Authorization'
        timeInterval: 5s
        tlsSkipVerify: true
      name: Prometheus
      secureJsonData:
        httpHeaderValue1: 'Bearer ${BEARER_TOKEN}'
      type: prometheus
      url: 'https://thanos-querier.openshift-monitoring.svc.cluster.local:9091'
  name: prometheus-grafanadatasource.yaml
```



Check Grafana Web UI
---------------

From the my-grafana namespace, navigate to Networking -> Routes and click on the Grafana URL to display the custom Grafana user interface.  Click on ‘Sign In’ from the bottom left menu of Grafana, and log in using the default username and password configured earlier.  Now, an editable Grafana interface appears and you can view your custom Grafana dashboards or create your own.  As a note, administrators should take caution with custom dashboards to query Prometheus as this will have an impact on the performance of the monitoring stack.


![](go11.PNG)

make sure we could query metrics from prometheus

![](go12.PNG)


Customizing Grafana Dashboard
--------

You can create dashboard from Grafana directly or using the "GrafanaDashboard" resource.



```yaml
apiVersion: integreatly.org/v1alpha1
kind: GrafanaDashboard
metadata:
  labels:
    app: grafana
  name: simple-dashboard
  namespace: brandon
spec:
  json: |
    {
      "annotations": {
        "list": [
          {
            "builtIn": 1,
            "datasource": "-- Grafana --",
            "enable": true,
            "hide": true,
            "iconColor": "rgba(0, 211, 255, 1)",
            "name": "Annotations & Alerts",
            "type": "dashboard"
          }
        ]
      },
      "editable": true,
      "gnetId": null,
      "graphTooltip": 0,
      "id": 6,
      "links": [],
      "panels": [
        {
          "cacheTimeout": null,
          "colorBackground": false,
          "colorValue": false,
          "colors": [
            "#299c46",
            "rgba(237, 129, 40, 0.89)",
            "#d44a3a"
          ],
          "datasource": null,
          "format": "none",
          "gauge": {
            "maxValue": 100,
            "minValue": 0,
            "show": false,
            "thresholdLabels": false,
            "thresholdMarkers": true
          },
          "gridPos": {
            "h": 9,
            "w": 12,
            "x": 0,
            "y": 0
          },
          "id": 2,
          "interval": null,
          "links": [],
          "mappingType": 1,
          "mappingTypes": [
            {
              "name": "value to text",
              "value": 1
            },
            {
              "name": "range to text",
              "value": 2
            }
          ],
          "maxDataPoints": 100,
          "nullPointMode": "connected",
          "nullText": null,
          "options": {},
          "postfix": "",
          "postfixFontSize": "50%",
          "prefix": "",
          "prefixFontSize": "50%",
          "rangeMaps": [
            {
              "from": "null",
              "text": "N/A",
              "to": "null"
            }
          ],
          "sparkline": {
            "fillColor": "rgba(31, 118, 189, 0.18)",
            "full": false,
            "lineColor": "rgb(31, 120, 193)",
            "show": false,
            "ymax": null,
            "ymin": null
          },
          "tableColumn": "",
          "targets": [
            {
              "expr": "version{job=\"prometheus-example-app\"}",
              "refId": "A"
            }
          ],
          "thresholds": "",
          "timeFrom": null,
          "timeShift": null,
          "title": "App Version",
          "type": "singlestat",
          "valueFontSize": "80%",
          "valueMaps": [
            {
              "op": "=",
              "text": "N/A",
              "value": "null"
            }
          ],
          "valueName": "avg"
        }
      ],
      "schemaVersion": 21,
      "style": "dark",
      "tags": [],
      "templating": {
        "list": []
      },
      "time": {
        "from": "now-6h",
        "to": "now"
      },
      "timepicker": {},
      "timezone": "",
      "title": "Example App Dashboard",
      "uid": "DAs6rtcGz",
      "version": 1
    }
  name: simple-dashboard.json


```

make sure the customize dashboard created

![](go13.PNG)


Reference
-----

https://www.redhat.com/en/blog/custom-grafana-dashboards-red-hat-openshift-container-platform-4