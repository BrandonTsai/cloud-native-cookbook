
Day 2: Dockerfile 實作範例
=======================


前言
------------

如果您無法從現有的 Image Registry 找到合適的 Container Image，與其將應用程序的靜態代碼掛載到正在運行的容器中，你應攥寫Dockerfile 並使用 Dockerfile 來構建自己的 Container Image。 該文章將首先介紹 Dockerfile 的基本格式，然後針對不同場景提出一些進階用法。

Build Images with Dockerfile
----------------------------

有關Dockerfile格式的更多詳細信息，可以參考 [Dockerfile官方文件]（https://docs.docker.com/engine/reference/builder/），了解更多有關Dockerfile的用法。

以下是基本的 Dockerfile 範例，它將基於 registry.access.redhat.com/ubi7/nodejs-12 構建自定義的 Node.js 應用程序映像。


```
FROM registry.access.redhat.com/ubi7/nodejs-12:1

COPY myapp/ /usr/src/app
ADD default-configs.tar.gz /usr/src/app/

WORKDIR /usr/src/app

RUN npm install && \
    mkdir -p /usr/logs && \
    chown -R "1001:0" /usr/logs && \
    chmod -R u+w /usr/logs

USER 1001

EXPOSE 8000

CMD ["node","hello-http.js"]
```


### FROM

```
FROM <image>[:<tag>]
```

``FROM``指令初始化一個新的構建階段並為後續指令設置基礎映像(Base Image)。 因此，有效的 Dockerfile 必須以 ``FROM`` 指令開頭。 該基礎映像(Base Image)可以是從任何公開的映像倉庫（Image Registry）中取得，例如 Dockerhub。


最常見的基礎映像(Base Image)是 Alpine。 但是，如果您的應用程序之前是部署在 Red Hat 系統上，那麼您也可以使用[Red Hat UBI Image]（https://www.redhat.com/en/blog/introducing-red-hat-universal-base-image）當作你的基礎映像。 您可以在[Red Har Registry]（https://catalog.redhat.com/software/containers/explore）上找到許多基於Red Hat UBI Image 構建的映像，並且可以將其中之一用作你的基礎映像。


### COPY & ADD

```
ADD [--chown=<user>:<group>] <src>... <dest>
COPY [--chown=<user>:<group>] <src>... <dest>
```

``ADD``和``COPY``指令都是從``<src>``複製新文件，目錄或遠程文件URL，並將它們添加到映像的文件系統中的路徑``<dest>``處。

但是 ``COPY`` 僅支持將本地文件基本複製到容器中，而``ADD``具有某些功能，例如僅tar檔解壓縮後放到目的路徑或從遠程URL下載。 例如：

```
ADD http://example.com/font.js /opt/
ADD my_big_lib.tar.gz /var/lib/myapp
```

對於不需要``ADD``的 tar 自動解壓縮功能的項目，您應始終使用 ``COPY``。 不鼓勵使用``ADD``從遠程URL獲取包/文件; 您應該使用``curl``或``wget``。



### WORKDIR

```
WORKDIR /path/to/workdir
```

``WORKDIR``指令為Dockerfile中設置工作目錄，之後的指令都在此工作目錄上運行。你可以想像為 ``cd /path/to/workdir``

### RUN

```
# shell form, the command is run by /bin/sh -c
RUN <command>

# exec form
RUN ["executable", "param1", "param2"]
```

``RUN``指令將在當前 Image 頂部的新層中執行任何命令並提交結果。

你應該確保首先執行最通用的步驟，然後執行最長的步驟，然後再將其緩存，從而使您能夠在快速重建你的 Image。


### ENTRYPOINT & CMD

```
# exec form, this will use "exec" as default ENTRYPOINT
CMD ["executable","param1","param2"]

# shell form, the default ENTRYPOINT is "/bin/sh -c"
CMD command param1 param2
```

``ENTRYPOINT``允許您配置任何可執行腳本當作容器啟動時運行的初始腳本。默認是``exec`` 或 ``/bin/sh -c``。
``CMD``為執行中的容器提供初始指令。 Dockerfile 中只能有一個 ``CMD``指令。 如果您列出多個``CMD``，則只有最後一個``CMD``才會生效。

``CMD``可用於為``ENTRYPOINT``指令提供默認參數，在這種情況下，``CMD``和``ENTRYPOINT``指令均應使用JSON Array。 例如：

```
# Use CMD as default parameters to ENTRYPOINT
ENTRYPOINT ["top", "-b"]
CMD ["-c"]
```


### USER

```
USER <user>[:<group>]
USER <UID>[:<GID>]
```

