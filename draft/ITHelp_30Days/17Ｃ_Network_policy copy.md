網路安全，人人有責： Network Policy
========

NetworkPolicy 是一組可用來定義及管理 pod 之間連線的規範。預設的情況下，Pod 透過 Service 開放端口後，平台內所有其他的 Pod 都可以透過此 Service 連接到此 Pod，即使是在不同的 Namespace。在 kubernetes，你必須要設定支援此功能的 CNI plugin 才可以使用 NetworkPolicy，而 OpenShift 則是已經將此功能內建於平台內。

如果Pod被一個或多個對像中的選擇器匹配，則 Pod 將僅接受那些 NetworkPolicy 對像中的至少一個允許的連接。 完全可以訪問任何 NetworkPolicy 對象未選擇的Pod。
NetworkPolicy 透過 Label 的方式來選擇目標要管理的Pod，被選中的Pod只會接受被 NetworkPolicy 允許的連線來源，不被允許的則無法連接到此 Pod。相對的，沒有被 NetworkPolicy 選擇中的Pod，則可以接受來自任何地方的連線。

以下我將放上幾個範例來透過 NetworkPolicy 管理連線。

範例一：只接受來自同一個 Project/Namespace 的 Pods 的連線。
-----------------------

如果 “UAT” Project 中的 Pod 知道 "PROD" 項目中另一個 Pod 的 clusterIP，則實際上它可以通過 clusterIP 連接到該 Pod。


```
## Go into UAT project
$ oc project uat
Now using project "uat" on server "https://api.crc.testing:6443".

## get the CLUSTER-IP of pods in PROD project
$ oc get svc -n prod
NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
nginx-prod   ClusterIP   172.25.45.167   <none>        80/TCP    6m23s

## query the PROD service inside a UAT pod
$ oc rsh nginx-uat-5fb55b7d79-hgzzl curl http://172.25.45.167
<html>
<head>
        <title>Test NGINX Prod passed</title>
</head>
<body>
<h1>NGINX Prod is working</h1>
</body>
</html>

```

為了避免此情況，我們需要讓 Project 是確實獨立隔離的，我們就可以在每個 Project 設定下列 NetworkPolicy：



```
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: allow-same-namespace
spec:
  podSelector:
  ingress:
  - from:
    - podSelector: {}
```

部署後你會發現不同 Project 的 Pod 無法再直接都透過 clusterIP 連接。

```
$ oc apply -f np-allow-same-project-only.yml -n prod
networkpolicy.networking.k8s.io/allow-same-namespace created

$ oc rsh nginx-uat-5fb55b7d79-hgzzl curl http://172.25.45.167
curl: (7) Failed connect to 172.25.45.167:80; Connection timed out
command terminated with exit code 7
```



範例二：只接受來自特定的 Pods 的連線。
-----------------------

譬如，對於某些後端的 Pods, 我們希望他只可以接受來自 nginx Pod 的連線，這時我們可以使用以下 NetworkPolicy 來規範。


```
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: allow-pod-and-namespace-both
spec:
  podSelector:
    matchLabels:
      app: api-gw-uat
  ingress:
    - from:
      - podSelector:
          matchLabels:
            app: nginx-uat
```


在部署前, 所有的 Pods 都可以連接 `api-gw-uat` 服務。 

```
$ oc rsh nginx-uat-5fb55b7d79-hgzzl curl -s -o /dev/null -w "%{http_code}\n"  http://api-gw-uat
200

$ oc rsh test-pod curl -s -o /dev/null -w "%{http_code}\n"  http://api-gw-uat
200

```

部署後，只有 nginx pods 可以成功連接到 `api-gw-uat` 服務。

```
$ oc apply -f np-allow-pod.yml
networkpolicy.networking.k8s.io/allow-pod-to-api-gw created

$ oc rsh nginx-uat-5fb55b7d79-hgzzl curl -s -o /dev/null -w "%{http_code}\n"  http://api-gw-uat
200

$ oc rsh test-pod curl -s -o /dev/null -w "%{http_code}\n"  http://api-gw-uat
000
command terminated with exit code 7
```


範例三：只允許 HTTP 及/或 HTTPS 連線。
------------------


可以設定其他Pod就是只能透過 HTTP 及/或 HTTPS 連線。

```
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: allow-https-only-to-nginx
spec:
  podSelector:
    matchLabels:
      app: nginx-uat
  ingress:
  - ports:
    - protocol: TCP
      port: 8443
```

以上範例只允許其他 Pod 對 目標 Pod 的 8443 端口溝通。

```
$ oc apply -f np-allow-https-only.yml
networkpolicy.networking.k8s.io/allow-https-only-to-nginx created

$ oc rsh test-pod curl -s -o /dev/null -w "%{http_code}\n" -k https://nginx-uat
200

## Other pods must use https port to connect to nginx pod
$ oc rsh test-pod curl -s -o /dev/null -w "%{http_code}\n" http://nginx-uat
000
command terminated with exit code 7

```




結論
------

在 Kubernetes 為基礎的其他平台, 你可以透過安裝 [Calico](https://docs.projectcalico.org/getting-started/kubernetes/) 來達到相同目的。



參考資料
----------

- https://docs.openshift.com/container-platform/4.5/networking/network_policy/about-network-policy.html