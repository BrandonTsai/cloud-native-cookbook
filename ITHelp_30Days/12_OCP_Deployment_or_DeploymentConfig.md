
OpenShift 系列 3: 到底是該用 DeploymentConfig，還是用 Deployment 呀，你搞得我好亂呀！
================================================

Deployment vs DeploymentConfig
-----------------------

OpenShift provide two similar but different methods for fine-grained management over common user applications - Deployment & DeploymentConfig.
They allows you to describe an application’s life cycle, such as which images to use for the app, the number of pods there should be, and the way in which they should be updated. 



| | Deployment | DeploymentConfig |
|--|------------|------------------|
| Building blocks | Using **`ReplicaSet`** | Using **`ReplicationController`** |
| Selector | Supports the new **set-based** selector. <BR> for eg: `environment in (production, qa)`  | only supports equality-based selector. <BR> for eg: `environment = production`  |
| CAP theorem |prefer availability <BR> During a failure it is possible for other masters to act on the same Deployment at the same time. | take consistency over availability <BR> if a node running a deployer Pod goes down, it will not get replaced. The process waits until the node comes back online or is manually deleted.  |
| Automatic rollbacks | Not Support | Can rolling back to the last successfully deployed ReplicationController in case of a failure. |
| Trigger | Not Support | every change in the pod template of a deployment automatically triggers a new rollout. If you do not want new rollouts on pod template changes, pause the deployment |
| Lifecycle hooks | Not Support | |
| Custom strategies | Not Support | Support user-specified Custom deployment strategies |


Both Kubernetes Deployments and OpenShift Container Platform-provided DeploymentConfigs are supported to manage pods. Despite the fact that DeploymentConfig has some extra featuers than Deployment, these features are not very useful in most case, it is recommended to use Deployments unless you need a specific feature or behavior provided by DeploymentConfigs.


DeploymentConfig Trial
-----------------------

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
        image: quay.io/brandon_tsai/testlab:1
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
```

```
$ oc apply -f dc.yml
deploymentconfig.apps.openshift.io/nginx-dc configured

$ oc get pods
NAME             READY   STATUS      RESTARTS   AGE
nginx-dc-1-deploy   0/1     Completed   0          86s
nginx-dc-1-mtrrb    1/1     Running     0          83s
```


Each time a deployment is triggered, whether manually or automatically, a deployer Pod manages the deployment (including scaling down the old ReplicationController, scaling up the new one, and running hooks). The deployment pod remains for an indefinite amount of time after it completes the Deployment in order to retain its logs of the Deployment. 


Scale up the DeploymentConfig to facilitate more load.
----------------------------------------------

You can use the `oc scale` command to manually scale a Deployment or a DeploymentConfig. For example, the following command sets the replicas in the nginx-dc DeploymentConfig to 2.

```
$ oc scale dc nginx-dc --replicas=2
deploymentconfig.apps.openshift.io/nginx-dc scaled

$ oc get pods
NAME                READY   STATUS      RESTARTS   AGE
nginx-dc-1-49cpw    1/1     Running     0          2m54s
nginx-dc-1-deploy   0/1     Completed   0          5m45s
nginx-dc-1-mtrrb    1/1     Running     0          5m42s
```

The number of replicas eventually propagates to the desired and current state of the deployment configured by the DeploymentConfig frontend.



Rollback to earlier Deployment revision 
-----------------------------------------

The `oc rollback` command only works for DeploymentConfig.
I suggest use `oc rollout undo` command instead of `oc rollback`

First, We change the nginx-dc deploymentConfig by editing the image tag.

```
$ oc set image dc/nginx-dc nginx=quay.io/brandon_tsai/testlab:2
```

Wait until the rolling update complete

```
$ oc rollout status -w dc/nginx-dc
```

Rollback to previous version and wait unit the rollback finish.

```
$ oc rollout undo --to-revision=1 dc/nginx/dc
deploymentconfig.apps.openshift.io/nginx-dc deployment #3 rolled back to nginx-dc-1

$ oc rollout status -w dc/nginx-dc
replication controller "nginx-dc-3" successfully rolled out
```

Check the rollout history.

$ oc rollout history dc/nginx-dc
deploymentconfig.apps.openshift.io/nginx-dc 
REVISION	STATUS		CAUSE
1		Complete	config change
2		Complete	config change
3		Complete	config change
```



Conclusion
--------


Reference
----------

https://docs.openshift.com/container-platform/4.5/applications/deployments/what-deployments-are.html