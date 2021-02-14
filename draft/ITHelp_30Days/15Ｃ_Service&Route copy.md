
免除 Ingress Controller 煩惱，擁抱 OpenShift Route 新世界
==============================


在傳統VM的環境中，如果要公開新服務，則需要為新VM分配IP，並將其添加到DNS服務器，設置Nginx以啟用TLS證書，更新HAProxy配置並重新加載HAProxy 服務來處理負載平衡。 在有些公司這並不是經常性的需求，所以不那麼讓人感到厭煩。 但是，如果您是在一個每天都會更新或發布許多新服務的公司，那麼您可能很快就會發現這種重複性工作會佔去你大部分的時間。

Kubernetes 中的網絡模型可通過解決以下問題來幫助您減輕這種痛苦：

- 在 Pod 內的容器可透過 127.0.0.1 來跟彼此溝通。
- 在不同 Pod 之間可以透過內部的 `Service` 來溝通 。
- 設定好 `Ingress Controller` 後，可以讓平台外的人透過 `Ingress` 資源連線到內部的 `Service`。

 OpenShift, 不像 Kubernetes, 你需要先設定一個 `Ingress Controller` 才可以使用 `Ingress` 資源, OpenShift 內建一個解決方案給外部連線使用，叫 `Route`.

概念圖如下:

![](images/04_OCP_INTRO/network_1.png)



Service 概念淺談
----------

Service 是將運行在一組 Pod 上的應用程序公開為網絡服務的抽象方法。
Kubernetes 為 Pods 提供一個 IP 地址和單個DNS名稱給一組Pod 使用，並且可以在它們之間進行負載平衡。


以下範例將創建一個名為 “nginx” 的 Service，該對象的目標是任何帶有`app = nginx` 標籤的Pod上的TCP端口8080。

```
apiVersion: v1
kind: Service
metadata:
  name: nginx
spec:
  selector:
    app: nginx
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
```


OpenShift 容器平台具有內置的 DNS 服務，在其他Pod中，他們可以通過此內置 DNS 服務訪問該新建的“nginx”服務，例如：

```
$ curl http://nginx
<html>
<head>
	<title>Test NGINX passed</title>
</head>
<body>
<h1>NGINX is working</h1>
</body>
</html>
```



OpenShift Route
----------------------

Kubernetes的Ingress資源公開了從外部到內部 Service 的 HTTP 和 HTTPS 路由，而且Ingress資源可以定義路由規則。

OpenShift 出於和 Kubernetes Ingress 相同目的創建了一個名為 “Route”的資源，但它具有其他功能，例如在多個內部 Service 之間分離流量等。在 OpenShift 上創建Route對象時，它會由內建的 HAProxy load balancer 來負責管理以公開所請求的服務，並使該服務在給定配置下可從外部連線使用。

Ingress 跟 Route 比對圖
![](images/04_OCP_INTRO/network_ingress_route.png)

以下我將示範幾種不同的情境建立 Route 的用法：


### 範例一： 僅開放 HTTP 連線

首先先建立基本的 nginx Service 你可以透過指令：

```
$ oc expose svc nginx
route.route.openshift.io/nginx exposed
```

或 透過 YAML 檔案

```
apiVersion: v1
kind: Route
metadata:
  name: nginx
spec:
  to:
    kind: Service
    name: nginx
```

部署 YAML 檔案
```
$ oc apply -f route.yml 
route.route.openshift.io/nginx created

$ oc get route
NAME    HOST/PORT                          PATH   SERVICES   PORT    TERMINATION   WILDCARD
nginx   nginx-myproject.apps-crc.testing          nginx      <all>                 None
```

如果你沒有特別指定 hostname， OpenShift 會自動產生 hostname  `<route-name>-<namespace>.<external-address>`。
你也可以客製化 hostname, 但客製化的 hostname 必須是 `external-address` 的 subdomain， 例如 `<my-hostname>.<external-address>`。


我們可以在 YAML 檔案中直接指定要客製化的 hostname。

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


檢驗我們可以從瀏覽器讀取此網址，或透過 curl 指令測試：

```
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

### 範例二：TLS with edge termination

最簡單快速開放 HTTPS 連線的方式是利用 edge termination，預設是使用平台的 "Widecard Certificate"，所以不需要再自己申請憑證（Certificate）.

指令：

```
$ oc create route edge --service=nginx --hostname=nginx-uat.apps-crc.testing
route.route.openshift.io/nginx created
```

YAML 檔案：

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
  tls:
    termination: edge
```




