
我不是針對誰，我是說在座各位都應該用 ConfigMap 或 Secret
===========================


許多應用程序都需要使用配置文件(Configuration File)和環境變數(Environment Variables)的來設定及運行應用程序。筆者看過許多該入坑的工程師會根據不同應用場景，直接將對應的配置文件和環境變數直接寫死並一起放入映像檔內，然後再根據應用場景來使用不同的映像檔。這不僅造成安全性上的風險，也完全沒有必要．這些配置文件和環境變量應與容器映像檔分離，以保持容器化應用程序的可攜性。你應該使用同一個映像檔，但在運行容器時在根據應用場景來嵌入這些配置文件和環境變數。

ＯpenShift 是建立在 Kubernetes 基礎上的平台，我們可以使用 `ConfigMap` 或 `Secret` 元件對運行在 Pod 內的容器嵌入這些配置文件和環境變數，已達到應用程序的可攜性。 `ConfigMap` 是用於存儲非機密性的數據。 如果您要存儲的數據是機密的，例如SSL金鑰，請使用 `Secret`，或使用其他第三方工具（例如 [Conjur](https://www.conjur.org/), 或 [Hashicorp Vault](https://www.vaultproject.io/)）來管理並保密機密數據並嵌入容器內。


從資料夾或檔案來建立 ConfigMap
---------------------------------------

提示:

- 檔案可以是 Binary。
- 無法從子資料夾或多層架構的資料夾建立 ConfigMap，你應該先平坦化配置文件的資料夾(Flatten directory structure)。
- 在你的 CI/CD Pipeline 中利用 "--dry-run=client -o yaml" 來產生對應的 YAML 然後再部署到 OpenShift, 例如： `oc create configmap nginx-config --from-file=configs/nginx -n myproject --dry-run=client -o yaml | oc apply -f -`


```bash
＃ From generic file.
$ oc create configmap nginx-config --from-file=configs/uat/nginx.conf -n myproject --dry-run=client -o yaml | oc apply -f -

# From Folder
$ oc create configmap nginx-icons --from-file=configs/uat/icons -n myproject --dry-run=client -o yaml | oc apply -f -

# From env file.
$ oc create configmap nginx-env --from-env-file=configs/uat/nginx.env -n myproject --dry-run=client -o yaml | oc apply -f -
```


從資料夾或檔案來建立 Secrets
---------------------------------------

跟建立 ConfigMap 很像，但 Secrets 針對不同三種用途提供三中不同的格式。

| Commands | Usage |
|----------|-------|
| docker-registry | 用來從私有映像倉庫拉映像檔 |
| tls             | 針對 TLS 金鑰 |
| generic         | 一般用途，可從檔案，資料夾 或 字面值（literal value） |

例如：

```bash

$ oc create secret tls nginx-tls --key="configs/uat/tls.key" --cert="configs/uat/tls.crt" -n myproject --dry-run=client -o yaml | oc apply -f -

$ oc create secret generic pgsql-secret --from-literal pgsql_user=brandon --from-literal pgsql_key=1234U987 --dry-run -o yaml | oc apply -f -
```



在 Pod 中使用 ConfigMap & Secrets 
-----------------------------

ConfigMap 跟Secret 可以被用來:

- 嵌入環境變數於容器內
- 以 Volume 方式掛載配置文件於容器內

請注意，必須先創建ConfigMap 和 Secret，然後才能在 Pod 中使用其內容。

如以下範例:


```yaml
kind: Deployment
apiVersion: apps/v1
metadata:
  name: nginx-app1
  namespace: myproject
  labels:
    app: nginx-app1
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx-app1
  template:
    metadata:
      labels:
        app: nginx-app1
    spec:
      containers:
      - name: nginx-app1
        image: registry.redhat.io/rhscl/nginx-114-rhel7:1
        ports:
        - containerPort: 8001
          protocol: TCP
        env:
          - name: OS_TYPE
            valueFrom:
              configMapKeyRef:
                name: example-env
                key: os.type
        envFrom:
          - configMapRef:
              name: nginx-env
          - secretRef:
              name: pgsql-secret
        volumeMounts:
        - name: nginx-tls
          mountPath: /opt/app-root/etc/nginx.d/ssl
        - name: nginx-config
          mountPath: /opt/app-root/etc/nginx.d/
        - name: nginx-icons
          mountPath: /opt/app-root/src/icons
      volumes:
        - name: nginx-tls
          secret:
            secretName: nginx-tls
        - name: nginx-config
          configMap:
            name: nginx-config
        - name: nginx-icons
          configMap:
            name: nginx-icons
```

使用 Secret 的風險
----------

- 如果您通過JSON或YAML文件來部署Secret，該文件的數據編碼為base64，則共享此文件或將其上傳到Git Repository意味著該秘密已被盜用。 Base64編碼不是一種加密方法，應該被認為與純文本相同。
- 容器從Secret讀取機密性的數據後，應用程序仍然需要保護機密性的數據，例如不要意外地將其記錄在Log中或傳輸給不受信任的一方。
- 可以創建Pod並使用Secret的用戶也可以透過此方式來看到該機密的值。
- 目前，在任何節點上具有root權限的任何人都可以通過模擬kubelet來從API服務器讀取任何秘密。
- 在 K8S 基本架構中, Secret 的資料是儲存在 `etcd` 中; 因此:
  - 系統管理員應該對整個 K8S 資料啟用加密(requires v1.13 or later).
  - 系統管理員應將對etcd的訪問權限限制為管理員用戶。
  - 對於運行etcd的磁盤不再使用時，管理員應該要將其粉碎使其無法復原或讀取。
  - 系統管理員應確保將SSL / TLS用於etcd間的通信。

目前 ＯpenShift 跟 Kubernetes 預設都沒有對 `etcd` 加密。 ＯpenShift的系統管理員在架構OpenShift時，應參考
[這篇文章](https://docs.openshift.com/container-platform/4.5/security/encrypting-etcd.html#enabling-etcd-encryption_encrypting-etcd)對保護系統內的資料，Kubernetes的系統管理員則應參考https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/。

有餘力且重視資料安全性的團隊，我建議考慮整合其他第三方工具（例如 [Conjur](https://www.conjur.org/), 或 [Hashicorp Vault](https://www.vaultproject.io/)）來管理並保密機密數據並嵌入容器內。


結論
---
不管有沒有用ＯpenShift/Kubernetes平台，保持容器化應用程序的可攜性都是很重要的，如果你只是單純在 Server 上跑 Docker 來部署服務，那麼你可以把放配置文件的資料夾在運行時再掛載上去。而對於許多習慣用多層架構的資料夾來存放配置文件的開發者，在轉移到ＯpenShift/Kubernetes時必定要先經過一番努力來平坦化配置文件的資料夾，目前我也沒有看到ＯpenShift/Kubernetes 是否有規劃之後要支援此類情況的物件，也沒看到第三方工具來幫助解決這個狀況。所以目前我只能說 『Don't Be Afraid to Get Your Hands Dirty！』，加油！
！[]http://www.ghibli.jp/gallery/chihiro027.jpg