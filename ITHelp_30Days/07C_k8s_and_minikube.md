為何我們要 Kubernetes?
---------------------

-----

儘管 Docker 為容器化的應用程序提供了標準的部屬方式，但仍然存在一些問題：
- 容器之間如何溝通協調？ 尤其當他們可能部署在不同機器時
- 如何無縫升級應用程序而不中斷服務？
- 如何監視應用程序的運行狀況，知道出了什麼問題並無縫地重新啟動它？


這就是為什麼我們需要 Kubernetes 來幫助我們管理這些容器化的應用程序，尤其是當您的應用程序正在快速擴展並且您只有很少的人員來維護這些服務時。 Kubernetes 提供了幫助您在實體機器或ＶＭ上調度和運行容器，並且提供可預測，高彈性和高可用性的方法來完全管理容器化應用程序的生命週期。 由於 Kubernetes 主要設計的目的就是要自動化和標準化操作營運服務時的任務，因此您可以做原本其他平台或系統可以做的許多相同的事情 - 前提是您要先容器化您的應用程序。


安裝 Kubernetes 的工具
--------------------

-----

**（1）Minikube**

Minikube 可以在 Windows 和 MacOS 上運行，因為它依賴於虛擬機（例如Virtualbox）在 Linux VM 中部署kubernetes，或您也可以在Linux上直接運行minikube而不用透過虛擬機。 但是 Minikube 當前僅限支援單一運算節點（ Computer Node），因此僅適用於應用程序開發和測試。


**（2）K3s**
K3s 可在任何 Linux 上運行，而無需任何其他外部工具。 Rancher將該產品作為輕量級的 Kubernetes 進行銷售，該產品適用於，IoT設備，CICD Pipeline，甚至是 Raspberry Pi的 ARM 設備。 K3s 使用 sqlite3 作為默認數據庫（而不是etcd）來實現其輕量級目標。 因此，這種輕量級的 Kubernetes 僅消耗 512 MB 的 RAM 和 200 MB 的硬碟空間。 K3s具有一些不錯的功能，例如安裝好即可使用 Helm Chart 來部署您的應用程序。

與minikube不同，K3s 可以支援多個運算節點。 但是，由於 SQLite 的技術限制，K3當前不支持高可用性（HA）。

**（3）Kubeadm**
官方的CNCF工具，有比較高的自由度來設定Kubernetes平台（例如，單運算節點，多運算節點，高可用性）。但安裝過程中有很多步驟需要手動設定。

**（4）kops**
kops 可幫助您利用命令行在AWS上創建，銷毀，升級和維護高可用性的Kubernetes服務。 如果想要在 AWS 上設置自己管理的kubernetes，他是一個很方便的工具。



此外，也有許多公司提供完整的 Kubernetes 部署和管理解決方案。 這些解決方案可以幫助您設置監控工具，提供管理介面和其他有用的插件，例如 [kubematic](https://www.kubermatic.com/products/kubermatic/)，[Rancher labs](https://rancher.com/) 和 [Red Hat OpenShift](https://www.openshift.com/)。



在 MAC 安裝 Minikube
------------------------

-----

以下是簡單的步驟在 MAC 安裝 Minikube 來做測試。

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

> 你如果要安裝在其他作業系統，可以參考 https://kubernetes.io/docs/tasks/tools/install-minikube/


在 Minikube 部署您的第一個 Pod
---------------------------------------

-----

Pod 是 K8S 中基本的單位，是由一個或多個容器以及有關如何運行容器的規範組成，同一個 Pod 裡面的 Container 能夠透過 localhost 互相的連線並且共享及使用同一個IP。所有的容器也都會同生共死，只要有一個容器運行失敗，那麼其他容器也會跟著失敗。

Kubernetes 可以透過 YAML 檔案來管理和部署所有的服務及資源。以下是部署一個包含 nginx 容器的 Pod 的範例

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

然後我們就可以用 kubectl 指令及此 YAML 檔案來創造一個新的 Pod。

```bash
$ alias k=kubectl
$ k apply -f test.yaml
pod/nginx-pod created
```

確定 Pod 有在運行。

```bash
$ k get pods
NAME        READY   STATUS    RESTARTS   AGE
nginx-pod   1/1     Running   0          116s
```

與docker非常相似，我們可以透過指令進入容器內進行除錯。

```bash
kubectl exec -it nginx-pod bash
```

也可以檢查 Logs

```
kubectl logs nginx-pod。
```


建立 Namespace
-------------

-----

在同一個 Kubernetes 內，我們可以藉由 Namespace 來提供了多個 virtual clusters．它可以幫助不同的項目，團隊或客戶共享同一個Kubernetes資源。

```
$ k get namespaces
NAME              STATUS   AGE
default           Active   2h
kube-node-lease   Active   2h
kube-public       Active   2h
kube-system       Active   2h
```

Kubernetes從四個初始 Namespaces 開始：
- default：默認的Namespaces。
- kube-system：Kubernetes 系統創建時用的 Namespaces。
- kube-public：此 Namespaces 是自動創建的，並且所有用戶（包括未經身份驗證的用戶）均可讀取。
- kube-node-lease：此 Namespaces 用於改善Node Heartbeat的性能。

讓我們透過YAML為我們的應用程序創建一個新的 Namespace


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

檢查新的 Namespace 是否創建了。

```
$ k get namespaces
NAME              STATUS   AGE
default           Active   2h
kube-node-lease   Active   2h
kube-public       Active   2h
kube-system       Active   2h
myapp             Active   26s
```

在新的 Namespaces 建立一個 Pod
--------------------------

-----


通過參數 `-n` 在新 Namespaces 上部署 Pod。

```
$ k apply -f nginx-pod.yaml -n myapp
pod/nginx-pod created
```

列出 "myapp" Namesapce 內所有的 Pod。
```
$ k get pods -n myapp
NAME        READY   STATUS    RESTARTS   AGE
nginx-pod   1/1     Running   0          20s
```


結論
-----------

-----

我希望本文能夠在介紹 OpenShift 之前讓您快速了解 Kubernetes 的基本概念。 如果您想學習更多有關 Kubernetes 的基本概念，可以參考[官方文件](https://kubernetes.io/docs/tutorials/)。