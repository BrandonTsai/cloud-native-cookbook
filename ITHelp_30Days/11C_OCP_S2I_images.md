

OpenShift 系列 2： 懶人救星 - Source-To-Image (S2I)
=================================================

Source-To-Image (S2I) 淺談
---------------------------

If your application is not complicate, instead of implement Dockerfile and work hard to figure out how to fix the permission isse, another way is using the Source-To-Image (S2I) image.
https://github.com/openshift/source-to-image

The main reasons one might be interested in using source builds are:

- Speed - with S2I, the assemble process can perform a large number of complex operations without creating a new layer at each step, resulting in a fast process.
- Patchability - S2I allows you to rebuild the application consistently if an underlying image needs a patch due to a security issue.
- User efficiency - S2I prevents developers from performing arbitrary yum install type operations during their application build, which results in slow development iteration.
- Ecosystem - S2I encourages a shared ecosystem of images where you can leverage best practices for your applications.
This article is about creating a simple S2I builder image.

S2I generates a new Docker image using source code and a  builder Docker image. The S2I project includes a number of commonly used Docker builder images, such as Python or Ruby, you can also extend S2I with your own custom scripts.


安裝 S2I 
----------

在 MAC 可以透過 Homebrew 來安裝

```
$ brew install source-to-image
```

透過 S2I 來建立映像檔
--------------------

We can use `s2i` command to get context from git repository to build our custom image based on the S2I image directly, you do not need to download the git repo and implement the Dockerfile. What a life saver for simple application!

使用方法:

**Build a Docker image from a remote Git repository**

```
$ s2i build <git-repo> <S2I Image Repository> <Your New Image Name>
```

**Build a Docker image from a remote Git repository which context are in particular folder**
```
$ s2i build <git-repo> --context-dir=<Path/To/Context> <S2I Image Repository> <Your New Image Name>
```

**Build from a local directory.  If this directory is a git repo then the current commit will be built.**

```
$ s2i build . <S2I Image Repository> <Your New Image Name>
```

For example:

```
$ s2i build https://github.com/sclorg/nginx-container.git --context-dir=1.16/test/test-app/ centos/nginx-116-centos7 quay-eu-uat.windmill.local/application-images/test:1
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

Push to Quay and so that we can use it in our OpenShift YAML file

```
$ docker push quay-eu-uat.windmill.local/application-images/test:1
```

結論
-----

S2I Images is very useful  


Reference
---------

- https://www.openshift.com/blog/create-s2i-builder-image