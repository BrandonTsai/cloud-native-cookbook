The applications were typically run standalone in a VM and use all of the resource.

The operators and developers need to choose the size of the VM for running the application.

But in Openshift, the pods/containers can be running on any machine which requires sharing the resources with others.

That is where the QoS (Quality of Service Classes) and Resource Quotas comes in. 



Resource Request and Limits
When you create a pod for your application, you can set requests and limits for CPU and memory for every container inside.

Properly setting these values is the way to instruct Kubernetes how to manage your applications.



For example,

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


Requests: The values is used for scheduling. It is the minimum amount of resources a container needs to run. The Pods will remain "Pending" state if no node has enough resources for the request. 

Limits: The maximum amount of this resource that the node will allow the containers to use. 

If a container attempts to exceed the specified limit, the system will throttle the container
If the container exceeds the specified memory limit, it will be terminated and potentially restarted dependent upon the container restart policy.
Quality of Service Classes (QoS)


A node is overcommitted when it has a pod scheduled that makes no request, or when the sum of limits across all pods on that node exceeds available machine capacity.
In an overcommitted environment, it is possible that the pods on the node will attempt to use more compute resource than is available at any given point in time.

When this occurs, the node must give priority to one container over another. Containers that have the lowest priority are terminated/throttle first. The facility used to make this decision is referred to as a Quality of Service (QoS) Class. 



Priority	Class Name	Description
1 (highest)	Guaranteed	If limits and optionally requests are set (not equal to 0) for all resources and they are equal.
2	Burstable	If requests and optionally limits are set (not equal to 0) for all resources, and they are not equal
3 (lowest)	BestEffort	If requests and limits are not set for any of the resources


Therefore, if developer do not declare CPU/Memory requests and limits, the container will be terminated first. We should Protect the critical pods in production projects by setting limits value so they are classified as Guaranted. BestEffort or Burstable pods should be used in developing projects only.

Quota 
Administrator can set the Project Quota to restrict the resource consumption. 
This has additional effect; If you set a Memory request in the quota, then all pods need to set a Memory request in their definition.
The new pod will not be scheduled and will remain in a pending state if it try to allocate resource more than the quota restriction.



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


Multi-Project Quota
A multi-project quota, defined by a ClusterResourceQuota object, allows quotas to be shared across multiple projects.



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


LimitRanges
A LimitRange is a policy to constrain resource by Pod or Container in a namespace. it can:

Set default request/limit for compute resources in a namespace and automatically inject them to Containers at runtime.
Enforce minimum and maximum compute resources usage per Pod or Container in a namespace.
Enforce minimum and maximum storage request per PersistentVolumeClaim in a namespace.
Enforce a ratio between request and limit for a resource in a namespace.


Following issue may happen after applying default value by LimitRange. Developers should be aware of the CPU/memory limits settings for each container.

The default cpu/memory request value may be too high for some namespaces
Pods may be in "Pending" state easily, and developers do not know why.
The default cpu/memory request value may be too low for other namespaces.
Pods may crash frequently, and developers do not know why.


GTS Openshift Project Quota, Multi-Project Quota Plan


Resource

Default Value for each container

CPU	Request	10m
Limit	400m

Memory	Request	100 Mi
Limit	800 Mi


Resource

Small (default)

Medium

Large

XLarge
CPU



Request

4

8

16

32
Limits

4	8	16	32
Memory	Request	8Gi	16Gi	32Gi	64Gi
Limits	8Gi	16Gi	32Gi	64Gi
Max Pods # with default LimitRange Value 	10	20	40	80


We could customize the Quota for those projects that require high memory/cpu resources.



How do developer know how many resource they need for each container?


Run "docker stats" on local development to understand CPU/Memory requests per container
They need a VM to run Docker container and test it locally.
Observe resource usage in Splunk and do performance test in DEV/SIT environment to get the value of CPU/Memory limits


Current Known issues.
Developers need to specify the cpu/memory request/limits for "EACH" container. But they only know how much resources they need in total, not for each process/container. As a result, they can not set the value appropriately.
If the value is too large, their pods will be in "Pending" state, and we need to keep scale out out clusters even they did not use that much resource.
If the value is too small, their pods will keep restart.
Sometimes they just launch deployment by "oc new-app" command, which does not specify the cpu/memory limits.
they do not want to implement openshift yaml for a testing application.
Set default cpu/memory request/limits by LimitRange is a very bad idea, because
each container require different amount of resource.
developers won't notice that there is a limitation for CPU/Memory resource.


Steps of Introducing Quota/Multi-Project Quota to Developes






Reference:
https://docs.openshift.com/container-platform/3.11/admin_guide/out_of_resource_handling.html
https://docs.openshift.com/container-platform/3.11/admin_guide/overcommit.html#qos-classes
