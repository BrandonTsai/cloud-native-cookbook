
還只會用 "RollingUpdate" 嗎？ 快來看看如何在 OpenShift 使用進階的 "Blue-Green" 跟 "A/B" 部署策略吧！
====

除了 RollingUpdate 和 Recreate 外，在 OpenShift 您還可以根據 Route 的設置輕易地實現 “Blue-Green” 和 “A/B” 部署策略，在部署和測試上更加靈活。本篇將介紹如何實現這兩種部署策略。

> 本篇文章圖片接來自 [DevOps with OpenShift](https://www.openshift.com/resources/ebooks/devops-with-openshift/) 這本免費的電子書，有興趣的人可以去下載來閱讀。

事前準備
---------

兩個版本的 Deployments 和 Service

版本 1

```
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-1
  labels:
    app: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx-1
  template:
    metadata:
      labels:
        app: nginx-1
    spec:
      containers:
      - name: nginx
        image: quay.io/brandon_tsai/testlab:1
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-1
spec:
  selector:
    app: nginx-1
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
```

版本 2

```
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx2-index
data:
  index.html: |
    <html>
    <head>
            <title>Test NGINX 2 passed</title>
    </head>
    <body>
    <h1>NGINX 2 is working</h1>
    </body>
    </html>
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-2
  labels:
    app: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx-2
  template:
    metadata:
      labels:
        app: nginx-2
    spec:
      containers:
      - name: nginx
        image: quay.io/brandon_tsai/testlab:1
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
        volumeMounts:
        - name: nginx2-index
          mountPath: /opt/app-root/src
      volumes:
        - name: nginx2-index
          configMap:
            name: nginx2-index
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-2
spec:
  selector:
    app: nginx-2
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080

```


一個最基本的 Route 對應到版本 1 

```
apiVersion: v1
kind: Route
metadata:
  name: nginx
spec:
  host: nginx-uat.apps-crc.testing
  to:
    kind: Service
    name: nginx
```

確定這個基本的 Route 是可以正確地運作。

```
$ oc apply -f route.yml
route.route.openshift.io/nginx created

$ oc get route
NAME    HOST/PORT                    PATH   SERVICES   PORT    TERMINATION   WILDCARD
nginx   nginx-uat.apps-crc.testing          nginx-1    <all>                 None

$ curl http://nginx-uat.apps-crc.testing
<html>
<head>
	<title>Test NGINX passed</title>
</head>
<body>
<h1>NGINX is working</h1>
</body>
</html>
```




Blue-Green 部署策略
--------------

Blue-Green 的部署策略通過確保在部署過程中有兩個可用的應用程序版本來最大程度地減少執行部署轉換所需的時間。
我們可以利用服務和路由層輕鬆地在兩個正在運行的應用程序版本之間進行切換，因此執行 Rollback 非常簡單快捷。
對於 `Stateless` 架構的應用程式，Blue-Green 部署很容易實現實現，因為您不必擔心需要與應用程序一起遷移或 Rollback 的儲存數據。 如下圖所示：

![](images/04_OCP_INTRO/route_gb.png)



對於之前部署的 Route，我們可以藉由更改目標 Service 的名稱來很快速的切換到 *版本2* 。

```
$ oc patch route/nginx -p '{"spec":{"to":{"name":"nginx-2"}}}'
route.route.openshift.io/nginx patched

$ curl http://nginx-uat.apps-crc.testing
<html>
<head>
        <title>Test NGINX 2 passed</title>
</head>
<body>
<h1>NGINX 2 is working</h1>
</body>
</html>
```

如果*版本2*有問題，你也可以快速的在切換回*版本1* 。

```
$ oc patch route/nginx -p '{"spec":{"to":{"name":"nginx-1"}}}'
route.route.openshift.io/nginx patched

$ curl http://nginx-uat.apps-crc.testing
<html>
<head>
	<title>Test NGINX passed</title>
</head>
<body>
<h1>NGINX is working</h1>
</body>
</html>

```

A/B 部署策略
------------


A / B部署策略提供了測試新應用程序功能的能力。 這樣，您可以讓部分客戶連接到新版本，以測試您的假設功能是正確還是錯誤，然後回滾到初始應用程序（版本1）或繼續使用新應用程序（版本2）。

如下圖所示：

![](images/04_OCP_INTRO/route_ab1.png)

![](images/04_OCP_INTRO/route_ab2.png)


我們可以利用指令來設定到各版本的流量的比例：

```
$oc annotate route/nginx haproxy.router.openshift.io/balance=roundrobin
$ oc set route-backends nginx nginx-1=50 nginx-2=50
```

測試是否有一半流量被導到*版本2*


```
$ for i in {1..10}; do curl -s http://nginx-uat.apps-crc.testing | grep "<h1>" ; done
<h1>NGINX is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX is working</h1>
<h1>NGINX 2 is working</h1>
```

也可以隨時 *版本1* 調整 *版本2* 的比例。

```
$ oc set route-backends nginx nginx-1=20 nginx-2=80
route.route.openshift.io/nginx backends updated

$ for i in {1..10}; do curl -s http://nginx-uat.apps-crc.testing | grep "<h1>" ; done
<h1>NGINX 2 is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX is working</h1>
```


結論
-----------
與 Kubernetes 相比，在 OpenShift 平台上實現 "Blue-Green" 跟 "A/B" 部署策略容易得多了。

參考資料
---------
free ebook: [DevOps with OpenShift](https://www.openshift.com/resources/ebooks/devops-with-openshift/)