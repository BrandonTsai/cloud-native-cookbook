new-app, template
============================

create new application without implement YAML file.


new-app
-------

with the support of source-to-images (S2I) and built-in CI/CD tools, developer can very easiy to deploy their app on OpenShift with one commend `oc new-app /path/to/source/code`

OpenShift Container Platform automatically detects whether the Docker, Pipeline or Source build strategy should be used, and in the case of Source builds, detects an appropriate language builder image.


![](image/new-app-lang.png)

for example
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


But if your application language can not be detected by OpenShift, you can specify the S2I image that you want to use to build your customize app, `oc new-app S2I_Image_Repo~/path/to/source/code`


for example

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


Besides,  you can also use docker images directly by `oc new-app Image_Repo`

for example
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


this `oc new-app` command will create following resource. 
- `BuildConfig`: A BuildConfig is created for each source repository specified in the command line. The BuildConfig specifies the strategy to use, the source location, and the build output location. 
- `ImageStream`: For BuildConfig, two ImageStreams are usually created. One represents the input image. With Source builds, this is the builder image. With Docker builds, this is the FROM image. The second one represents the output image. If a container image was specified as input to new-app, then an image stream is created for that image as well. 
- `Deployment` and `Service`: to run your application.

you can check the status of these resource by `oc get all` command.

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


A template describes a set of objects that can be parameterized and processed to produce a list of objects for creation by OpenShift Container Platform. A template can be processed to create anything you have permission to create within a project, for example services, build configurations, and DeploymentConfigs. A template may also define a set of labels to apply to every object defined in the template.

You can create a list of objects from a template using the CLI or, if a template has been uploaded to your project or the global template library, using the web console.



### Creating an application by using the web console



### Creating objects from templates by using the CLI

You can list the available default Instant App and Quickstart templates with:


$ oc get templates -n openshift

oc process --parameters -n <project> <template_name>

OpenShift has many buildin Template. Developer can deploy application easily from Web UI or through CLI `oc process`




Clean up
--------

The drawback of 'oc new-app' and 'template' is that you can not clean all resource related to the application eaisy. you have to get all resouce and delete these resource one by one.


Conclusion
----------

I do not encourage developers and operators deploy production applications via `oc new-app` command as it does not provide accurate resource requests and limits. this should be used in development environment to test new features without implement any YAML file.



Reference
----------

- https://docs.openshift.com/container-platform/4.5/openshift_images/using-templates.html