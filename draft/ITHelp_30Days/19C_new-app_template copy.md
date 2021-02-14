免 YAML 部署 App 到 OpenShift: new-app 跟 Template 淺談
============================

在一般 K8S 平台，一般你必須要撰寫 YAML 檔或 Helm 來部署 App。

但是在 OpenShift 你可以透過 `oc new-app` 跟 `Template` 這內建的兩個功能，快速的部署一般常見簡易的服務。


new-app
-------

在 Source-To-Images（S2I）和內建的 CI/CD 工具的支持下，開發人員可以非常輕鬆地透過指令 `oc new-app /path/to/source/code` 將其應用程序部署到 OpenShift 上。

OpenShift 會自動檢測程式碼的語言，並使用適當的 S2I 映像檔來快速建構客製化的映像檔，然後再將其部署到你的 Project 內。以下是支援的程式語言。

![](images/04_OCP/new-app-lang.png)


範例：

```
$ oc new-app https://github.com/sclorg/cakephp-ex
--> Found image 988e5d4 (2 months old) in image stream "openshift/php" under tag "7.3" for "php"

    Apache 2.4 with PHP 7.3 
    ----------------------- 
    PHP 7.3 available as container is a base platform for building and running various PHP 7.3 applications and frameworks. PHP is an HTML-embedded scripting language. PHP attempts to make it easy for developers to write dynamically generated web pages. PHP also offers built-in database integration for several commercial and non-commercial database management systems, so writing a database-enabled webpage with PHP is fairly simple. The most common use of PHP coding is probably as a replacement for CGI scripts.

    Tags: builder, php, php73, rh-php73

    * The source repository appears to match: php
    * A source build using source code from https://github.com/sclorg/cakephp-ex will be created
      * The resulting image will be pushed to image stream tag "cakephp-ex:latest"
      * Use 'oc start-build' to trigger a new build

--> Creating resources ...
    imagestream.image.openshift.io "cakephp-ex" created
    buildconfig.build.openshift.io "cakephp-ex" created
    deployment.apps "cakephp-ex" created
    service "cakephp-ex" created
--> Success
    Build scheduled, use 'oc logs -f bc/cakephp-ex' to track its progress.
    Application is not exposed. You can expose services to the outside world by executing one or more of the commands below:
     'oc expose svc/cakephp-ex' 
    Run 'oc status' to view your app.
```


如果OpenShift無法檢測到您的應用程序語言，則可以透過 `oc new-app S2I_Image_Repo~/path/to/source/code` 指定要用於構建應用程序映像檔的S2I基礎映像檔．


例如：

```
$ oc new-app centos/nginx-116-centos7~https://github.com/sclorg/nginx-container.git --context-dir=1.16/test/test-app/
--> Found container image 28684f2 (2 weeks old) from Docker Hub for "centos/nginx-116-centos7"

    Nginx 1.16 
    ---------- 
    Nginx is a web server and a reverse proxy server for HTTP, SMTP, POP3 and IMAP protocols, with a strong focus on high concurrency, performance and low memory usage. The container image provides a containerized packaging of the nginx 1.16 daemon. The image can be used as a base image for other applications based on nginx 1.16 web server. Nginx server image can be extended using source-to-image tool.

    Tags: builder, nginx, rh-nginx116

    * An image stream tag will be created as "nginx-116-centos7:latest" that will track the source image
    * A source build using source code from https://github.com/sclorg/nginx-container.git will be created
      * The resulting image will be pushed to image stream tag "nginx-container:latest"
      * Every time "nginx-116-centos7:latest" changes a new build will be triggered

--> Creating resources ...
    imagestream.image.openshift.io "nginx-116-centos7" created
    imagestream.image.openshift.io "nginx-container" created
    buildconfig.build.openshift.io "nginx-container" created
    deployment.apps "nginx-container" created
    service "nginx-container" created
--> Success
    Build scheduled, use 'oc logs -f bc/nginx-container' to track its progress.
    Application is not exposed. You can expose services to the outside world by executing one or more of the commands below:
     'oc expose svc/nginx-container' 
    Run 'oc status' to view your app.
```


當然,  你也可以直接使用已經建構好的映像檔，透過 `oc new-app Image_Repo` 指令來自動產生 Deployment 跟 Service，並部署到OpenShift 上。

譬如：

