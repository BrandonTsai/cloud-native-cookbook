
工欲善其事，必先裝 Krew
=====================

在安裝完 Kubernetes 服務後，你可能會發現你有時需要寫 shell 腳本來執行 kubectl 的某些複雜行為或取得特定資訊，這時您可以考慮將插件（plugins）安裝到 kubectl命令中，以實現自定義的功能，簡化管理步驟。而管理插件最簡單的方法是[Krew](https://github.com/kubernetes-sigs/krew)。 本文將透過 Krew 來安裝並介紹一些有用的插件。


安裝 Krew
-------------

1. 確定 git 已經安裝。

2. 執行以下指令來下載並安裝 krew。

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

3. 把 $HOME/.krew/bin 路徑加到 PATH 環境變數. 為此，請更新 `.bashrc` 或 `.zshrc` 文件，並添加：

```bash
export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"
```

4. 重新連線到終端機並執行 `kubectl krew` 確定有安裝成功。

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

搜索並獲取特定插件的詳細信息
--------------------------

搜索所有插件

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

獲取特定插件的詳細信息 ``kubectl krew info <Plugin Name>``

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

插件介紹
--------

### **ns**

也稱為"kubens"，用於設置及切換當前 Namespace。

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

也稱為 "kubectx", 當你有很多 Kubernetes clusters 時，你可以利用此插件來切換到不同的cluster。

```bash
$ kubectl krew install ctx
$ kubectl ctx
minikube
$ kubectl ctx minikube
Switched to context "minikube"
```

### **get-all**

此插件類似 ``kubectl get all`` ，但是真的可以列出所有的資源。
例如，原始的 ``kubectl get all`` 命令無法在 ``kube-public`` Namespaces 中找到任何資源。

```bash
$ kubectl get all -n kube-public
No resources found in kube-public namespace.
```

使用此插件，我們可以獲得所有資源的列表，包括 configmaps, secrets 和 serviceaccount 等.

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

該插件顯示了 Kubernetes 中使用的容器映像檔。


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

該插件可用於顯示當前登入的帳號。

```bash
$ kubectl krew install whoami
$ kubectl whoami
kubecfg:certauth:admin
```

### **rbac-lookup**

Role 和 ClusterRole 是 Kubernetes 用來管理權限的資源。
該插件可以輕鬆找到與任何 User，Group 或 service account 綁定的 Roles 和 Cluster Roles。

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

這個插件可以通過匹配資源的名稱來過濾 Kubernetes 資源列表，而不是通過複雜的指令搜索資源。

在安裝此插件之前，我們通常通過以下命令過濾 Pod 列表，

```bash
$ kubectl get pods            | grep "keyword"
$ kubectl get pods -o wide    | grep "keyword"
```

安裝此插件後，您可以輕鬆地使用 ``kubectl grep`` 來達到相同目的

```bash
$ kubectl grep pods "keyword"
$ kubectl grep pods "keyword" -o wide
```

Conclusion
----------

使用Kubernetes時，Kubectl 插件可以使生活更輕鬆。 如果找不到適合您的插件，可以參考官方手冊：[如何編寫自己的插件](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/).
