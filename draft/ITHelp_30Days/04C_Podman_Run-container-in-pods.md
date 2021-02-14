Day 3: Podman 淺談 - Pod 功能介紹
==============================

何謂Pod？
------

-----


Pod 概念由[Kubernetes]（https://kubernetes.io/docs/concepts/workloads/pods/pod/）引入。 Pod 是由多個運行在一起的容器組成。 Podman 可以使用類似的概念來管理本地服務器上的容器。 一個 Pod 中的所有容器共享相同的網絡，因此它們可以輕鬆地通過 localhost 相互通信，而無需設定任何額外的服務端口。 您可以參考 [Podman：在本地端管理Pod和容器]（https://developers.redhat.com/blog/2019/01/15/podman-managing-containers-pods/）了解有關技術的更多詳細信息。在這篇文章中，我們將重點介紹如何在本地服務器上運行和管理Pod。

手動設定及運行 Pod
----------------

-----


首先要做的是創建一個新的 Pod。

```
# sudo podman pod create -n my-app -p 8081:80
```

然後把新增一個容器到此 Pod。要注意到，在這個步驟你不能綁定一個服務端口到特定的容器. 你必須之後將服務端口綁定在 Pod。
```
sudo podman run -dt --pod my-app -v /opt/http:/usr/share/nginx/html:ro --security-opt="seccomp=unconfined" --name hello-nginx nginx
```


你可以透過指令 `podman pod ps` 列出所有的 Pod 及包含的容器數量。

```
# sudo podman pod ps
POD ID         NAME     STATUS    CREATED          # OF CONTAINERS   INFRA ID
75d943416fc8   my-app   Created   26 minutes ago   1                 30138c8d0d1c
```


如果你停止一個 Pod ，那麼它包含的所有容器也都會停止運行。


```
$ sudo podman pod stop my-app
a2edfd1095760b1e2946271184743cce6f621665878b618ddc83d73b295070ba
$ sudo podman ps -a
CONTAINER ID  IMAGE                           COMMAND               CREATED        STATUS                    PORTS                 NAMES
cacdc75990b0  docker.io/library/nginx:latest  nginx -g daemon o...  2 minutes ago  Exited (0) 7 seconds ago  0.0.0.0:8082->80/tcp  hello-nginx
4dce350e01cf  k8s.gcr.io/pause:3.1                                  3 minutes ago  Exited (0) 7 seconds ago  0.0.0.0:8082->80/tcp  a2edfd109576-infra
```

同理，啟動一個 Pod  也會啟動裡面的所有容器。

```
$ sudo podman pod start my-app
a2edfd1095760b1e2946271184743cce6f621665878b618ddc83d73b295070ba
$ sudo podman ps -a
CONTAINER ID  IMAGE                           COMMAND               CREATED        STATUS            PORTS                 NAMES
cacdc75990b0  docker.io/library/nginx:latest  nginx -g daemon o...  4 minutes ago  Up 5 seconds ago  0.0.0.0:8082->80/tcp  hello-nginx
4dce350e01cf  k8s.gcr.io/pause:3.1
```


利用 YAML 檔案來建立及設定 Pod。
-----------------------------

-----


Podman 可以透過跟 Kubernetes 相容的 YAML 檔案來定義一個。
而且支援掛載 [hostPath](https://kubernetes.io/docs/concepts/storage/volumes/#hostpath) 到容器內。
如以下範例：

```yaml
# my-app.yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-app
spec:
  containers:
  - name: ng1
    image: nginx
    ports:
      - containerPort: 8001
        hostPort: 8001
        protocol: TCP
    volumeMounts:
      - name: html1-volume
        mountPath: /opt/html
      - name: config1-volume
        mountPath: /etc/nginx/conf.d
  - name: ng2
    image: nginx
    ports:
      - containerPort: 8002
        hostPort: 8002
        protocol: TCP
    volumeMounts:
      - name: html2-volume
        mountPath: /opt/html
      - name: config2-volume
        mountPath: /etc/nginx/conf.d
  volumes:
    - name: html1-volume
      hostPath:
        path: /opt/myapp/html1
        type: Directory
    - name: config1-volume
      hostPath:
        path: /opt/myapp/config1
        type: Directory
    - name: html2-volume
      hostPath:
        path: /opt/myapp/html2
        type: Directory
    - name: config2-volume
      hostPath:
        path: /opt/myapp/config2
        type: Directory
```

有了YAML檔之後，我們就可以透過以下指令來運行Pod。

```
$ sudo podman play kube ./my-app.yaml
```


確認透過YAML檔設定的 Pod 及容器有成功運行
```
$ podman ps -a
CONTAINER ID  IMAGE                           COMMAND               CREATED        STATUS                        PORTS                             NAMES
2268e5ab9b61  nginx                           nginx -g daemon o...  2 minutes ago  Up 2 minutes ago              0.0.0.0:8001-8002->8001-8002/tcp  ng2
19dba831eeae  nginx                           nginx -g daemon o...  2 minutes ago  Up 2 minutes ago              0.0.0.0:8001-8002->8001-8002/tcp  ng1
42c150972ddb  k8s.gcr.io/pause:3.1                                  2 minutes ago  Up 2 minutes ago              0.0.0.0:8001-8002->8001-8002/tcp  4ae6b24effb5-infra
```

> 注意，Pod 裡面的所有容器都共享相同的 IP 127.0.0.1, 所以他們必須設定不同的服務端口．否則容器會無法成功啟動．



結論
-----

-----


對於 Kubernetes 用戶，因為 Podman 可以使用 Kubernetes 相容的 Pod 部署方式，他們應該可以非常熟練的實作 YAML 檔來在本地端管理多個容器。

但是，與 docker-compose 相比，Podman 的 YAML 檔 卻不能用於同時構建多個映像檔。 即使有第三方工具[podman-compose]（https://github.com/muayyad-alsadi/podman-compose）可能帶來類似的功能。 但是我會建議直接用其他CI/CD工具或Shell Script來自動建立映像檔，然後使用Podman的Pod功能進行容器管理。

> 本篇英文版發佈在 https://darumatic.com/blog/podman_run_container_in_pods