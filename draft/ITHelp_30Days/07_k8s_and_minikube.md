Kubernetes Introduction
------------------------

While Docker provided an open standard for packaging and distributing containerized applications, there are some problems:
- How would all of these containers be coordinated and scheduled?
- How do you seamlessly upgrade an application without any interruption of service?
- How do you monitor the health of an application, know when something goes wrong and seamlessly restart it?

That is why we need kubernetes to help us to manage these Dockerize applications, especially when your applications are rapid expanding and you have only few operaters to maintain these service.

It gives you the platform to schedule and run containers on clusters of physical or virtual machines (VMs) in production environments and it also can completely manage the life cycle of containerized applications and services using methods that provide predictability, scalability, and high availability. Furthermore, because Kubernetes is all about automation of operational tasks, you can do many of the same things other application platforms or management systems let you do — but for your containers.

Tools to install K8S
--------------------


**Minikube**
Minikube can run on Windows and MacOS, because it relies on virtualization (e.g. Virtualbox) to deploy a kubernetes cluster in a Linux VM. You can also run minikube directly on linux with or without virtualization. It also has some developer-friendly features, like add-ons.

Minikube is currently limited to a single-node Kubernetes cluster, therefore it should be used for app development and test only.


**K3s**
K3s runs on any Linux distribution without any additional external dependencies or tools. It is marketed by Rancher as a lightweight Kubernetes offering suitable for edge environments, IoT devices, CI pipelines, and even ARM devices, like Raspberry Pi's. K3s achieves its lightweight goal by stripping a bunch of features out of the Kubernetes binaries (e.g. legacy, alpha, and cloud-provider-specific features), replacing docker with containerd, and using sqlite3 as the default DB (instead of etcd). As a result, this lightweight Kubernetes only consumes 512 MB of RAM and 200 MB of disk space. K3s has some nice features, like Helm Chart support out-of-the-box.

Unlike minikube, K3s can do multiple node Kubernetes cluster. However, due to technical limitations of SQLite, K3s currently does not support High Availability (HA), as in running multiple master nodes. The K3s team plans to address this in the future.

**Kubeadm**
The official CNCF tool for provisioning Kubernetes clusters in a variety of shapes and forms (e.g. single-node, multi-node, HA, self-hosted)

Although this is the most manual way to create and manage a cluster of all the offerings listed here.

**kops**
kops helps you create, destroy, upgrade and maintain production-grade, highly available, Kubernetes clusters on AWS from the command line. It is very convenient if you want to set up self-own kubernetes cluster on AWS.


Furthermore, there are many Startups provide fully Kubernetes cluster deployment and management solution. These solutions can help you to set up monitoring tools, management UI, and other useful plugins, such as [kubematic](https://www.kubermatic.com/products/kubermatic/),  [Rancher labs](https://rancher.com/) and [Red Hat OpenShift](https://www.openshift.com/)


Install Minikube on Mac
------------------------

You can refer https://kubernetes.io/docs/tasks/tools/install-minikube/ to install Minikube on your Machine.

Following is the script to install Minikube on MacOS

```bash
$ brew install minikube
$ sudo mv minikube /usr/local/bin
$ curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/darwin/amd64/kubectl"
$ chmod +x ./kubectl
$ sudo mv ./kubectl /usr/local/bin/kubectl
$ kubectl version --client
$ brew install hyperkit
$ brew link --overwrite hyperkit
$ minikube start --driver=hyperkit
$ minikube status
```


Deploy your first Pods
---------------------------------------

A Pod (as in a pod of whales or pea pod) is a group of one or more containers, with shared storage/network resources, and a specification for how to run the containers. A Pod's contents are always co-located and co-scheduled, and run in a shared context. A Pod models an application-specific "logical host": it contains one or more application containers which are relatively tightly coupled. In non-cloud contexts, applications executed on the same physical or virtual machine are analogous to cloud applications executed on the same logical host.

Following is the a nginx example.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  labels:
    app: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
```

Apply the yaml file to create the new pod.

```bash
$ alias k=kubectl
$ k apply -f nginx-pod.yaml
pod/nginx-pod created
```

Check the pods is running

```bash
$ k get pods
NAME        READY   STATUS    RESTARTS   AGE
nginx-pod   1/1     Running   0          116s
```

if you want to go into the container for debuging, it is very similar to docker command.

```bash
k exec -it nginx-pod bash
```

and you can check the logs as well

```
k logs nginx-pod
```


Create another namespace
---------

Kubernetes namespaces provide multiple virtual clusters backed by the same physical cluster. It help different projects, teams, or customers to share a Kubernetes cluster.

```
$ k get namespaces
NAME              STATUS   AGE
default           Active   2h
kube-node-lease   Active   2h
kube-public       Active   2h
kube-system       Active   2h
```

Kubernetes starts with four initial namespaces:

- default: The default namespace for objects with no other namespace.
- kube-system: The namespace for objects created by the Kubernetes system.
- kube-public: This namespace is created automatically and is readable by all users (including those not authenticated). This namespace is mostly reserved for cluster usage, in case that some resources should be visible and readable publicly throughout the whole cluster. The public aspect of this namespace is only a convention, not a requirement.
- kube-node-lease: This namespace for the lease objects associated with each node which improves the performance of the node heartbeats as the cluster scales.

let's start creating a new namespace through YAML for our app

```YAML
apiVersion: v1
kind: Namespace
metadata:
  name: myapp
```

```
$ k create -f ./my-namespace.yaml
namespace/myapp created
```

Check the new namespace has created.

```
$ k get namespaces
NAME              STATUS   AGE
default           Active   2h
kube-node-lease   Active   2h
kube-public       Active   2h
kube-system       Active   2h
myapp             Active   26s
```

Apply Pod to New Namespaces
--------------------------

Apply YAML file on new namespace via argument `-n`

```
$ k apply -f nginx-pod.yaml -n myapp
pod/nginx-pod created
```

Get the pods resource in myapp namesapce
```
$ k get pods -n myapp
NAME        READY   STATUS    RESTARTS   AGE
nginx-pod   1/1     Running   0          20s
```


Conclusion
-----------

I hope this article give you a quick review about the basic concept of Kubernetes before introducing the openshift. If you want to learn more about Kubernetes, you can refer the [official tutorials](https://kubernetes.io/docs/tutorials/).
