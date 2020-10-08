

這次我要安裝跟 Quay 有關的一個 Operator 叫 "Container Security Operator"。
它可以幫忙檢查目前在 OpenShift 平台上運行的 Pod 是否有 Vulnerability。
我們可以從 Web UI 的 OperatorHub 來直接安裝，或透過指令來安裝。

透過指令來安裝 Operator 
----------


（1.） 找尋要安裝的 Operator 

```
# Search operator
$ oc get packagemanifests -n openshift-marketplace | grep "security"
container-security-operator                 Red Hat Operators     22d
```

（2.） 確認 Operator 內容以取得 `channel` 跟 `CSV` 資訊

```
# check operator info
$ oc describe packagemanifests container-security-operator -n openshift-marketplace

```


（3.） 建立 `Subscription` 物件：

YAML 檔

```
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: container-security-operator
  namespace: openshift-operators
spec:
  channel: quay-v3.3
  installPlanApproval: Automatic
  name: container-security-operator
  source: redhat-operators
  sourceNamespace: openshift-marketplace
  startingCSV: container-security-operator.v3.3.1
```

部署到 OpenShift

```
$ oc apply -f sub.yml 
subscription.operators.coreos.com/container-security-operator created
```


（4.） 檢查 Operator 的 Pod 是否有正常運行

```
$ oc get pods
NAME                                           READY   STATUS              RESTARTS   AGE
container-security-operator-5d8c9c64d6-cf99d   1/1     Running             0          16s
devworkspace-controller-84f877b-4d9x9          1/1     Running             0          2d3h
devworkspace-webhook-server-7d6645dc7b-6nkd9   1/1     Running             1          2d2h
```


從 Web UI 測試 "Container Security Operator"
------------------------------

1. 確定 Operator 有安裝成功

![](o1.png)

2. 從 “OverView” 頁面 看整個系統有多少 "Image Vulnerabilities"

![](02.png)

3. 檢查一個 Pod 的 Vulnerabilities

![](3.png)


從 CLI 測試 "Container Security Operator"
------------------------------

搜尋整個系統所有找到的 CVEs

```
$ oc get imagemanifestvuln --all-namespaces
NAMESPACE                   NAME                                                                      AGE
myproject                   sha256.b176867581c15c7bf937757df9207dcd25924789a640af2ed1837a317f3ace25   71m
openshift-cluster-version   sha256.7ad540594e2a667300dd2584fe2ede2c1a0b814ee6a62f60809d87ab564f4425   26h
```

列出偵測到有 Vulnerabilities 的 Pod

```
$ oc get imagemanifestvuln -o json | jq '.items[].status.affectedPods' | jq -r 'keys[]' | sort -u
myproject/nginx-77f76d44b8-552mh
```

列出偵測到有 Vulnerabilities 的 Image 名稱

```
$ oc get imagemanifestvuln -o json | jq -r '.items[].spec.image' | sort -u
quay.io/brandon_tsai/testlab
```


結論
-----------

偵測風險不應該只有在建立映像檔後及部署到 OpenShift 前做而已，部署到 OpenShift 後持續地監控整個系統的風險是非常重要的。透過這個 Operator，我們可以很輕易地找出系統中潛在的風險以方便做進一步的管理。
