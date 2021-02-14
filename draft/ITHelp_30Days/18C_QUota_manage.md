Openshift 運算資源管理淺談： Pod QoS, Quota, Multi-Project Quotas, 跟 LimitRange.
=================

傳統上，運營工程師和開發人員通常需要先選擇虛擬機的大小以運行應用程序，這些應用程序會在虛擬機中獨立運行，並使用該虛擬機所有資源。但是在Openshift中，Pod 可以在任何與他人共享資源的運算節點上運行，這導致了可能在某些時候，該運算節點的資源會不夠所有在該節點上運行的 Pod 使用，而出現問題。 這就凸顯了事先規劃 QoS（服務質量等級）和資源配額的重要性。


Resource Request and Limits
------------------------------
為應用程序創建Pod時，可以為其中的每個容器的 CPU 和 Memory 等資源宣告請求(Request)和限制(Limits)。正確設置這些值是讓 Kubernetes 有效管理系統資源並分配給應用程序的第一步。

宣告方式如

```
spec:
  containers:
  - image: openshift/hello-openshift
    name: hello-openshift
    resources:
      requests:
        cpu: 100m
        memory: 200Mi
        ephemeral-storage: 1Gi
      limits:
        cpu: 300m
        memory: 800Mi
        ephemeral-storage: 2Gi
 ```

**Requests**: 這些值用是宣告容器需要運行的最少資源，OpenShift 可根據此值來決定要將此容器運行在哪個運算節點上， 如果沒有節點有足夠的資源來處理請求，則Pod會保持“待處理（Pending）“狀態。 

**Limits**: 運算節點將允許容器使用的最大資源量。
- 如果容器嘗試使用的資源超過指定的限制，則系統將阻擋分配給此容器更多資源．
- 如果容器超過了指定的 Memory 限制，則它將被終止，並有可能根據容器重新啟動策略重新啟動。

Quality of Service Classes (QoS)
------------------


當節點運行了一或多個沒有宣告資源請求的Pod時，或者該節點上所有Pod的限制總和超過該節點可用的計算資源時，該節點被稱為“過度使用”。
在過度使用的環境中，節點上的Pod可能在任何時間點會嘗試使用比該節點可用的計算資源更多的計算資源。
發生這種情況時，該節點對容器設定優先順序。 優先級最低的容器首先被終止/節流。 用來做出此決定的工具稱為服務質量（QoS）。



| Priority | Class Name | Description |
|----------|------------|-------------|
| 1 (highest) | Guaranteed | If limits and optionally requests are set (not equal to 0) for all resources and they are equal. |
| 2 | Burstable | If requests and optionally limits are set (not equal to 0) for all resources, and they are not equal |
|3 (lowest) | BestEffort | If requests and limits are not set for any of the resources |


因此，如果開發人員未聲明任何資源請求和限制，則該容器將首先終止。 
我們應該透過設置請求和限制值來保護正式環境中的主要容器，以便將其QoS分類為“Guaranteed”。 “BestEffort” 或“ Burstable”只能在開發環境中使用。


Quota 
------

管理員可以針對每個 Project 設定配額（Quota）以限制資源的使用量。
這具有附加作用； 如果您在配額中設置了Memory請求總量，那麼所有Pod都需要在其定義中設置Memory請求。如果新的Pod嘗試分配的資源超出配額限制，則將不會進行調度，並且將保持在“待處理（pending）”狀態。



```
apiVersion: v1
kind: ResourceQuota
metadata:
  name: compute-resources
  namespace: ocp-backup
spec:
  hard:
    # pods: "4" 
    requests.cpu: "1" 
    requests.memory: 2Gi 
    limits.cpu: "1" 
    limits.memory: 2Gi 
```

Multi-Project Quota
------

類似配額（Quota），但 Quota 是針對每個 Project分配資源配額，而 Multi-Project quota 顧名思義就是多個Project共享資源配額。



```
apiVersion: v1
kind: ClusterResourceQuota
metadata:
  name: gts-dev
spec:
  quota: 
    hard:
        requests.cpu: "2"
        requests.memory: "4Gi"
        limits.cpu: "2"
        limits.memory: "4Gi"
  selector:
    annotations: 
      gts.automation/clusterResourceQuota: gts-dev
    labels: null
```

LimitRanges
------

LimitRange可用來對Project中的Pod或Container近一步管理資源使用量。 它可以對 Project 內的 每個 Container：

- 設定預設的 request/limit 值。
- 規定request/limit允許的最小和最大值。
- 規定 request 和 limit 值之間最大的比例，例如 `limit值/request值` 必須小於3。


結論
-----

Warning: 筆者建議系統管理員不要在正式環境中透過 LimitRange 設定預設 request/limit 值。 應該要強迫開發者和營運工程師在正式環境中自己對每個 container 設定合適的request/limit值。因為如果預設值太高，那麼 Pod 很容易不會被調度，並且一直保持在“待處理（pending）”狀態。相反的，如果預設值太低，那麼 Pod 很容易被砍掉重練。


參考資料
---------

https://docs.openshift.com/container-platform/3.11/admin_guide/out_of_resource_handling.html
https://docs.openshift.com/container-platform/3.11/admin_guide/overcommit.html#qos-classes
