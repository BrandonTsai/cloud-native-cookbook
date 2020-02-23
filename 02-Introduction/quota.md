

The applications were typically run standalone in a VM and use all of the resource.
The operators and developers need to choose the size of the VM for running the application.
But in Openshift, the pods/containers can be running on any machine which requires sharing the resources with others.
That is where the QoS (Quality of Service Classes) and Resource Quota comes in. 



## Resource Request and Limits
When you create a pod for your application, you can set requests and limits for CPU and memory for every container inside.
Properly setting these values is the way to instruct Kubernetes how to manage your applications.



For example,
```yaml
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
        cpu: 200m
        memory: 400Mi
        ephemeral-storage: 2Gi
```

**Requests**: The values is used for scheduling. It is the minimum amount of resources a container needs to run. The Pods will remain "Pending" state if no node has enough resources for the request. 

**Limits**: The maximum amount of this resource that the node will allow the containers to use. 

- If a container attempts to exceed the specified limit, the system will throttle the container
- If the container exceeds the specified memory limit, it will be terminated and potentially restarted dependent upon the container restart policy.


## Quality of Service Classes (QoS)


A node is overcommitted when it has a pod scheduled that makes no request, or when the sum of limits across all pods on that node exceeds available machine capacity.
In an overcommitted environment, it is possible that the pods on the node will attempt to use more compute resource than is available at any given point in time.

When this occurs, the node must give priority to one container over another. Containers that have the lowest priority are terminated/throttle first. The facility used to make this decision is referred to as a Quality of Service (QoS) Class. 


| Priority |	Class Name | Description |
| -------- | ----------- | ----------- |
| 1 (highest)	| Guaranteed | If limits and optionally requests are set (not equal to 0) for all resources and they are equal. |
| 2           | Burstable  | If requests and optionally limits are set (not equal to 0) for all resources, and they are not equal |
| 3 (lowest)  | BestEffort | If requests and limits are not set for any of the resources |

Therefore, if developer do not declare CPU/Memory requests and limits, the container will be terminated first. We should Protect the critical pods in production projects by setting limits value so they are classified as Guaranted. BestEffort or Burstable pods should be used in developing projects only.


## Project Quota and Limit Ranges:

Administrator can set the Project Quota to restrict the resource consumption. 
This has additional effect; If you set a Memory request in the quota, then all pods need to set a Memory request in their definition.
The new pod will not be scheduled and will remain in a pending state if it try to allocate resource more than the quota restriction.


A limit range is a policy to constrain resource by Pod or Container in a namespace. it can:

- Set default request/limit for compute resources in a namespace and automatically inject them to Containers at runtime.
- Enforce minimum and maximum compute resources usage per Pod or Container in a namespace.
- Enforce minimum and maximum storage request per PersistentVolumeClaim in a namespace.
- Enforce a ratio between request and limit for a resource in a namespace.




## What should we monitor for managing cluster resource?


### Node Status

Make sure all nodes are in "Ready" state



### Percentage of resource (CPU/Memory/Disk) allocated from total available resource in the cluster

A good warning threshold would be (n-1)/n*100, where n is the number of nodes.

Over this threshold, you wouldn't be able to reallocate your workloads in the rest of the nodes.



### Percentage of Resource (CPU/Memory) Usage in the node

The OS Kernel invokes OOMKiller when Memory usage come under pressure in the node, even all of the containers being under their limits.

CPU Pressure will restrain processes and affect performance.


A warning threshold to notify administrator that this node may have performance issue or reach the "Eviction Policies".

  - Check the "Eviction Policies" setting. Make sure get alerts before reach eviction-hard policy. 


### The CPU and Memory Usage/Request/Limit vs Capacity in the node

Following warning thresholds to notify administrator that this node may not able to allocated new pods 

- Less than 200m CPU can be allocate to CPU Request/Limits
- Less than 200Mi Memory can be allocate to Memory Request/Limits


if n-1 nodes can not allocate new pods, then it is time to scale up or check the CPU/Memory requests are too high or not.



### Disk Space in the node

If the node runs out of disk, it will try to free docker space with a fair chance of pod eviction



### Memory and CPU usage per container

Because Openshift limits are per container, not per pod. So don't waste time on them if you have performance issue.

Ideally, the container should usage exactly the amount of resource it requested.
If your usage is much lower than your request, you are waste money and it would be too hard to allocate new pods.
if it is higher, you are risking performance issue in the nodes



## How do developer know how many resource they need for each container?


- Run "docker stats" on local development to understand CPU/Memory requests per container
- Do performance test in UAT environment to get the value of CPU/Memory limits


## Reference:
https://docs.openshift.com/container-platform/3.11/admin_guide/out_of_resource_handling.html
https://docs.openshift.com/container-platform/3.11/admin_guide/overcommit.html#qos-classes