```
$ oc new-app quay.io/brandon_tsai/testlab:1
--> Found container image 3d97f35 (10 days old) from quay.io for "quay.io/brandon_tsai/testlab:1"

    quay.io/brandon_tsai/testlab:1 
    ------------------------------ 
    Nginx is a web server and a reverse proxy server for HTTP, SMTP, POP3 and IMAP protocols, with a strong focus on high concurrency, performance and low memory usage. The container image provides a containerized packaging of the nginx 1.16 daemon. The image can be used as a base image for other applications based on nginx 1.16 web server. Nginx server image can be extended using source-to-image tool.

    Tags: builder, nginx, rh-nginx116

    * An image stream tag will be created as "testlab:1" that will track this image

--> Creating resources ...
    imagestream.image.openshift.io "testlab" created
    deployment.apps "testlab" created
    service "testlab" created
--> Success
    Application is not exposed. You can expose services to the outside world by executing one or more of the commands below:
     'oc expose svc/testlab' 
    Run 'oc status' to view your app.
```


`oc new-app` 指令通常會在你的專案空間（Project）內自動產生以下資源： 

- `BuildConfig`: 為命令行中指定的每個原始碼指定要使用的策略，原始碼位置和構建後輸出的位置。它類似於 Pipeline，是 OpenShift 用來支援 CI/CD 的工具之一。
- `ImageStream`: 對於BuildConfig，通常有兩個 ImageStream會被建立。 一個代表輸入的S2I 映像檔。 第二個代表建構好的應用程序映像檔。 
- `Deployment` and `Service`: 就是用來跑你的應用程序。


透過 `oc new-app` 指令部署上去後，你可以透過 `oc get all` 來檢查資源目前的狀態。

```
$ oc get all
NAME                              READY   STATUS      RESTARTS   AGE
pod/cakephp-ex-1-build            0/1     Completed   0          7m15s
pod/cakephp-ex-5486bcb578-shhpg   1/1     Running     0          113s

NAME                 TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)             AGE
service/cakephp-ex   ClusterIP   172.25.252.30   <none>        8080/TCP,8443/TCP   7m15s

NAME                         READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/cakephp-ex   1/1     1            1           7m15s

NAME                                    DESIRED   CURRENT   READY   AGE
replicaset.apps/cakephp-ex-5486bcb578   1         1         1       113s
replicaset.apps/cakephp-ex-f9687bdc     0         0         0       7m15s

NAME                                        TYPE     FROM   LATEST
buildconfig.build.openshift.io/cakephp-ex   Source   Git    1

NAME                                    TYPE     FROM          STATUS     STARTED         DURATION
build.build.openshift.io/cakephp-ex-1   Source   Git@377fe8f   Complete   7 minutes ago   5m23s

NAME                                        IMAGE REPOSITORY                                                         TAGS     UPDATED
imagestream.image.openshift.io/cakephp-ex   default-route-openshift-image-registry.apps-crc.testing/uat/cakephp-ex   latest   About a minute ago

```


Template
----------

“模板(Template)” 描述了一組可參數化和處理的資源，以用來快速針對相似的應用程序，快速生成及部署到不同的專案空間（Project）。他跟 Ｈelm 很像，但是他是 OpenShift 內建的功能，你不用再做另外的設定，此外你可以針對每個參數提供描述跟規範，這讓使用者可以很快速在使用者介面直接透過“模板(Template)” 部署新服務，而不用再另外閱讀說明文件。



### 範例： 在使用者介面，透過 Template 建立 Apache 服務。

(1) 登入使用者介面後，在你的專案空間下，選擇從 Catalog 新增應用程序。

![](images/04_OCP/t1.png)

(2) 然後你就可以看到很多內建的 Template。

![](images/04_OCP/t2.png)

(3) 選擇 "Apache HTTP Server" 做測試。


![](images/04_OCP/t3.png)

(4) 輸入參數後按新增。

![](images/04_OCP/t4.png)


(5) 然後就可以檢查應用程序的狀態啦，一但 Pod 變成 Running 的狀態，表示該應用程序 已經正常可以運作了。

![](images/04_OCP/t5.png)




### 範例： 從指令透過 Template 建立 Apache 服務。


(1) 可以先透過 `oc get templates -n openshift` 列出所有預設可用的 templates

```
$ oc get templates -n openshift | grep "Apache"
httpd-example                                   An example Apache HTTP Server (httpd) application that serves static content....   9 (3 blank)       5
```

(2) 透過`oc process --parameters openshift//<Template Name>`，看有哪些參數需要設定

