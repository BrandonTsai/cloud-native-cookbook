
etcd is the key-value store for OpenShift Container Platform, which persists the entire state of the cluster: its configuration, specifications, and the statuses of the running workloads.
System administrater should Back up cluster’s etcd data regularly and store in a secure location ideally outside the OpenShift



Check Cluster health
--------

you can ssh to master node and run `etcdctl endpoint health` command

```
$ ssh -i id_rsa core@masternode1
$ source /etc/kubernetes/static-pod-resources/etcd-certs/configmaps/etcd-scripts/etcd-common-tools
$ source /etc/kubernetes/static-pod-resources/etcd-certs/configmaps/etcd-scripts/etcd.env
$ etcdctl endpoint health
https://192.168.50.10:2379 is healthy: successfully committed proposal: took = 12.932144ms
https://192.168.50.11:2379 is healthy: successfully committed proposal: took = 14.816544ms
https://192.168.50.12:2379 is healthy: successfully committed proposal: took = 15.857068ms
```

or you can check health via oc command

```
$ oc get etcd -o=jsonpath='{range .items[0].status.conditions[?(@.type=="EtcdMembersAvailable")]}{.message}{"\n"}'
3 members are available
```


Backing up etcd cluster
-------------


```
# /usr/local/bin/cluster-backup.sh /var/home/core/backup/
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

In this example, two files are created in the /var/home/core/backup/ directory on the master host:

- snapshot_<datetimestamp>.db: This file is the etcd snapshot.
- static_kuberesources_<datetimestamp>.tar.gz: This file contains the resources for the static Pods. If etcd encryption is enabled, it also contains the encryption keys for the etcd snapshot.




Fix failed members
---------
etcd cluster achieves high availability by tolerating minor member failures. However, to improve the overall health of the cluster, replace failed members immediately.
When multiple members fail, replace them **one by one**.
Replacing a failed member involves two steps: removing the failed member and adding a new member.
Though etcd keeps unique member IDs internally, it is recommended to use a unique name for each member to avoid human errors.



1. Get the member ID of the failed member:

```
etcdctl  member list -w table
+------------------+---------+-----------------+----------------------------+----------------------------+
|        ID        | STATUS  |      NAME       |         PEER ADDRS         |        CLIENT ADDRS        |
+------------------+---------+-----------------+----------------------------+----------------------------+
|  38cea7b2f31b828 | started | gbvleuaacdopm10 | https://10.248.150.39:2380 | https://10.248.150.39:2379 |
|  e478d9c72b17c7f | started | gbvleuaacdopm11 | https://10.248.150.40:2380 | https://10.248.150.40:2379 |
| c8a14199c8670745 | started | gbvleuaacdopm12 | https://10.248.150.41:2380 | https://10.248.150.41:2379 |
+------------------+---------+-----------------+----------------------------+----------------------------+
```


For example:

```
# etcdctl -C https://10.248.164.14:2379 --ca-file=/etc/etcd/ca.crt  --cert-file=/etc/etcd/peer.crt  --key-file=/etc/etcd/peer.key cluster-health                
member 62716e30614d7108 is healthy: got healthy result from https://10.248.164.14:2379
member e3c8143436eddd4d is unreachable: no available published client urls
member 9f71e4c3be9d2b7c is healthy: got healthy result from https://10.248.164.16:2379
```


2. Remove the failed member:

```
etcdctl member remove c8a14199c8670745
```

3. Delete and recreate the master machine. After this machine is recreated, a new revision is forced and etcd scales up automatically.


4. Verify that all etcd Pods are running properly:

In a terminal that has access to the cluster as a cluster-admin user, run the following command:

```
$ oc get pods -n openshift-etcd | grep etcd
etcd-masternode0                 3/3     Running     0          16d20h
etcd-masternode1                 3/3     Running     0          16d20h
etcd-masternode2                 3/3     Running     0          16d20h
```

Restore from snapshot
----------------

> Notices: The procedure to restore the data MUST be performed on a SINGLE etcd host. Then the etcd cluster Operator handles scaling to the remaining master hosts.

 please refer https://docs.openshift.com/container-platform/4.5/backup_and_restore/disaster_recovery/scenario-2-restoring-cluster-state.html to recover to previous state




Refer
------


 https://medium.com/better-programming/a-closer-look-at-etcd-the-brain-of-a-kubernetes-cluster-788c8ea759a5

OpenShift documentation related to ETCD
- https://docs.openshift.com/container-platform/3.11/admin_guide/assembly_restoring-cluster.html
- https://docs.openshift.com/container-platform/3.11/admin_guide/assembly_replace-etcd-member.html




