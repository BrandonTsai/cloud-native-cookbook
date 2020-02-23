

The applications were typically run standalone in a VM and use all of the resources.
The operators and developers need to choose the size of the VM for running the application.
But in Kubernetes, the pods/containers can be running on any machine which requires sharing the resources with others.
That is where the QoS (Quality of Service Classes) and Resource Quota comes in.



## Resource Request and Limits
When you create a pod for your application, you can set requests and limits for CPU and memory for every container inside.
Properly setting these values is the way to instruct Kubernetes on how to manage your applications.



For example,
```yaml
spec:
  containers:
  - image: k8s/hello-k8s
    name: hello-k8s
    resources:
      requests:
        cpu: 100m
        memory: 200Mi
      limits:
        cpu: 200m
        memory: 400Mi
```

**Requests**: The values are used for scheduling. It is the minimum amount of resources a container needs to run. The Pods will remain "Pending" state if no node has enough resources for the request.

**Limits**: The maximum amount of this resource that the node will allow the containers to use.

- If a container attempts to exceed the specified limit, the system will throttle the container
- If the container exceeds the specified memory limit, it will be terminated and potentially restarted dependent upon the container restart policy.


## Quality of Service Classes (QoS)


A node is overcommitted when it has a pod scheduled that makes no request, or when the sum of limits across all pods on that node exceeds available machine capacity.
In an overcommitted environment, the pods on the node may attempt to use more compute resources than available at any given point in time.

When this occurs, the node must give priority to one container over another. Containers that have the lowest priority are terminated/throttle first. The facility used to make this decision is referred to as a Quality of Service (QoS) Class.


| Priority |	Class Name | Description |
| -------- | ----------- | ----------- |
| 1 (highest)	| Guaranteed | If limits and optionally requests are set (not equal to 0) for all resources and they are equal. |
| 2           | Burstable  | If requests and optionally limits are set (not equal to 0) for all resources, and they are not equal |
| 3 (lowest)  | BestEffort | If requests and limits are not set for any of the resources |

Therefore, if the developer does not declare CPU/Memory requests and limits, the container will be terminated first. We should Protect the critical pods in production projects by setting limits value so they are classified as *Guaranted*. *BestEffort* or *Burstable* pods should be used in developing projects only.


## Project Quota and Limit Ranges:

The administrator can set the Project Quota to restrict resource consumption.
This has an additional effect; If you set a Memory request in the quota, then all pods need to set a Memory request in their definition.
The new pod will not be scheduled and will remain pending if it tries to allocate resources more than the quota restriction.


A limit range is a policy to constrain resources by Pod or Container in a namespace. it can:

- Set default request/limit for computing resources in a namespace and automatically inject them to Containers at runtime.
- Enforce minimum and maximum resources usage per Pod or Container in a namespace.
- Enforce minimum and maximum storage requests per PersistentVolumeClaim in a namespace.
- Enforce a ratio between request and limit for a resource in a namespace.




## What should we monitor for managing cluster resources?


### Node Status

Make sure all nodes are in "Ready" state



### Percentage of resource (CPU/Memory/Disk) allocated from the total available resource in the cluster

A good warning threshold would be (n-1)/n*100, where n is the number of nodes.

Over this threshold, you wouldn't be able to reallocate your workloads in the rest of the nodes.



### Percentage of Resource (CPU/Memory) Usage in the node

The OS Kernel invokes OOMKiller when Memory usage comes under pressure in the node, even all of the containers being under their limits.

CPU Pressure will restrain processes and affect performance.


A warning threshold to notify the administrator that this node may have the issues or reach the "Eviction Policies".

  - Check the "Eviction Policies" setting. Make sure alerts have triggered before reaching the eviction-hard thresholds.


### The CPU and Memory Usage/Request/Limit vs Capacity in the node

Add the following warning thresholds to notify the administrator that this node may not able to allocate new pods.

- Less than 200m CPU can be allocated to CPU Request/Limits
- Less than 200Mi Memory can be allocated to Memory Request/Limits


If n-1 nodes can not allocate new pods, then it is time to scale up or check the CPU/Memory requests are too high or not.



### Disk Space in the node

If the node runs out of disk, it will try to free docker space with a fair chance of pod eviction



### Memory and CPU usage per container

Because Kubernetes limits are per container, not per pod. So it is not necessary to monitor the usage per pod.

Ideally, the container should use exactly the amount of resources that it requested.
If your usage is much lower than your request, you are waste money and it would be too hard to allocate new pods.
if it is higher, you are risking the performance issue in the nodes



## How do a developer know how many resources they need for each container?

- Run "docker stats" on local development to understand CPU/Memory requests per container
- Do the  test in UAT environment to get the value of CPU/Memory limits



## Conclution

It is important to make sure requests and limits are declared and tested before deploying to production namespace. Cluster admin can set up namespace quota to enforce all of the workloads in the namespace to have a request and limit in every container. A good configuration of requests and limits will allow the application much stable.

Appropriate monitoring and alerts will help the cluster admin to reduce the waste of the cluster resources and avoid the performance issue. Ask us today if you need help to monitor your Kubernetes system.


## Reference:
- https://docs.openshift.com/container-platform/3.11/admin_guide/out_of_resource_handling.html
- https://docs.openshift.com/container-platform/3.11/admin_guide/overcommit.html#qos-classes
