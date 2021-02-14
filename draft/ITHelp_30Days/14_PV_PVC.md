儲存者聯盟： Storage Class, PV, PVC
===========

When do we need Persistent Volume（PV）?
---------------------
- massive data generate by app you want to keep when pod recreated
- large files that are not suitable to keep in configMap/Secrets because configMap and Secrets are keeped in etcd Cluster


Understanding Storage
----------

Ideally, a developer deploying their apps on Kubernetes should never have to know what kind of storage technology is used underneath, the same way they don’t have to know what type of physical servers are being used to run their pods.
When a developer needs a certain amount of persistent storage for their application, they can request it from Kubernetes, the same way they can request CPU, memory, and other resources when creating a pod. The system administrator can configure the cluster so it can give the apps what they request.

To enable apps to request storage in a Kubernetes cluster without having to deal with infrastructure specifics,
The process and components are described as the following photo.


![](https://i.stack.imgur.com/k3WkN.png)


### PersistentVolume (PV)
A PersistentVolume (PV) is a piece of storage in the cluster that has been manually provisioned by an administrator, or dynamically provisioned by Kubernetes using a StorageClass. A PersistentVolumeClaim (PVC) is a request for storage by a user that can be fulfilled by a PV. PersistentVolumes and PersistentVolumeClaims are independent from Pod lifecycles and preserve data through restarting, rescheduling, and even deleting Pods.
PV resources are not scoped to any single project; they can be shared across the entire OpenShift Container Platform cluster and claimed from any project.

### PVC
PVCs are specific to a project, and are created and used by developers as a means to use a PV. 
After a PV is bound to a PVC, that PV can not then be bound to additional PVCs. This has the effect of scoping a bound PV to a single namespace, that of the binding project.


### StorageClass
The Storage Object in Use Protection feature ensures that PVCs in active use by a Pod and PVs that are bound to PVCs are not removed from the system, as this can result in data loss.
Each StorageClass contains the fields provisioner, parameters, and reclaimPolicy, which are used when a PersistentVolume belonging to the class needs to be dynamically provisioned.
Administrators can specify a default StorageClass just for PVCs.

The reclaim policy of a PersistentVolume tells the cluster what to do with the volume after it is released. A volume’s reclaim policy can be Retain, or Delete.

- Retain: allows manual reclamation of the resource for those volume plug-ins that support it.
- Delete: deletes both the PersistentVolume object from OpenShift Container Platform and the associated storage asset in external infrastructure, such as AWS EBS or VMware vSphere.


### Provisioner
A provisioner that determines what volume plugin is used for provisioning actual storage and PVs.
According to the storage service provider, 3 different access modes are supported, as following

| Access Mode| CLI abbreviation | Description |
| ReadWriteOnce | RWO | The volume can be mounted as read-write by a single node. |
| ReadOnlyMany | ROX | The volume can be mounted as read-only by many nodes.|
| ReadWriteMany | RWX | The volume can be mounted as read-write by many nodes. |




Example 
---------

In our Client, we use "VMware vSphere volumes" by default.


```YAML
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  labels:
    app: postgres-db
  name: postgres-db
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
  storageClassName: vsphere-standard
```

```YAML
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  annotations:
  labels:
    app: postgres-db
  name: postgres-db
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres-db
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: postgres-db
    spec:
      containers:
      - env:
        - name: POSTGRESQL_USER
          value: reference_data_user
        - name: POSTGRESQL_DATABASE
          value: payments-testdb
        - name: POSTGRESQL_PASSWORD
          value: testing
        - name: POSTGRESQL_ADMIN_PASSWORD
          value: testing
        image: quay-ap.windmill.local/gts-base-images/postgresql-96-rhel7:1-52-release
        imagePullPolicy: IfNotPresent
        name: postgres-db
        ports:
        - containerPort: 5432
          name: postgres-db
          protocol: TCP
        volumeMounts:
        - mountPath: /var/lib/postgresql/data
          name: postgresql
      restartPolicy: Always
      volumes:
      - name: postgresql
        persistentVolumeClaim:
          claimName: postgres-db

```


apply this yaml file

```
$ oc apply -f /tmp/pvc.yaml 
persistentvolumeclaim/postgres-db unchanged
deployment.extensions/postgres-db created


$ oc get pods
NAME                           READY     STATUS    RESTARTS   AGE
postgres-db-754f6dc5c7-ck8nf   1/1       Running   0          9m

$ oc get pvc
NAME          STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS       AGE
postgres-db   Bound     pvc-e9c12dab-fed2-11ea-9a08-0050569e649e   5Gi        RWO            vsphere-standard   9m

$ oc get pv pvc-e9c12dab-fed2-11ea-9a08-0050569e649e
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS    CLAIM                     STORAGECLASS       REASON    AGE
pvc-e9c12dab-fed2-11ea-9a08-0050569e649e   5Gi        RWO            Retain           Bound     gts-lab-dev/postgres-db   vsphere-standard             10m
```


What Happen if we scale up deployment to 2 pods?
--------------



```
$ oc scale deploy postgres-db --replicas=2 
deployment.extensions/postgres-db scaled

$ oc get pods
NAME                           READY     STATUS              RESTARTS   AGE
postgres-db-754f6dc5c7-7khnc   0/1       ContainerCreating   0          8m
postgres-db-754f6dc5c7-ck8nf   1/1       Running             0          42m
```

That is because our storage Class is vSphere, it does not support ReadWriteMany or ReadOnlyMany mode.

```
Events:
  Type     Reason              Age              From                                     Message
  ----     ------              ----             ----                                     -------
  Normal   Scheduled           8m               default-scheduler                        Successfully assigned gts-lab-dev/postgres-db-754f6dc5c7-7khnc to sgvlapaacdopa02.windmill.local
  Warning  FailedAttachVolume  8m               attachdetach-controller                  Multi-Attach error for volume "pvc-e9c12dab-fed2-11ea-9a08-0050569e649e" Volume is already used by pod(s) postgres-db-754f6dc5c7-ck8nf
  Warning  FailedMount         1m (x3 over 6m)  kubelet, sgvlapaacdopa02.windmill.local  Unable to mount volumes for pod "postgres-db-754f6dc5c7-7khnc_gts-lab-dev(c38af83c-fed7-11ea-b1e4-0050569e6f56)": timeout expired waiting for volumes to attach or mount for pod "gts-lab-dev"/"postgres-db-754f6dc5c7-7khnc". list of unmounted volumes=[postgresql]. list of unattached volumes=[postgresql default-token-5bblj
```



Clean up
-------------

```
$ oc delete deploy postgres-db
deployment.extensions "postgres-db" deleted

$ oc delete pvc postgres-db
persistentvolumeclaim "postgres-db" deleted

$ oc get pvc               
No resources found.

$ oc get pv pvc-e9c12dab-fed2-11ea-9a08-0050569e649e
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS     CLAIM                     STORAGECLASS       REASON    AGE
pvc-e9c12dab-fed2-11ea-9a08-0050569e649e   5Gi        RWO            Retain           Released   gts-lab-dev/postgres-db   vsphere-standard             53m
```

The PV still exist, because our reclaimPolicy of storageclass vsphere-standard is Retain
we need to delete this PV manually.



Can I resize the PV as Data growth with time?
--------------------------------------------

It depends on your storage provider, and you also need some extra work for this feature. please refer  https://docs.openshift.com/container-platform/4.5/storage/expanding-persistent-volumes.html


Conclusion
----------

There is no different between Kubernetes and Openshift on Storage. 
Just be noticed that different Storage Provider has different features and backup/recovery strategy. 
It is very important to do a fully test before integrate it into your production environment, such as stress test, performace test during backup period, and make sure you know how to do recovery when data lost. 


Reference
---------

- buildin provisioner list: refer https://docs.openshift.com/container-platform/4.5/storage/understanding-persistent-storage.html#types-of-persistent-volumes_understanding-persistent-storage
- access mode: refer https://docs.openshift.com/container-platform/4.5/storage/understanding-persistent-storage.html#pv-access-modes_understanding-persistent-storage





