
Instead of implementing shell scripts to perform some complex behavior with kubectl, you can think of installing kubectl plugins to extended kubectl command to achieve custom functionality. The easiest way to manage plugins is [Krew](https://github.com/kubernetes-sigs/krew). Krew is a tool that aims to ease plugin discovery, installation, upgrade, and removal on multiple operating systems. This article will show you how easy it is to grab and experiment with existing plugins.



Install Krew
-------------

Please refer: https://krew.sigs.k8s.io/docs/user-guide/setup/install/ to install Krew

1. Make sure that git is installed.

2. Run this command in your terminal to download and install krew:

```bash
(
  set -x; cd "$(mktemp -d)" &&
  curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/krew.{tar.gz,yaml}" &&
  tar zxvf krew.tar.gz &&
  KREW=./krew-"$(uname | tr '[:upper:]' '[:lower:]')_amd64" &&
  "$KREW" install --manifest=krew.yaml --archive=krew.tar.gz &&
  "$KREW" update
)
```

3. Add $HOME/.krew/bin directory to your PATH environment variable. To do this, update your .bashrc or .zshrc file and append the following line:

```bash
export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"
```

4. restart your shell and verify running kubectl krew works.

```bash
$ kubectl krew
krew is the kubectl plugin manager.
You can invoke krew through kubectl: "kubectl krew [command]..."

Usage:
  krew [command]

Available Commands:
  help        Help about any command
  info        Show information about a kubectl plugin
  install     Install kubectl plugins
  list        List installed kubectl plugins
  search      Discover kubectl plugins
  uninstall   Uninstall plugins
  update      Update the local copy of the plugin index
  upgrade     Upgrade installed plugins to newer versions
  version     Show krew version and diagnostics

Flags:
  -h, --help      help for krew
  -v, --v Level   number for the log level verbosity

Use "krew [command] --help" for more information about a command.
```


Search plugins and get the detail of particular plugin
-------------------------------------------------------

List all plugins

```bash
$ kubectl krew search
NAME                            DESCRIPTION                                         INSTALLED
access-matrix                   Show an RBAC access matrix for server resources     no
advise-psp                      Suggests PodSecurityPolicies for cluster.           no
apparmor-manager                Manage AppArmor profiles for cluster.               no
auth-proxy                      Authentication proxy to a pod or service            no
bulk-action                     Do bulk actions on Kubernetes resources.            no
ca-cert                         Print the PEM CA certificate of the current clu...  no
capture                         Triggers a Sysdig capture to troubleshoot the r...  no
change-ns                       View or change the current namespace via kubectl.   no
...

```

Get plugin detail by **kubectl krew info \<Plugin Name\>**

```bash
$ kubectl krew info grep
NAME: grep
URI: https://github.com/guessi/kubectl-grep/releases/download/v1.2.2/kubectl-grep-Darwin-x86_64.tar.gz
SHA256: cef6f2642ba8f284e2a675314cfb352f53ce2b9ea0202348e4a9015a5f8f66be
VERSION: v1.2.2
HOMEPAGE: https://github.com/guessi/kubectl-grep
DESCRIPTION:
Filter Kubernetes resources by matching their names

Examples:

List all pods in all namespaces
$ kubectl grep pods --all-namespaces

List all pods in namespace "star-lab" which contain the keyword "flash"
$ kubectl grep pods -n star-lab flash

No more pipe, built-in grep :-)

CAVEATS:
\
 | This plugin requires an existing KUBECONFIG file, with a `current-context` field set.
 | 
 | Usage:
 | 
 |   $ kubectl grep # output help messages
 | 
 | More Info:
 | - https://github.com/guessi/kubectl-grep
/

```

Useful Plugins
---------------

### **ns**

Also known as "kubens", a utility to set your current namespace and switch
between them.

```bash
$ kubectl krew install ns
$ kubectl ns kube-system
Context "minikube" modified.
Active namespace is "kube-system".
$ kubectl get pods
NAME                               READY   STATUS    RESTARTS   AGE
coredns-66bff467f8-7xkvh           1/1     Running   2          41d
coredns-66bff467f8-wkktf           1/1     Running   2          41d
etcd-minikube                      1/1     Running   2          41d
kube-apiserver-minikube            1/1     Running   2          41d
kube-controller-manager-minikube   1/1     Running   2          41d
kube-proxy-vwm6d                   1/1     Running   2          41d
kube-scheduler-minikube            1/1     Running   2          41d
storage-provisioner                1/1     Running   3          41d
```

### **ctx**

Also known as "kubectx", a utility to switch between context entries in
your kubeconfig file efficiently.

This plugin is very useful when you need to manage multiple Kubernetes clusters

```bash
$ kubectl krew install ctx
$ kubectl ctx
minikube
$ kubectl ctx minikube
Switched to context "minikube"
```

### **get-all**

This plugin is similar to "kubectl get all" but really everything

For example, no resources can found in kube-public namespace by original "kubectl get all" command

```bash
$ kubectl get all -n kube-public
No resources found in kube-public namespace.
```

With this plugin, we can get a list of all resources including configmaps, secrets and serviceaccount, etc.

```bash
$ kubectl krew install get-all
$ kubectl get-all -n kube-public
NAME                                                                        NAMESPACE    AGE
configmap/cluster-info                                                      kube-public  41d  
secret/default-token-lg6sh                                                  kube-public  41d  
serviceaccount/default                                                      kube-public  41d  
rolebinding.rbac.authorization.k8s.io/kubeadm:bootstrap-signer-clusterinfo  kube-public  41d  
rolebinding.rbac.authorization.k8s.io/system:controller:bootstrap-signer    kube-public  41d  
role.rbac.authorization.k8s.io/kubeadm:bootstrap-signer-clusterinfo         kube-public  41d  
role.rbac.authorization.k8s.io/system:controller:bootstrap-signer           kube-public  41d  

```

### **images**

This plugin shows container images used in the Kubernetes cluster in a
table view. You can show all images or show images used in a specified
namespace.

```bash
$ kubectl krew install images
$ kubectl images -A
[Summary]: 2 namespaces, 9 pods, 9 containers and 8 different images
+-----------------------------------+-------------------------+------------------------------------------------+
|              PodName              |      ContainerName      |                 ContainerImage                 |
+-----------------------------------+-------------------------+------------------------------------------------+
| nginx-deployment-6b768b47c4-99g55 | nginx                   | nginx:1.14.2                                   |
+-----------------------------------+-------------------------+------------------------------------------------+
| coredns-66bff467f8-7xkvh          | coredns                 | k8s.gcr.io/coredns:1.6.7                       |
+-----------------------------------+                         +                                                +
| coredns-66bff467f8-wkktf          |                         |                                                |
+-----------------------------------+-------------------------+------------------------------------------------+
| etcd-minikube                     | etcd                    | k8s.gcr.io/etcd:3.4.3-0                        |
+-----------------------------------+-------------------------+------------------------------------------------+
| kube-apiserver-minikube           | kube-apiserver          | k8s.gcr.io/kube-apiserver:v1.18.2              |
+-----------------------------------+-------------------------+------------------------------------------------+
| kube-controller-manager-minikube  | kube-controller-manager | k8s.gcr.io/kube-controller-manager:v1.18.2     |
+-----------------------------------+-------------------------+------------------------------------------------+
| kube-proxy-vwm6d                  | kube-proxy              | k8s.gcr.io/kube-proxy:v1.18.2                  |
+-----------------------------------+-------------------------+------------------------------------------------+
| kube-scheduler-minikube           | kube-scheduler          | k8s.gcr.io/kube-scheduler:v1.18.2              |
+-----------------------------------+-------------------------+------------------------------------------------+
| storage-provisioner               | storage-provisioner     | gcr.io/k8s-minikube/storage-provisioner:v1.8.1 |
+-----------------------------------+-------------------------+------------------------------------------------+
```



### **whoami**

This plugin can be used to show the subject that's currently authenticated as.

```bash
$ kubectl krew install whoami
$ kubectl whoami
kubecfg:certauth:admin
```

### **rbac-lookup**

This plugin allows the cluster manager to easily find Kubernetes roles and cluster roles bound to any user, service account, or group name.

```bash
$ kubectl krew install rbac-lookup
$ kubectl rbac-lookup --kind user
SUBJECT                           SCOPE          ROLE
system:anonymous                  kube-public    Role/kubeadm:bootstrap-signer-clusterinfo
system:kube-controller-manager    kube-system    Role/extension-apiserver-authentication-reader
system:kube-controller-manager    kube-system    Role/system::leader-locking-kube-controller-manager
system:kube-controller-manager    cluster-wide   ClusterRole/system:kube-controller-manager
system:kube-proxy                 cluster-wide   ClusterRole/system:node-proxier
system:kube-scheduler             kube-system    Role/extension-apiserver-authentication-reader
system:kube-scheduler             kube-system    Role/system::leader-locking-kube-scheduler
system:kube-scheduler             cluster-wide   ClusterRole/system:kube-scheduler
system:kube-scheduler             cluster-wide   ClusterRole/system:volume-scheduler
```

### **grep**

This plugin can filter Kubernetes resources by matching their names instead of search resource by pipe.

Before installing this plugin, we usually filter pods by the following commands,

```bash
$ kubectl get pods            | grep "keyword"
$ kubectl get pods -o wide    | grep "keyword"
```

With this plugin installed, you can filter pod with kubectl grep easily

```bash
$ kubectl grep pods "keyword"
$ kubectl grep pods "keyword" -o wide
```

Conclusion
----------

Kubectl command plugins can make life easier when working with kubernetes. If you can not find a plugin that is suitable for you, there is a document in the Kubernetes repo that describes how to [write your own custom plugin](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/).