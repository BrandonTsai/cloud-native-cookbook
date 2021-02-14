
OpenShift 系列 3: 到底是該用 DeploymentConfig，還是用 Deployment 呀，你搞得我好亂呀！
================================================

Deployment vs DeploymentConfig
-----------------------

-----

OpenShift 提供兩種相似但不同的方法來對應用程序進行管理 - `Deployment` 跟 `DeploymentConfig`。
你可以透過他們來描述應用程序的生命週期，例如該應用程序應該運行的Pod數量，更新Pod的方式，及部署失敗時如何快速回到上一個版本。以下是兩者的比較

| | Deployment | DeploymentConfig |
|--|------------|------------------|
| 管理Pod數量的物件 |  [ReplicaSet](https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/) | [ReplicationController](https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller/) |
| 選擇算符(Selector) | 支援最新的 **set-based** selector. <BR> 例如: `environment in (production, qa)`  | 僅支援 equality-based selector. <BR> 例如: `environment = production`  |
| CAP 理論 |偏向數據可用性（Availability） <BR> 發生故障期間，其他主服務器可能會同時對同一 Deployment 進行操作。 | 偏向數據一致性（Consistency） <BR> 如果運行Deployer Pod的節點發生故障，它不會被其他節點替換。 該程序會一直等到該節點重新上線或被手動刪除。  |
| 自動 rollbacks | 不支援 | 如果發生故障，可以Rollback到最近成功部署的ReplicationController。 |
| Trigger | 不支援 | 當你使用OpenShift內建的CI/CD 工具或 ImageStrem 功能，則當有任何變動時，它可以自動讓DeploymentConfig 進行更新。  |
| Lifecycle hooks | 不支援 | 可以在 Ｌifecycle 的各個階段自定義要執行的任務。 |
| 自行定義更新策略 | 不支援  | 支援 |


Kubernetes 支援的 Deployments 和 OpenShift 提供的 DeploymentConfig 均可以用來有效的管理Pod。 儘管DeploymentConfig 比 Deployment 具有更多功能，但這些功能其實在大多數情況下都不是很有用，除非您非常需要DeploymentConfig 提供的特定功能或行為，否則建議使用 Deployment。



DeploymentConfig Trial
-----------------------

-----

以下是 DeploymentConfig 範例。

```YAML
apiVersion: v1
kind: DeploymentConfig
metadata:
  name: nginx-dc
  labels:
    app: nginx
spec:
  replicas: 1
  selector:
    app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: quay-eu-uat/application-images/test:1
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
```

部署到Openshift上

```
$ oc apply -f dc.yml
deploymentconfig.apps.openshift.io/nginx-dc configured

$ oc get pods
NAME             READY   STATUS      RESTARTS   AGE
nginx-dc-1-deploy   0/1     Completed   0          86s
nginx-dc-1-mtrrb    1/1     Running     0          83s
```

跟 Kubernetes 支援的 Deployment 不同，每次 DeploymentConfig 被部署或更新時，都會產生一個 `Deployer Pod` 來管理部署（包括縮減舊的 ReplicationController，擴大新的 ReplicationController以及運行Lifecycle hooks）。 `Deployer Pod` 在完成部署後會保留不確定的時間，以保留其部署日誌。



透過 DeploymentConfig 來增加 Pod 數量
----------------------------------------------

-----

你可以使用 `oc scale` 指令來手動設定一個 Deployment 或 DeploymentConfig 管理的 Pod 數量. 例如, 下列指令將 nginx-dc 的 DeploymentConfig 的 Pod 目標數量手動設為 2.

```
$ oc scale dc nginx-dc --replicas=2
deploymentconfig.apps.openshift.io/nginx-dc scaled

$ oc get pods
NAME                READY   STATUS      RESTARTS   AGE
nginx-dc-1-49cpw    1/1     Running     0          2m54s
nginx-dc-1-deploy   0/1     Completed   0          5m45s
nginx-dc-1-mtrrb    1/1     Running     0          5m42s
```



Rollback 到之前的版本
-----------------------

-----

因為 `oc rollback` 命令僅適用於 DeploymentConfig。 我建議平時使用`oc rollout undo` 命令代替 `oc rollback`。

首先，我們通過更改映像檔標籤來更新 “nginx-dc” deploymentConfig。

```
$ oc set image dc/nginx-dc nginx=quay.io/brandon_tsai/testlab:2
```

等待直到更新完成

```
$ oc rollout status -w dc/nginx-dc
```

Rollback 到前一版本並等待直到 Rollback 完成.

```
$ oc rollout undo --to-revision=1 dc/nginx/dc
deploymentconfig.apps.openshift.io/nginx-dc deployment #3 rolled back to nginx-dc-1

$ oc rollout status -w dc/nginx-dc
replication controller "nginx-dc-3" successfully rolled out
```

檢查更新歷史紀錄.

```
$ oc rollout history dc/nginx-dc
deploymentconfig.apps.openshift.io/nginx-dc
REVISION	STATUS		CAUSE
1		Complete	config change
2		Complete	config change
3		Complete	config change
```



結論
---------

-----


希望本篇有讓你了解到 Deployment 跟 DeploymentConfig 之間的微小差距。
記得，除非您非常需要 DeploymentConfig 提供的特定功能或行為，否則建議一律使用 Kubernetes 支援的 Deployment。



參考資料
---------

-----

https://docs.openshift.com/container-platform/4.5/applications/deployments/what-deployments-are.html