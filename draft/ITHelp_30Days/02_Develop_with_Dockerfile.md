Introduction
------------
If you can not find suitable images from image registry for your application, instead of mounting the static code of your application to a running container, you should always implement a Dockerfile and build your own images with the Dockerfile. This blog will first introduce the elemental format of Dockerfile, and will then propose some advanced usage for different scenarios.


Build Images with Dockerfile
----------------------------

Following is the basic Dockerfile example, it will build a customize Nodejs application images base on image registry.access.redhat.com/ubi7/nodejs-12.

For more details about Dockerfile format, you can refer [Dockerfile reference](https://docs.docker.com/engine/reference/builder/) for more usage of Dockerfile.


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

The ``FROM`` instruction initializes a new build stage and sets the Base Image for subsequent instructions. As such, a valid Dockerfile must start with a ``FROM`` instruction. The image can be any valid image – it is especially easy to start by pulling an image from the Public Repositories.


The most common base image is Alpine. But you can also use [Red Hat UBI images](https://www.redhat.com/en/blog/introducing-red-hat-universal-base-image) if your application was deploied on Red Hat System. You can find numerous images that are built based on the Red Hat UBI image on [Red Har Registry](https://catalog.redhat.com/software/containers/explore), and you can use one of them as your base image.



### COPY & ADD

```
ADD [--chown=<user>:<group>] <src>... <dest>
COPY [--chown=<user>:<group>] <src>... <dest>
```

The ``ADD`` and ``COPY`` instruction copies new files, directories or remote file URLs from ``<src>`` and adds them to the filesystem of the image at the path ``<dest>``.

``COPY`` only supports the basic copying of local files into the container, while ``ADD`` has some features like local-only tar extraction and remote URL support. for example:

```
ADD http://example.com/font.js /opt/
ADD my_big_lib.tar.gz /var/lib/myapp
```

For other items (files, directories) that do not require ADD’s tar auto-extraction capability, you should always use ``COPY``. Using ``ADD`` to fetch packages/files from remote URLs is discouraged; you should use ``curl`` or ``wget``.



### WORKDIR

The ``WORKDIR`` instruction sets the working directory for any RUN, CMD, ENTRYPOINT, COPY and ADD instructions that follow it in the Dockerfile.

### RUN

```
# shell form, the command is run by /bin/sh -c
RUN <command>
# exec form
RUN ["executable", "param1", "param2"]
```

The ``RUN`` instruction will execute any commands in a new layer on top of the current image and commit the results.

making sure the most general steps, and the longest are first, that will then cached, allowing you to fiddle with the last lines of your Dockerfile (the most specific commands) while having a quick rebuild time.

### USER

```
USER <user>[:<group>]
USER <UID>[:<GID>]
```

The ``USER`` instruction sets the user name (or UID) and optionally the user group (or GID) to use when running the image and for any RUN, CMD and ENTRYPOINT instructions that follow it in the Dockerfile.


### EXPOSE

```
EXPOSE <port>[/<protocol>]
```

The ``EXPOSE`` instruction informs Docker that the container listens on the specified network ports at runtime. You can specify whether the port listens on TCP or UDP, and the default is TCP if the protocol is not specified.



### ENTRYPOINT & CMD

```
# exec form, this will use "exec" as default ENTRYPOINT
CMD ["executable","param1","param2"]

# shell form, the default ENTRYPOINT is "/bin/sh -c"
CMD command param1 param2
```

An ``ENTRYPOINT`` allows you to configure a container that will run as an executable.

An ``CMD`` is to provide defaults for an executing container.

There can only be one ``CMD`` instruction in a Dockerfile. If you list more than one CMD then only the last CMD will take effect.


``CMD`` can be used to provide default arguments for the ``ENTRYPOINT`` instruction, In this case, both the ``CMD`` and ``ENTRYPOINT ``instructions should be specified with the JSON array format. For example:

```
# as default parameters to ENTRYPOINT
ENTRYPOINT ["top", "-b"]
CMD ["-c"]
```


### Build Image

To build the image, you can just simply run

```
docker build  -t my-nodejs-app:0.1.0 .
```


Build with argument
--------------------

If you do not want to implement different Dockerfile for different environments, you can use "ARG" instruction.

The ARG instruction defines a variable that users can pass at build-time to the builder with the docker build command using the ``--build-arg <varname>=<value>`` flag.

> It is not recommended to use build-time variables for passing secrets like github keys, user credentials etc. Build-time variable values are visible to any user of the image with the docker history command. Refer to the “[build images with BuildKit](https://docs.docker.com/develop/develop-images/build_enhancements/#new-docker-build-secret-information)” section to learn about secure ways to use secrets when building images.


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


Build docker images with build argument

```
docker build --build-arg APP_ENV=production -t my-app:0.1.0 .
```



Multi-stage Build
------------------

One of the most challenging things about building images is keeping the image size down. Each instruction in the Dockerfile adds a layer to the image, and you need to remember to clean up any artifacts you don’t need before moving on to the next layer. To write a really efficient Dockerfile, you have traditionally needed to employ shell tricks and other logic to keep the layers as small as possible and to ensure that each layer has the artifacts it needs from the previous layer and nothing else.

> This Multi-stage Build requiring Docker 17.05 or higher version. If you are working on RHEL7 system, I strongly recommend you using Podman instead of Docker.

Following is the example for Multi-stage build

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

To build the image

```
docker build --build-arg APP_ENV=production -t my-app:0.1.0 .
```

Conclusion
-----------

This essay has cover the Dockerfile format and advanced usage. I hope this gives you a basic sense of Dockerizing an application.



