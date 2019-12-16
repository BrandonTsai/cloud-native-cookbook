Podman - Part II :  Run a Pod 
==============================

Pods
-----

The Pod concept was introduced by [Kubernetes](https://kubernetes.io/docs/concepts/workloads/pods/pod/). A pods is a group of containers that operate together. podman use the similar concept to manage a group of containers on local server. All containers inside the pod share the same network namespace, so they can easily talk to each other over localhost without the need to export any extra ports. You can refer [Podman: Managing pods and containers in a local container runtime](https://developers.redhat.com/blog/2019/01/15/podman-managing-containers-pods/) for more details about the technicals that podman used. In this artical, we will focus on how to run and manage pods on local server.

### Create pod manually

The first thing to be done is the creation of a new pod.
```
# sudo podman pod create -n my-app -p 8081:80
```

And then add a container to a pod
```
sudo podman run -dt --pod my-app -v /opt/http:/usr/share/nginx/html:ro --security-opt="seccomp=unconfined" --name hello-nginx nginx
```

Notices that you can not run a container which binding port to a container that run in a pod.
you have to bind the port to the pod instead. And there is an issue when you try to export multiple port in a pods.

you can list all pods by `podman pod ps`
```
# sudo podman pod ps
POD ID         NAME     STATUS    CREATED          # OF CONTAINERS   INFRA ID
75d943416fc8   my-app   Created   26 minutes ago   1                 30138c8d0d1c
```

If you stop a pod, all containers in the pod will be stopped as well.
```
$ sudo podman pod stop my-app
a2edfd1095760b1e2946271184743cce6f621665878b618ddc83d73b295070ba
$ sudo podman ps -a
CONTAINER ID  IMAGE                           COMMAND               CREATED        STATUS                    PORTS                 NAMES
cacdc75990b0  docker.io/library/nginx:latest  nginx -g daemon o...  2 minutes ago  Exited (0) 7 seconds ago  0.0.0.0:8082->80/tcp  hello-nginx
4dce350e01cf  k8s.gcr.io/pause:3.1                                  3 minutes ago  Exited (0) 7 seconds ago  0.0.0.0:8082->80/tcp  a2edfd109576-infra
```

Similarly, start a pod will start all containers in the pod
```
$ sudo podman pod start my-app
a2edfd1095760b1e2946271184743cce6f621665878b618ddc83d73b295070ba
$ sudo podman ps -a
CONTAINER ID  IMAGE                           COMMAND               CREATED        STATUS            PORTS                 NAMES
cacdc75990b0  docker.io/library/nginx:latest  nginx -g daemon o...  4 minutes ago  Up 5 seconds ago  0.0.0.0:8082->80/tcp  hello-nginx
4dce350e01cf  k8s.gcr.io/pause:3.1
```


### Create Pod by kubernetes style yaml file.

https://mkdev.me/en/posts/dockerless-part-3-moving-development-environment-to-containers-with-podman

podman support setting a pod via Kubernetes-compatible pod definition yaml file.
And you can mount volume by useing [hostPath](https://kubernetes.io/docs/concepts/storage/volumes/#hostpath)


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

To create a new pod with yaml file
```
$ sudo podman play kube ./my-app.yaml
```


Check all containers are running
```
$ podman ps -a
CONTAINER ID  IMAGE                           COMMAND               CREATED        STATUS                        PORTS                             NAMES
2268e5ab9b61  nginx                           nginx -g daemon o...  2 minutes ago  Up 2 minutes ago              0.0.0.0:8001-8002->8001-8002/tcp  ng2
19dba831eeae  nginx                           nginx -g daemon o...  2 minutes ago  Up 2 minutes ago              0.0.0.0:8001-8002->8001-8002/tcp  ng1
42c150972ddb  k8s.gcr.io/pause:3.1                                  2 minutes ago  Up 2 minutes ago              0.0.0.0:8001-8002->8001-8002/tcp  4ae6b24effb5-infra
```

> Notices that all the containers in a pod will share the same local ip 127.0.0.1, they must running on different ports, otherwise some container will fail to start due to the port conflict.


Conclution
-----------

The ability for Podman to handle Kubernetes-compatible pod deployment is a clear differentiator to other container runtimes. For the kubernetes users, they should feel familiar to implement the YAML file to manage a group of containers locally.

However, compare to docker-compose, Podman Pod can not be used to build multiple images at the same time. There is a third party tool [podman-compose](https://github.com/muayyad-alsadi/podman-compose) that might bring this functionality. But I would sugget implement a script for building images and use Podman Pod for containers management.
