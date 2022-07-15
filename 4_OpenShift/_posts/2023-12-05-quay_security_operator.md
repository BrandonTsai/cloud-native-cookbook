---
title: "#5 Scan Openshift vulnerabilities using Quay Operator "
author: Brandon Tsai
---

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