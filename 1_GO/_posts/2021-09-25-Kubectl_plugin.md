
In this blog, I will show to how to work with kubernetes resource via Go client
we are going to implement a kubectl plug which will scan the vuls of the images of the pods running on the namespace


```
k trivy = Get image scan result in current namespaces
k trivy -n myapp = Get image scan result in myapp namespace
```



Go Client
---------



Kubectl Plugin
--------------

```bash
$ go build  -o  kubectl-trivy  .
$ cp ./kubectl-trivy /usr/local/bin/
$ kubectl trivy -n  kube-system
There are 9 pods in the cluster
Found 9 podin namespace kube-system
Remote Trivy Server:  127.0.0.1:8080
+--------------------------------------------+----------------------------------------------------+------+--------+-----+--------+
| IMAGE                                      | PODS                                               | HIGH | MEDIUM | LOW | UNKNOW |
+--------------------------------------------+----------------------------------------------------+------+--------+-----+--------+
| k8s.gcr.io/kube-proxy:v1.21.5              | kube-proxy-wcg2b,                                  |   27 |     15 |  73 |      0 |
| k8s.gcr.io/coredns/coredns:v1.8.0          | coredns-558bd4d5db-rkdmc,coredns-558bd4d5db-74bg4, |    3 |      4 |   0 |      1 |
| docker/desktop-vpnkit-controller:v2.0      | vpnkit-controller,                                 |    3 |      1 |   0 |      3 |
| docker/desktop-storage-provisioner:v2.0    | storage-provisioner,                               |    3 |      1 |   0 |      1 |
| k8s.gcr.io/etcd:3.4.13-0                   | etcd-docker-desktop,                               |    0 |      0 |   0 |      5 |
| k8s.gcr.io/kube-scheduler:v1.21.5          | kube-scheduler-docker-desktop,                     |    0 |      0 |   0 |      0 |
| k8s.gcr.io/kube-apiserver:v1.21.5          | kube-apiserver-docker-desktop,                     |    0 |      0 |   0 |      0 |
| k8s.gcr.io/kube-controller-manager:v1.21.5 | kube-controller-manager-docker-desktop,            |    0 |      0 |   0 |      0 |
+--------------------------------------------+----------------------------------------------------+------+--------+-----+--------+

```