
Day 1: Docker 淺談
=================

容器 （Container）vs 虛擬機（Virtual Machine）
-------------------------------

![VM-vs-Docker](https://geekflare.com/wp-content/uploads/2019/09/traditional-vs-new-gen.png "VM-vs-Docker")


在雲架構方面，虛擬機(Virtual Machine)已成為其眾多優勢的首選標準。 但是，如果您能有更輕巧，經濟和可擴展的虛擬機替代方案呢？ 這就是 Docker Container。

在這篇文章中，我將解釋 Docker 的基本概念以及 Docker CLI 的基本用法。


Docker 基本概念
-----------

![Docker Concept](images/01_docker/Docker_concept.png "Docker Concept")

容器的系統需要一個底層操作系統，該操作系統使用虛擬內存支持進行隔離，從而為所有容器化應用程序提供基本服務。 Docker包含以下組件。

**Docker Engine**

它是整個Docker系統的核心部分。 Docker Engine 是遵循客戶端-服務器架構的應用程序。它安裝在主機上。

**映像（Image）**

Image 是由多層組成，用於執行 Docker Container 的代碼。Image 本質上是根據指令的應用程序來構建的，該指令依賴於主機OS內核來完成應用程序的完整版本和可執行版本。

**Dockerfile**

Dockerfile 是一個文本文檔，其中包含用戶可以在命令行上調用以組裝Docker Image 的所有命令。用戶可以使用 Dockerfile 構建自定義Docker Image。


**映像倉庫（Image Registry）**

Image Registry 是一個系統用來存儲 Images，這些 Images 具有不同的標記版本(Tags)。最受歡迎的 Image Registry 是 [Docker Hub]（https://hub.docker.com/)，你也可以在 [Red Har Registry](https://catalog.redhat.com/software/containers/explore) 找到 Red Hat 支援的 Image。如果你想要在自己的資料中心建立私有的Image Registry，你可以考慮安裝 [Harbor](https://github.com/goharbor/harbor) 或 [Red Hat Quay](https://www.openshift.com/products/quay)，這兩個方案提供的功能都非常完整。


**容器（Container）**

運行 Docker Image 後，它將創建一個 Docker Container。所有應用程序及其環境都在此 Container 中運行。您可以使用 Docker API 或 CLI 來啟動，停止，刪除 Container。


如何在 Red Hat 系統上安裝 Docker
------------------------------

非常簡單，只要執行以下指令

```
sudo yum install docker
```

請注意，Red Hat 官方不再支持最新版本的 Docker，他們建議用 Podman 取代 Docker，因為它被認為比 Docker 更安全。 但是，您仍然可以在 RHEL7 上安裝 Docker 的社群版本。 請參考 [這裡]（https://computingforgeeks.com/install-docker-ce-on-rhel-7-linux/）在您的 RHEL7 系統上安裝 docker-ce。


運行 Docker Container
-----------------------

在本節中，我將介紹 Docker CLI 的一些基本用法。
在運行 Container 之前，您需要從 Image Registry 中提取（Pull） Image，例如：

```
sudo docker pull registry.access.redhat.com/rhscl/postgresql-10-rhel7:1
```

並確定 Image 已經存在本地端

```
sudo docker images
```

以下是一個簡單的 Shell script 範例，用來在本地端運行 Postgresql Container。

```
#!/bin/bash

data_folder="/opt/postgresql"
pgsql_user="myadmin"
pgsql_password="myPassword"
pgsql_database="myapp"

sudo mkdir -p ${data_folder}
sudo chmod 777 ${data_folder}

sudo docker run -d --name postgresql \
-e POSTGRESQL_USER="${pgsql_user}" \
-e POSTGRESQL_PASSWORD="${pgsql_password}" \
-e POSTGRESQL_DATABASE="${pgsql_database}" \
-v "${data_folder}:/var/lib/pgsql/data:Z" \
-p "5432:5432" \
registry.access.redhat.com/rhscl/postgresql-10-rhel7:1
```


### 增加或複寫環境變數 （Environment Variables）

您可以使用 -e 或 --env-file 來在容器中設置簡單的環境變數，或覆蓋在您的 Dockerfile 中定義的變量。


```
docker run \
-e TZ="${cluster_timezone}" \
-e POSTGRESQL_USER="${pgsql_user}" \
-e POSTGRESQL_PASSWORD="${pgsql_password}" \
-e POSTGRESQL_DATABASE="${pgsql_database}" \
...
```

### Volumes


Volume 可用於持久存儲容器生成或需要使用的數據。
使用 -v 或 -volume 將主機文件夾掛載到容器中。

```
docker run \
-v /host_folder:/opt/config:ro \
...
```

請注意，您必須確保容器內的用戶有權訪問已安裝的主機文件夾。有兩種方法可以實現此目的：

1. 在本地端執行 ``chown -R "uid:gid" /host_folder`` (建議)
2. 在本地端執行 ``chmod -R 777 /host_folder``


### Expose ＆ Publish


Dockerfile 中的 EXPOSE 指令通知 Docker 該容器在運行時偵聽指定的網絡端口。例如 Postgresql 通常需要監聽 5432 端口。

```
EXPOSE 5432/tcp
```

If you EXPOSE a port, the service in the container is not accessible from outside Docker, but from inside other Docker containers.
Use the -p flag to actually publish the port, the service in the container is accessible from anywhere, even outside Docker.

如果當你 EXPOSE 一個服務端口，你只可以從其他 Docker 容器內部訪問該服務端口，無法從 Docker 外部訪問此容器中的服務端口。
因此我們在運行時，需要使用 -p 來實際發布端口，讓容器中的服務端口可以從 Docker 外部訪問。


```BASH
docker run -p "5432:5432"
docker run -p "127.0.0.1:8080:80
```

### ENTRYPOINT & CMD


ENTRYPOINT 用於標識啟動容器時應運行哪個可執行檔
CMD 則是容器啟動後通過 ENTRYPOINT 運行的第一個指令。

如果未指定 ENTRYPOINT，則默認為 "/bin/sh -c" 或 "exec"

根據 Ｄockerfile 中，CMD 指令寫的方式不同，Docker 會用不一樣的 ENTRYPOINT 來執行。

```
# ENTRYPOINT 是 exec （推薦）
CMD ["executable","param1","param2"]

# ENTRYPOINT 是 shell， 以下指令會透過 /bin/sh -c 執行
CMD command param1 param2
```

如果您希望容器每次都運行相同的可執行文件，
那麼您應該考慮結合使用 ENTRYPOINT 和 CMD

```
ENTRYPOINT /my_container_entrypoint.sh
CMD ["param1", "param2"]
```

這也解釋了為什麼設置環境變量可以讓 Postgresql 容器自動創建數據庫， 你可以透過以下指令去查看一個 Image 設定的 ENTRYPOINT 跟 CMD 為何，這可以幫助你了解容器啟動後的行為。

```
sudo docker inspect registry.access.redhat.com/ubi7/nodejs-12:1
```


Docker Exec
-----------

最後來講下 ``docker exec`` 這個指令，該指令可以讓用戶在正在運行的容器中運行新命令。 我們可以使用此命令進入正在運行的容器，並檢查容器中的文件和服務狀態，這在 Debug 時很有用。 例如，

```
docker exec -it postgresql bash
```


結論
-----

總體而言，Docker 允許將應用程序與所有依賴檔案一起打包，從而大大簡化了部署過程，並擁有完整的可複制環境。
希望這個文章讓您對 Docker 容器感到興奮，從而您不再需要透過繁雜的腳本來將服務部署到每台機器上．