``USER`` 指令設置運行映像時使用的使用者（或UID）以及可選的群組（或GID）。
Dockerfile 中在此指令之後的任何RUN，CMD 和 ENTRYPOINT 指令都會以該使用者來執行。


### EXPOSE

```
EXPOSE <port>[/<protocol>]
```

``EXPOSE``指令通知 Docker 容器在運行時監聽指定的網絡端口。 您可以指定端口是偵聽TCP還是UDP，如果未指定協議，則默認值為TCP。



### Build Image

要用 Dockfile 構建你的映像檔，您只需執行

```
docker build  -t my-nodejs-app:0.1.0 .
```

如果 Dockerfile 是在不同資料夾或檔名不是預設的Dockerfile，你可以利用 ``-f`` 來標示你的 Dockerfile 位置

```
docker build  -t my-nodejs-app:0.1.0 -f docker/Dockerfile
```


Build with argument
--------------------

如果您不想為不同的環境實現不同的Dockerfile，則可以使用``ARG``指令。ARG指令定義了一個參數，用戶可以在使用docker build命令時加上 ``--build-arg <varname> = <value>``，將其參數值傳遞給Docker。

> 不建議參數來傳遞諸如github密鑰，用戶憑據等機密。構建時參數值對於使用docker history命令的任何用戶都是可見的。 請參閱“ [使用BuildKit生成映像]（https://docs.docker.com/develop/develop-images/build_enhancements/#new-docker-build-secret-information）”部分，了解在構建圖像時可以使用的安全方法 。

以下是一個簡單的範例示範如何透過參數來生成給不同環境執行的映像檔（Image）

```
FROM registry.access.redhat.com/ubi7/nodejs-12:1

ARG APP_ENV=qa
ENV REACT_APP_ENV=${APP_ENV}

COPY myapp/ /usr/src/app
ADD default-configs.tar.gz /usr/src/app/

WORKDIR /usr/src/app

RUN npm install && \
    rm -f .npmrc

RUN chown -R "1001:0" /usr/src/app && \
    chmod -R u+w /usr/src/app && \
    mkdir -p /usr/logs && \
    chown -R "1001:0" /usr/logs && \
    chmod -R u+w /usr/logs

USER 1001

EXPOSE 8000

CMD ["node","hello-http.js"]
```


然後透過以下指令來加上參數，構建Docker Image

```
docker build --build-arg APP_ENV=production -t my-app:0.1.0 .
```



Multi-stage Build
------------------

關於構建Image的最具挑戰性的事情之一是減小Image檔案大小。 Dockerfile 中的每條指令都會在 Image 上添加一層，您需要記住在移至下一層之前清除不需要的任何文件。 為了編寫一個真正有效的 Dockerfile，傳統上，您需要使用 Shell 技巧和其他邏輯來使各層盡可能小，並確保每一層都具有上一層所需的文件，而沒有其他任何東西。

當我們要運行容器時，往往可能只需要整個構建過程最後生成的檔案。而在構建過程過可能還需要額外安裝一些套件過可能還需要額外安裝一些最後不需要的套件。 通過 Multi-stage Build，您可以在 Dockerfile 中使用多個``FROM``語句將構建過程區分成不同的階段。每個FROM指令可以使用不同的基礎映像(Base Image)，而您可以將檔案從一個階段複製到另一個階段，從而在最終的Image中僅留下需要的所有內容。

>此功能需要 Docker 17.05 或更高版本。 如果您在RHEL7系統上工作，強烈建議您使用 Podman 或安裝 社群版本 (docker-ce) 而不是預設的Docker。

以下是Multi-stage Build的簡單範例：


```
FROM registry.access.redhat.com/ubi7/nodejs-12:1 as builder

WORKDIR /usr/src/app

COPY . /usr/src/app
RUN yarn global add node-gyp && yarn

ARG APP_ENV=qa
ENV REACT_APP_ENV=${APP_ENV}

RUN npm run compile && \
    rm -rf node_modules && \
    yarn


FROM registry.access.redhat.com/ubi7/nodejs-12:1

ARG APP_ENV=qa
ENV REACT_APP_ENV=${APP_ENV}

WORKDIR /usr/src/app
COPY --from=builder /usr/src/app .

EXPOSE 9000
CMD npm run production
```

然後用一般的docker build指令即可建構此映像：

```
docker build --build-arg APP_ENV=production -t my-app:0.1.0 .
```

結論
-----

本文涵蓋了 Dockerfile ㄧ些基礎指令介紹和一點進階用法。 Docker的世界博大精深，建議參考官方文件繼續精進．願讀者們都能在這片茫茫大海中找到一點小小的樂趣。