如果你想要使用自定義的 Certificate 而非平台提供的 "Widecard Certificate"， 你可以透過以下指令 `oc create route edge --service=nginx --hostname=nginx-uat.apps-crc.testing  --key=nginx-uat.key --cert=nginx-uat.crt`。

如果你想要把 http 連線都自動導到 https， 你還可以加上參數 `--insecure-policy=Redirect`。


然後就可以透過 curl 參數測試：

```
$ curl -k https://nginx-myproject.apps-crc.testing
<html>
<head>
	<title>Test NGINX passed</title>
</head>
<body>
<h1>NGINX is working</h1>
</body>
</html>
```


### 範例三：TLS with Passthrough Termination

如果你想要到內部 Pod 這端都一直保持 HTTPS 連線，你可以使用 `Passthrough Termination`。

在這個模式，你的 deployment 一定要開啟 8443 端口，而對應的 service 也必須開啟 443 端口。


首先我們把修改 nginx 設定檔讓它支援 https 。

```
server {
        listen       8443 ssl http2 default_server;
        listen       [::]:8443 ssl http2 default_server;
        server_name  _;
        root         /opt/app-root/src;

        ssl_certificate "/opt/app-root/tls/tls.crt";
        ssl_certificate_key "/opt/app-root/tls/tls.key";


        location / {
        }

        error_page 404 /404.html;
            location = /40x.html {
        }

        error_page 500 502 503 504 /50x.html;
            location = /50x.html {
        }
    }
```

從指令建立 ConfigMap。

```
$ oc create configmap nginx-config --from-file=nginx/https.conf -n myproject --dry-run=client -o yaml | oc apply -f -
```

從指令為 TLS certificate 建立 Secret。

```
$ oc create secret tls nginx-tls --key="nginx-uat.key" --cert="nginx-uat.crt" -n myproject --dry-run=client -o yaml | oc apply -f -
```

更新 YAML，把 ConfigMap 和 Secret 掛載到 deployment。

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: quay.io/brandon_tsai/testlab:1
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
        - containerPort: 8443
        volumeMounts:
        - name: nginx-tls
          mountPath: /opt/app-root/tls
        - name: nginx-config
          mountPath: /opt/app-root/etc/nginx.d/
      volumes:
        - name: nginx-tls
          secret:
            secretName: nginx-tls
        - name: nginx-config
          configMap:
            name: nginx-config
```

更新 service 設定

```
apiVersion: v1
kind: Service
metadata:
  name: nginx
spec:
  selector:
    app: nginx
  ports:
    - name: port-http
      protocol: TCP
      port: 80
      targetPort: 8080
    - name: port-https
      protocol: TCP
      port: 443
      targetPort: 8443
```


最後透過指令建立 Route

```
$ oc create route passthrough --service=nginx --hostname=nginx-uat.apps-crc.testing --port=port-https
route.route.openshift.io/nginx created
```

或透過 YAML 檔

```
apiVersion: route.openshift.io/v1
kind: Route
metadata:
  name: nginx
spec:
  host: nginx-uat.apps-crc.testing
  port:
    targetPort: port-https
  tls:
    termination: passthrough
  to:
    kind: Service
    name: nginx
status: {}
```





### 範例四： TLS with Re-encryption Termination**

Re-encryption 是指從外部到 Route 端的HTTPS路由是由一個憑證加密，然後它再Route這邊解密後會再重新由另一個憑證加密內部到Pod端的連線。


這次使用跟範例三一樣的 deployment 跟 service , 但我們把 Route 改為 reencrypt。
外部用戶到 Openshift 之間的連線會是透過平台預設的 "Widecard Certificate" 加密。
OpenShift 內部的連線則是透過我們自己的憑證加密。

指令：
```
oc create route reencrypt --service=nginx --hostname=nginx-uat.apps-crc.testing --port=port-https  --dest-ca-cert=rootCA.crt
```

YAML 範例：
```
apiVersion: route.openshift.io/v1
kind: Route
metadata:
  name: nginx
