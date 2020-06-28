Podman - Part I:  Introduction
===============================

Introduction
-------------

As you may have noticed, Red Hat replaces the Docker Daemon with CRI-O/Podman since RHEL 8.
So what is Podman? According to the definition at [Podman official website](https://podman.io/),

> Podman is a daemonless container engine for developing, managing, and running OCI Containers on your Linux System. Containers can either be run as root or in rootless mode.

Why does Red Hat want to get rid of the Docker Daemon? This is because there are few problems with running Docker with Docker Daemon

- Single point of failure issue, once the daemon died, all containers died.
- This daemon process owned all the child processes for the running containers.
- All Docker operations had to be conducted by a user with the same full root authority.
- Building containers could lead to security vulnerabilities.

So Podman solves the above issues by directly interacting with Image registry, containers and image storage instead of work through a daemon. And the rootless mode allows a user to run containers without the full root authority.
Besides, it also provides a Docker compatible command-line experience enabling users to pull, build, push and run containers.


#### Docker vs Podman

![](images/Docker_vs_Podman.png)

Podman interacts with Linux kernel to manage containers through the runC container runtime process instead of a daemon
. The [buildah](https://buildah.io/)  utility is used to replace Docker build as the container images build tool and Docker push is replaced by [skopeo](https://github.com/containers/skopeo) for moving container images between registries and container engines.


Installation and Setup
------------------------

For RHEL7, subscribe rhel-7-server-extras-rpms yum repository and then enable Extras channel and install Podman.

```
sudo subscription-manager repos --enable=rhel-7-server-extras-rpms
sudo yum-config-manager --enable rhel-7-server-extras-rpms
sudo yum -y install podman
```

#### Rootless mode?
Podman supports rootless mode, for more details to set up rootless mode on Redhat 7, please refer:
https://www.redhat.com/en/blog/preview-running-containers-without-root-rhel-76

This article will focus on the basic usage of Podman.

Basic Usage
------------

Most Podman commands are similar to Docker commands. If you’ve used the Docker CLI, you will be quite familiar with Podman.

```bash
# Pull mage
sudo podman pull nginx

# List images
sudo podman images

# Run container
sudo podman run -dt -p 8081:80/tcp -v /opt/http:/usr/share/nginx/html:ro --name hello-nginx nginx
```


Checkpoint and Restore
------------------------

![](images/container-cli.png)

If you are running containers with [tmpfs volume](https://docs.docker.com/storage/tmpfs/), then export/import can not be the backup solution for that container because `export` does not back up the memory content. The files in the tmpfs volume will be lost when you `import` the container from the tar file. You can use instead `checkpoint/restore`. For Docker, you will need to turn on 'experimental features to enable this. Podman can use these features directly without doing any change.

#### Run container that supports checkpoint

The `criu` package is required to do checkpoint/restore. And you have to add --security-opt="seccomp=unconfined" when running a container on RHEL because CRIU cannot correctly handle seccomp on RHEL7.

```
sudo yum install -y criu
sudo podman run -dt --tmpfs /tmp -v /opt/http:/usr/share/nginx/html --security-opt="seccomp=unconfined" --name hello-nginx nginx
```

#### Create file in /tmp/ folder
```
sudo podman exec -it hello-nginx touch /tmp/test-01
```

#### Create checkpoint and export as a file.
```
sudo podman container checkpoint --leave-running --export=/tmp/backup.tar hello-nginx
```

#### Restore from file
```
sudo podman stop hello-nginx
sudo podman rm hello-nginx
sudo podman container restore --import=/tmp/backup.tar
```

#### Verify data does not lost in /tmp/ folder

```
$ sudo podman exec -it hello-nginx ls /tmp/test-01
/tmp/test-01
```

Conclusion
-----------

As the replacement of Docker, Podman provides the same developer experience as Docker while doing things in a slightly more secure way in the background. You could alias Docker with Podman and never notice that there is a completely different tool for managing your local containers. Besides, with the daemonless design and the rootless mode, Podman is more isolated and secure than Docker. You should consider using Podman instead of installing Docker-ce on your local machine.


