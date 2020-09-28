

OpenShift 系列 2： ”懶“之呼吸 第一式 - Source-To-Image (S2I)
=================================================

Source-To-Image (S2I) 淺談
---------------------------

-----

如果您的應用程序不是非常複雜，比起自己撰寫 Dockerfile 並努力解決權限問題，另一種更方便的方法是使用 [Source-To-Image (S2I)](https://github.com/openshift/source-to-image). S2I 使用建構者映像檔（Builder Image）和你的程式代碼來生成新的Docker映像檔。 S2I 的建構者映像檔包括許多常用程式語言和服務，例如Python或Ruby，您也可以使用自己定義的腳本來擴展 S2I 社群資源。


使用 S2I 來建構映像檔有以下好處:

- 速度 - 借助 S2I，建構過程可以執行大量複雜的操作，而無需自己在 Dockerfile 創建新的步驟來解決問題，從而實現了快速的開發過程。
- 用戶效率 - S2I 阻止開發人員在其應用程序構建期間執行任意的 yum 安裝類型操作，避免過度安裝太多多餘的套件導致映像檔太過肥大。
- 可修補性 - 如果由於安全問題而基礎映像檔需要修補，則 S2I 您在重建所有應用程序的映像檔中就完成。
- 生態系統 - S2I 鼓勵共享映像檔的生態系統，您可以在其中找到並利用最佳實踐的方式來為你的應用程序建構映像檔。


本文會介紹如何透過 S2I 創建一個簡單的 nginx 服務映像檔。


安裝 S2I
----------

-----

在 MAC 可以透過 Homebrew 來安裝

```
$ brew install source-to-image
```

透過 S2I 來建立映像檔
--------------------

-----

我們可以利用 `s2i` 命令從 git repository 獲取程式碼以直接在 Builder Image 上構建我們的映像檔，您無需從 git repository 下載程式碼，也不用撰寫 Dockerfile。 對於一般的應用程序，這可是大大節省了時間和精力啊！

使用方法:

1) 從遠端 Git Repository 抓取原始碼來建構映像檔。

```
$ s2i build <git-repo> <S2I Builder Image Repository> <Your New Image Name>
```

2) 抓取遠端 Git Repository 的某個資料夾內原始碼來建構映像檔。
```
$ s2i build <git-repo> --context-dir=<Path/To/Context> <S2I Builder Image Repository> <Your New Image Name>
```

3) 當然你也可以從本地的程式碼來建構映像檔，不一定要透過 git repository。

```
$ s2i build . <S2I Builder Image Repository> <Your New Image Name>
```

例如:

```
$ s2i build https://github.com/sclorg/nginx-container.git --context-dir=1.16/test/test-app/ centos/nginx-116-centos7 quay-eu-uat/application-images/test:1
Submodule 'common' (https://github.com/sclorg/container-common-scripts.git) registered for path 'common'
Cloning into '/private/var/folders/7z/k_5hdgqx3vq1qrtxrk3619rw0000gn/T/s2i364743917/upload/tmp/common'...
Submodule path 'common': checked out '91d4ac4ceb89c7bced5c7f5ec552dbb45d637e7d'
---> Installing application source
---> Copying nginx.conf configuration file...
'./nginx.conf' -> '/etc/opt/rh/rh-nginx116/nginx/nginx.conf'
---> Copying nginx configuration files...
'./nginx-cfg/default.conf' -> '/opt/app-root/etc/nginx.d/default.conf'
---> Copying nginx default server configuration files...
'./nginx-default-cfg/alias.conf' -> '/opt/app-root/etc/nginx.default.d/alias.conf'
---> Copying nginx start-hook scripts...
Build completed successfully
```

然後我們就可以把建構好的映像檔放到 Quay ，並在 OpenShift 使用。

```
$ docker push quay-eu-uat/application-images/test:1
```

結論
-----

-----

對於一般不太複雜的服務，我們都可以透過 Source-To-Image (S2I) 來建構我們的映像檔。 在撰寫 Dockerfile 之前，不仿先找找是否有現成可以用的 S2I Image 可以利用，以節省時間．當個高效率的懶人工程師。


參考資料
---------

-----

- https://www.openshift.com/blog/create-s2i-builder-image