```
$ oc process --parameters openshift//httpd-example
NAME                     DESCRIPTION                                                                                               GENERATOR           VALUE
NAME                     The name assigned to all of the frontend objects defined in this template.                                                    httpd-example
NAMESPACE                The OpenShift Namespace where the ImageStream resides.                                                                        openshift
MEMORY_LIMIT             Maximum amount of memory the container can use.                                                                               512Mi
SOURCE_REPOSITORY_URL    The URL of the repository with your application source code.                                                                  https://github.com/sclorg/httpd-ex.git
SOURCE_REPOSITORY_REF    Set this to a branch name, tag or other ref of your repository if you are not using the default branch.                       
CONTEXT_DIR              Set this to the relative path to your project if it is not in the root of your repository.                                    
APPLICATION_DOMAIN       The exposed hostname that will route to the httpd service, if left blank a value will be defaulted.                           
GITHUB_WEBHOOK_SECRET    Github trigger secret.  A difficult to guess string encoded as part of the webhook URL.  Not encrypted.   expression          [a-zA-Z0-9]{40}
GENERIC_WEBHOOK_SECRET   A secret string used to configure the Generic webhook.                                                    expression          [a-zA-Z0-9]{40}
```


(3) 透過 `oc process <template_name> PARM1=VALUE1 PARM2=VALUE2` 來產生可用於部署應用程序的YAML檔，並部署到目標的專案空間內。

```
$ oc process openshift//httpd-example MEMORY_LIMIT=128Mi  | oc create -n myproject -f -
service/httpd-example created
route.route.openshift.io/httpd-example created
imagestream.image.openshift.io/httpd-example created
buildconfig.build.openshift.io/httpd-example created
deploymentconfig.apps.openshift.io/httpd-example created
```

(4) 透過 `oc status` 來檢查部署狀態。

```
$ oc status
In project myproject on server https://api.crc.testing:6443

http://httpd-example-myproject.apps-crc.testing (svc/httpd-example)
  dc/httpd-example deploys istag/httpd-example:latest <-
    bc/httpd-example source builds https://github.com/sclorg/httpd-ex.git on openshift/httpd:2.4 
      build #1 running for 28 seconds - 72f17c2: Merge pull request #24 from multi-arch/master (Honza Horak <hhorak@redhat.com>)
    deployment #1 waiting on image or update

View details with 'oc describe <resource>/<name>' or list everything with 'oc get all'.
```


Clean up
--------

'oc new-app' 和 'template' 的缺點是您無法一鍵清除與應用程序有關的所有資源。 您必須手動一一刪除這些資源。

```
$ oc get all -o NAME --no-headers | xargs oc delete
pod "httpd-example-1-build" deleted
pod "httpd-example-1-deploy" deleted
pod "httpd-example-1-q525s" deleted
replicationcontroller "httpd-example-1" deleted
service "httpd-example" deleted
deploymentconfig.apps.openshift.io "httpd-example" deleted
buildconfig.build.openshift.io "httpd-example" deleted
build.build.openshift.io "httpd-example-1" deleted
imagestream.image.openshift.io "httpd-example" deleted
route.route.openshift.io "httpd-example" deleted
```

結論
----------

我不鼓勵開發人員和操作員通過`oc new-app`命令部署正式環境的應用程序，因為它沒有提供準確的資源請求/限制，以及 HealthCheck。 他只應該在開發環境中使用它來測試新功能，而免去修改任何YAML文件或Dockerfile的煩惱。

Openshift 提供 98 個可直接使用的 Template，這在部署一般常見的服務，例如資料庫或 Redis 時非常有用，因為 [Helm Hub](https://hub.helm.sh/) 很多會因為權限問題，無法直接部署在 Openshift 上，此時，我們就可以透過內建的 Template 來處理。

然而，雖然 `new-app` 跟 `Template` 都可以快速的部署應用程序，但他們都無法針對該應用程序產生的資源進行版本更新或移除，此時用 Helm 來管理應用程序明顯是較佳的選擇。 我會建議建立一個開發用的專案空間（Project）, 然後在此專案空間利用`new-app` 跟 `Template` 來快速產生要開發及測試的環境。


參考資料
----------

- https://docs.openshift.com/container-platform/4.5/applications/application_life_cycle_management/creating-applications-using-cli.html
- https://docs.openshift.com/container-platform/4.5/openshift_images/using-templates.html