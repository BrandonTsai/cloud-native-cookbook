---
title: "#7 OpenShift: Operator & OperatorHub"
author: Brandon Tsai
---


What is Operator?
-----------------


The conventional wisdom of Kubernetes’ earlier days was that it was very good at managing stateless apps. 
But for stateful applications such as databases, it wasn’t such an open-and-shut case: These apps required more hand-holding.

“Adding or removing instances may require preparation and/or post-provisioning steps – for instance, changes to its internal configuration, communication with a clustering mechanism, interaction with external systems like DNS, and so forth,” Thompson explains. “Historically, this often required manual intervention, increasing the DevOps burden and increasing the likelihood of error. Perhaps most importantly, it obviates one of Kubernetes’ main selling points: automation.”

That’s a big problem. Fortunately, the solution emerged back in 2016, when coreOS introduced Operators to extend Kubernetes’ capabilities to stateful applications. (Red Hat acquired coreOS in January 2018, expanding the capabilities of the OpenShift container platform.) Operators became even more powerful with the launch of the Operator Framework for building and managing Kubernetes native applications (Operators by another name) in March 2018. Amidst the broader excitement about Kubernetes, the importance of Operators seems understated – when in fact it would be hard to overstate their importance.


The idea is that when you have an application, like a database like Postgres or Cassandra…any complex application needs a lot of domain-specific knowledge,” Sebastien Pahl, director of engineering, Red Hat OpenShift, explained in a Kubecon 2018 presentation. You can use Kubernetes tools to try and solve some of those problems



Operator Framework
------------------
Operator Framework is an open source toolkit designed to manage Kubernetes Operators, in a more effective, automated, and scalable way.  It is not just about writing code; testing, delivering, and updating Operators is just as important. The Operator Framework components consist of open source tools to tackle these problems:

Operator SDK
The Operator SDK assists Operator authors in bootstrapping, building, testing, and packaging their own Operator based on their expertise without requiring knowledge of Kubernetes API complexities.

Operator Lifecycle Manager
Operator Lifecycle Manager (OLM) controls the installation, upgrade, and role-based access control (RBAC) of Operators in a cluster. Deployed by default in OpenShift Container Platform 4.5.

Operator Registry
The Operator Registry stores ClusterServiceVersions (CSVs) and Custom Resource Definitions (CRDs) for creation in a cluster and stores Operator metadata about packages and channels. It runs in a Kubernetes or OpenShift cluster to provide this Operator catalog data to OLM.

OperatorHub
OperatorHub is a web console for cluster administrators to discover and select Operators to install on their cluster. It is deployed by default in OpenShift Container Platform.

Operator Metering
Operator Metering collects operational metrics about Operators on the cluster for Day 2 management and aggregating usage metrics.

These tools are designed to be composable, so you can use any that are useful to you.



Deploy first Operator
---------------------


When Login as kubeadmin, you can see OperatorHub as following

![](images/Operator/1.png)

Search Web Terminal Operator

![](images/Operator/2.png)

Install Web Terminal Operator with default settings

![](images/Operator/3.png)

Check Operator status

![](images/Operator/4.png)

Make sure user can connect into pod via Terminal from Web UI.

![](images/Operator/5.png)

Reference
--------

It is very useful for Operator to deploy Production level, High Availaity Service for complexity application such as Storage service and Database.

I will introduce some useful Operator for managing OpenShift Cluster.

https://enterprisersproject.com/article/2019/2/kubernetes-operators-plain-english?page=1


Install Operator via command 
----------


1. Search Operator you want to install

```
# Search operator
$ oc get packagemanifests -n openshift-marketplace | grep "security"
container-security-operator                 Red Hat Operators     22d
```

2. check operator information to get `channel` and `CSV` information

```
# check operator info
$ oc describe packagemanifests container-security-operator -n openshift-marketplace

```


3. Create Subscription resource

Implement Subscription YAML file

```
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: container-security-operator
  namespace: openshift-operators
spec:
  channel: quay-v3.3
  installPlanApproval: Automatic
  name: container-security-operator
  source: redhat-operators
  sourceNamespace: openshift-marketplace
  startingCSV: container-security-operator.v3.3.1
```

Apply the YAML file

```
$ oc apply -f sub.yml 
subscription.operators.coreos.com/container-security-operator created
```


4. Check the pod is running.

```
$ oc get pods
NAME                                           READY   STATUS              RESTARTS   AGE
container-security-operator-5d8c9c64d6-cf99d   1/1     Running             0          16s
devworkspace-controller-84f877b-4d9x9          1/1     Running             0          2d3h
devworkspace-webhook-server-7d6645dc7b-6nkd9   1/1     Running             1          2d2h
```


Test container-security-operator from Web UI
------------------------------

1. make sure the Operator is installed

![](o1.png)

2. Check all "Image Vulnerabilities" from OverView page

![](02.png)

3. Check "Vulnerabilities" of a Pod

![](3.png)


Test container-security-operator from CLI
------------------------------

Get a list of all detected CVEs in pods running on the cluster:

```
$ oc get imagemanifestvuln --all-namespaces
NAMESPACE                   NAME                                                                      AGE
myproject                   sha256.b176867581c15c7bf937757df9207dcd25924789a640af2ed1837a317f3ace25   71m
openshift-cluster-version   sha256.7ad540594e2a667300dd2584fe2ede2c1a0b814ee6a62f60809d87ab564f4425   26h
```

Get a list of all the pods affected by vulnerable images detected by the Operator:

```
$ oc get imagemanifestvuln -o json | jq '.items[].status.affectedPods' | jq -r 'keys[]' | sort -u
myproject/nginx-77f76d44b8-552mh
```

Get a list of all the images that contain vulnerabilities

```
$ oc get imagemanifestvuln -o json | jq -r '.items[].spec.image' | sort -u
quay.io/brandon_tsai/testlab
```


Conclusion
-----------

It is important to check is there any new vulnerabilities in the Production images. we should always patch the images to fix any vulnerabilities found by this operator.   