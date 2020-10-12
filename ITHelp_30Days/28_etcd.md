What is etcd?
Refer: https://medium.com/better-programming/a-closer-look-at-etcd-the-brain-of-a-kubernetes-cluster-788c8ea759a5



OpenShift documentation related to ETCD
- https://docs.openshift.com/container-platform/3.11/admin_guide/assembly_restoring-cluster.html
- https://docs.openshift.com/container-platform/3.11/admin_guide/assembly_replace-etcd-member.html



Backing up etcd cluster
Manually backup the etcd data on a healthy node:

```
ETCDCTL_API=3
etcdctl --cert /etc/etcd/peer.crt --key /etc/etcd/peer.key --cacert /etc/etcd/ca.crt --endpoints 'gbvleqaacopm01p.windmill.local:2379,https://gbvleqaacopm02p.windmill.local:2379,https://gbvleqaacopm03p.windmill.local:2379' snapshot save  /tmp/snapshot_$(date +%Y%m%d).db
```




Check Cluster Health
This command will check the cluster health status

```
etcdctl -C https://10.248.150.19:2379 \
  --ca-file=/etc/etcd/ca.crt     \
  --cert-file=/etc/etcd/peer.crt     \
  --key-file=/etc/etcd/peer.key cluster-health
```

List members of the ETCD cluster
This command will list the members of the etcd cluster

```
etcdctl --cert-file=/etc/etcd/peer.crt \
    --key-file=/etc/etcd/peer.key \
    --ca-file=/etc/etcd/ca.crt \
    --peers="https://172.18.1.18:2379,https://172.18.9.202:2379,https://172.18.0.75:2379" \
    member list
```



Fix failed members
---------
etcd cluster achieves high availability by tolerating minor member failures. However, to improve the overall health of the cluster, replace failed members immediately.
When multiple members fail, replace them **one by one**.
Replacing a failed member involves two steps: removing the failed member and adding a new member.
Though etcd keeps unique member IDs internally, it is recommended to use a unique name for each member to avoid human errors.



1. Get the member ID of the failed member:

```
etcdctl -C https://<surviving host IP>:2379 --ca-file=/etc/etcd/ca.crt --cert-file=/etc/etcd/peer.crt --key-file=/etc/etcd/peer.key cluster-health
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
etcdctl -C https://10.248.164.14:2379 --ca-file=/etc/etcd/ca.crt --cert-file=/etc/etcd/peer.crt --key-file=/etc/etcd/peer.key member remove e3c8143436eddd4d
```

3. On the failed etcd member

Stop the etcd service by removing the etcd pod definition:

```
$ mkdir -p /etc/origin/node/pods-stopped
$ mv /etc/origin/node/pods/etcd.yaml /etc/origin/node/pods-stopped/
$ sudo service docker restart
```

Remove the contents of theetcddirectory
```
# mkdir -p /var/lib/etcd-old-data
# mv /var/lib/etcd/* /var/lib/etcd-old-data/
```

4. Add the member back to the cluster:
```
# etcdctl -C https://${CURRENT_ETCD_HOST}:2379 --ca-file=/etc/etcd/ca.crt --cert-file=/etc/etcd/peer.crt --key-file=/etc/etcd/peer.key member add ${NEW_ETCD_HOSTNAME} https://${NEW_ETCD_IP}:2380
```

For example
```
# etcdctl -C https://10.248.164.14:2379 --ca-file=/etc/etcd/ca.crt --cert-file=/etc/etcd/peer.crt --key-file=/etc/etcd/peer.key member add gbvleqaacopm02p.windmill.local https://10.248.164.15:2380
Added member named gbvleqaacopm02p.windmill.local with ID 4e1db163a21d7651 to cluster
 
ETCD_NAME="gbvleqaacopm02p.windmill.local"
ETCD_INITIAL_CLUSTER="gbvleqaacopm01p.windmill.local=https://10.248.164.14:2380,gbvleqaacopm02p.windmill.local=https://10.248.164.15:2380,gbvleqaacopm03p.windmill.local=https://10.248.164.16:2380"
ETCD_INITIAL_CLUSTER_STATE="existing"
```

5. On the failed etcd member, 

Update the /etc/etcd/etcd.conf file, replace the following values with the values generated in the previous step:

- ETCD_NAME
- ETCD_INITIAL_CLUSTER
- ETCD_INITIAL_CLUSTER_STATE


Start  etcd service
```
# mv /etc/origin/node/pods-stopped/etcd.yaml /etc/origin/node/pods/
# sudo service docker restart
```

Wait until etcd container is running. and check the logs of the etcd container.

```
/usr/local/bin/master-logs etcd etcd
```


6. Check the etcd cluster is healthy
```
# etcdctl -C https://10.248.164.14:2379 --ca-file=/etc/etcd/ca.crt --cert-file=/etc/etcd/peer.crt --key-file=/etc/etcd/peer.key cluster-health
member 62716e30614d7108 is healthy: got healthy result from https://10.248.164.14:2379
member 650b55864907513f is healthy: got healthy result from https://10.248.164.15:2379
member 9f71e4c3be9d2b7c is healthy: got healthy result from https://10.248.164.16:2379
cluster is healthy
```


Restore from snapshot
----------------

> Notices: The procedure to restore the data MUST be performed on a SINGLE etcd host. You can then add the rest of the nodes to the cluster.


Please refer: https://docs.openshift.com/container-platform/3.11/admin_guide/assembly_restoring-cluster.html#restoring-etcd-v3-snapshot



Update ETCD_INITIAL_CLUSTER_STATE
------------------


On the master nodes do the following:

```
#Edit etcd.conf file
sudo vi /etc/etcd/etcd.conf
ETCD_INITIAL_CLUSTER_STATE=new
#Change to
ETCD_INITIAL_CLUSTER_STATE=existing
# Check that this flag is blank
ETCD_INITIAL_CLUSTER=
#Find PID of etcd container
sudo docker ps | grep k8s_etcd_master-etcd-
#Kill etcd container
sudo docker kill [PID]
#Check that etcd got restarted
sudo docker ps | grep k8s_etcd_master-etcd-
 
#Verify etcd endpoint on local node is up and cluster is healthy.
etcdctl -C https://[NODE-IP]:2379 \
  --ca-file=/etc/etcd/ca.crt     \
  --cert-file=/etc/etcd/peer.crt     \
  --key-file=/etc/etcd/peer.key cluster-health

````




