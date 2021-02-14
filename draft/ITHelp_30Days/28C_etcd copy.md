
etcd 是 OpenShift平台的鍵值存儲資料庫（key-value store），可儲存整個系統每個資源的狀態，譬如配置，規格以及運行中的工作負載的狀態。系統管理員應定期備份 etcd 數據，並最好備份到在 OpenShift 外部的安全位置。

本篇文章將介紹如何備份 OpenShift 4.5 的 etcd 數據。

確認 etcd Cluster 是健康的
--------

-----

方法一： ssh 到其中一個 master 節點並執行 `etcdctl endpoint health` 指令

```
$ ssh -i id_rsa core@masternode1
$ source /etc/kubernetes/static-pod-resources/etcd-certs/configmaps/etcd-scripts/etcd-common-tools
$ source /etc/kubernetes/static-pod-resources/etcd-certs/configmaps/etcd-scripts/etcd.env
$ etcdctl endpoint health
https://192.168.50.10:2379 is healthy: successfully committed proposal: took = 12.932144ms
https://192.168.50.11:2379 is healthy: successfully committed proposal: took = 14.816544ms
https://192.168.50.12:2379 is healthy: successfully committed proposal: took = 15.857068ms
```

方法二： 透過 oc 指令

```
$ oc get etcd -o=jsonpath='{range .items[0].status.conditions[?(@.type=="EtcdMembersAvailable")]}{.message}{"\n"}'
3 members are available
```


備份 etcd 數據。
-------------


-----

在任何一個健康的 master 節點執行以下指令

```
＄ /usr/local/bin/cluster-backup.sh /var/home/core/backup/
512a3e830ede6af4472474ae1ab90ac7c5fb8c9e60b20b96bbb90a30f8c8a97c
etcdctl version: 3.3.18
API version: 3.3
found latest kube-apiserver-pod: /etc/kubernetes/static-pod-resources/kube-apiserver-pod-63
found latest kube-controller-manager-pod: /etc/kubernetes/static-pod-resources/kube-controller-manager-pod-39
found latest kube-scheduler-pod: /etc/kubernetes/static-pod-resources/kube-scheduler-pod-37
found latest etcd-pod: /etc/kubernetes/static-pod-resources/etcd-pod-21
Snapshot saved at /var/home/core/backup//snapshot_2020-10-13_062754.db
snapshot db and kube resources are successfully saved to /var/home/core/backup/
```

上面指令會產生兩個檔案在 /var/home/core/backup/ 資料夾

- snapshot_<datetimestamp>.db: etcd 的主要 snapshot.
- static_kuberesources_<datetimestamp>.tar.gz: 這個檔案包含了 master 節點上 static Pods 的資訊。

執行完後只要將該這兩個檔案上傳到 Nexus 或其他 OpenShift 外部的位置保存即可。


如何修復 failed 的 etcd 成員
---------

-----

etcd 可容忍一小部分成員故障來實現高可用性。 但是，為了改善 etcd 的整體運行狀況，請立即更換發生故障的成員。如果多個成員失敗，請“一個接一個”的替換它們，切勿同時替換所有故障的成員。替換失敗的成員主要涉及兩個步驟：刪除失敗的成員和添加新成員。

1. 取得故障成員的 ID :

```
etcdctl  member list -w table
+------------------+---------+-----------------+----------------------------+----------------------------+
|        ID        | STATUS  |      NAME       |         PEER ADDRS         |        CLIENT ADDRS        |
+------------------+---------+-----------------+----------------------------+----------------------------+
|  38cea7b2f31b828 | started | masternode1 | https://192.168.50.10:2380 | https://192.168.50.10:2379 |
|  e478d9c72b17c7f | started | masternode2 | https://192.168.50.11:2380 | https://192.168.50.11:2379 |
| c8a14199c8670745 | unreachable | masternode3 | https://192.168.50.12:2380 | https://192.168.50.12:2379 |
+------------------+---------+-----------------+----------------------------+----------------------------+
```


2. 移除故障成員:

```
etcdctl member remove c8a14199c8670745
```

3. 刪除並重新創建 master 節點. 當新的 master 節點建立之後, etcd 會自動 scales up 。


4. 確定所有瘩 etcd Pods 狀態都是 Running:

```
$ oc get pods -n openshift-etcd | grep etcd
etcd-masternode0                 3/3     Running     0          16d20h
etcd-masternode1                 3/3     Running     0          16d20h
etcd-masternode2                 3/3     Running     0          16d20h
```

Restore from snapshot
----------------

-----

請參考[官方文件](https://docs.openshift.com/container-platform/4.5/backup_and_restore/disaster_recovery/scenario-2-restoring-cluster-state.html) 來透過之前備份的 snapshot 復原 etcd 到之前的狀態。


> 注意: 這個步驟必須先移除其他 etcd 成員，確定只有一個主要的節點有正常的運行的 etcd 。 復原到之前的狀態漏，再透過 etcd cluster Operator 將 etcd 擴充到剩餘的 master 節點.


結論
------

-----

如果有要對 master 節點做安全性更新而需要重新啟動 VM，務必在那之前先做好備份。




參考資料
------

- https://medium.com/better-programming/a-closer-look-at-etcd-the-brain-of-a-kubernetes-cluster-788c8ea759a5
- https://docs.openshift.com/container-platform/4.5/backup_and_restore/disaster_recovery/scenario-2-restoring-cluster-state.html