spec:
  host: nginx-uat.apps-crc.testing
  port:
    targetPort: port-https
  tls:
    destinationCACertificate: |
      -----BEGIN CERTIFICATE-----
      MIIFRDCCAywCCQCe2OOaDBEvLjANBgkqhkiG9w0BAQsFADBkMQswCQYDVQQGEwJB
      VTEMMAoGA1UECAwDTlNXMQ8wDQYDVQQHDAZTeWRuZXkxEjAQBgNVBAoMCURhcnVt
      YXRpYzELMAkGA1UECwwCSVQxFTATBgNVBAMMDGJyYW5kb24udGVzdDAeFw0yMDA5
      MjcwNzU5NDBaFw0yMzA3MTgwNzU5NDBaMGQxCzAJBgNVBAYTAkFVMQwwCgYDVQQI
      DANOU1cxDzANBgNVBAcMBlN5ZG5leTESMBAGA1UECgwJRGFydW1hdGljMQswCQYD
      VQQLDAJJVDEVMBMGA1UEAwwMYnJhbmRvbi50ZXN0MIICIjANBgkqhkiG9w0BAQEF
      AAOCAg8AMIICCgKCAgEAxl7AtuZa/kXDQKNsgIYbHCvDhOUXW7Jvz8WVMAL94/Fe
      lcktvieClHzIkBYk599G3INpsEBEiersKGyPjlIBPqmrfDJmlZSnpZwnWFhrBbIs
      /EouQe4t6LsqUg+Jj9WpTPSFGAzxqn8OZrMUoMOLj8xRxp8p85ziV9t6CZtfwET6
      laj+Cv7MznsNn8R+cgK2YW+516W7YQgg1szoucBlldoKRR4Xya7h4VcfNa3s4uKx
      RPoBUmLnV5Edes/BCUjwFtC7lenzNjc+mO9El75XGPJxZY+NtTonQ0v5L4rgzUsW
      3Nz2nR36NwOwR+buq/tfodRwR29ZqlJ4mHBDrJntmWmfqnR+WAu6Dsbwt0YTFj1i
      XDRrbXSTHz5Efu+2IQv6wiUdczHD958MZBBzNTpCr7Ss+4gvSBgiVM2yvwiZQJeg
      2I3147d+hz/57J0mLc07tbDQ0pGnTWyAHXvEm7KlO0yZIaTLY6SReTVQsqq16uJT
      flhZWEz13fn3axQpD/OTSXJIRbRyusVJrJKglbFpuUeGjbR/I3K32sZZti4fBpDZ
      ldSnCpuR/z27iBoTpHHg5Aa6SBRO5TjI91yUBdNxw2NHtVzSoY3Z0fGMNfG8hSrK
      DYD2tXJx4k1upiT68HMMYA2kdGVblisBPygL3VM9eZXtoQmcy9BjixhaTWOCyfMC
      AwEAATANBgkqhkiG9w0BAQsFAAOCAgEAtxnHufjkm8ZsQ3aftdsrm8sHL5XUzlYf
      RyH60QLK0Gjl19FkG/sS/XminoZZO0PFFb/Z78L+KVezMj6bd6FNc3ULiKmssQA0
      9Pvzr4c6dyXRapMRWArGCrfYns8vPy8TAJ9DDGV+VNHI2L0VTPk9h/a1qp8qAXmp
      XM0tVCZQrVFc7e6DeCfYYZ/ukAj2n70jUm5iuDTkM5OgbE9XQrbgJeJnGEhm5XrY
      mmE9+G+VXxoAkV2EaNVAzHTg3AeywrLlWArKPL4vl/pG15u5xDDRG074Tkb8gyRS
      UZ+Nc9lNDs2Rw0dsP5E7njnNYkQU81XxgsIu96XSbws0Z5GTGpeHk+CNVycQ9wOV
      Kdyc7aosKxzGUQi69gTa4xmn+EGsboUOappo4fTkP3TeetUYk/79q7AZxxgOTHkr
      fLgrNcrjiUiW91U4ma6PtHbnlzNCl7MYZyy+sxLogR3NFnO8xOtK/1Xdkrtf9/YI
      NmzFsSaumCtxbVsYrTvZMt7eVkJUKL3Kx4K1Vs51emEwtsPB/HJh5ozY2fk3rjlw
      GQ/TM3lH3dViUkhh5DJiGnYU05lOP7aZKR2yWxlqMkdMpnrUq6tF392s1585YdSO
      ohRM5gMVGrL95F/vXln3e4vX4mA4bgr9LLj5xkELJR7UEIXHV6nCeCSutohzejzs
      CwdaB0xYUqI=
      -----END CERTIFICATE-----
    termination: reencrypt
  to:
    kind: Service
    name: nginx
  ```


結論
----------

這篇文章介紹了 OpenShift 中 Service 和 Route 的基本概念。 文章中包含了4個範例，適用用於大多數情況，下一篇我會介紹基於 OpenShift Route 的兩種進階部署策略。


參考資料
---------
- https://docs.openshift.com/container-platform/3.11/dev_guide/routes.html
https://docs.openshift.com/container-platform/3.11/architecture/networking/routes.html#route-types
- https://www.openshift.com/blog/kubernetes-ingress-vs-openshift-route