
OpenShift 系列 1： 裝了OCP就想部署Pod? 沒這麼容易!
========================================

Apply the first pod
-------------------------------

In the previous [blog](), we have launch the local OpenShift cluster.
Let's apply the same YAML file we used in [minikube]() to create a new pod.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  labels:
    app: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    imagePullPolicy: Always
    ports:
    - containerPort: 80
```

```
$ oc login -u developer -p developer https://api.crc.testing:6443
Login successful.

$ oc project myproject
Now using project "myproject" on server "https://api.crc.testing:6443".

$ oc apply -f ./test.yaml
pod/nginx-pod created

```

and check the status.

```
$ oc get pods
NAME        READY   STATUS             RESTARTS   AGE
nginx-pod   0/1     CrashLoopBackOff   3          81s
```

OOps, the pods is keep CrashLoopBackOff. Why? Let's check the logs

```
$ oc logs nginx-pod
2020/09/22 12:32:45 [warn] 1#1: the "user" directive makes sense only if the master process runs with super-user privileges, ignored in /etc/nginx/nginx.conf:2
nginx: [warn] the "user" directive makes sense only if the master process runs with super-user privileges, ignored in /etc/nginx/nginx.conf:2
2020/09/22 12:32:45 [emerg] 1#1: mkdir() "/var/cache/nginx/client_temp" failed (13: Permission denied)
nginx: [emerg] mkdir() "/var/cache/nginx/client_temp" failed (13: Permission denied)

```

The user in the container has the permission issue? Why did not we get this issue when we apply the same YAML to kubernetes? Actually, it is all about ``Security Context Constraints`` policy.



Security Context Constraints
----------------------------

Similar to the way that RBAC resources control user access, administrators can use Security Context Constraints (SCCs) to control permissions for pods. You can use SCCs to define a set of conditions that a pod must run with in order to be accepted into the system.

The cluster contains eight default SCCs:

- anyuid
- hostaccess
- hostmount-anyuid
- hostnetwork
- node-exporter
- nonroot
- privileged
- restricted


By default, cluster administrators, nodes, and the build controller are granted access to the privileged SCC. All authenticated users are granted access to the **restricted** SCC.

With the restricted SCC , a pod must run as a user in a pre-allocated range of UIDs. This range is defined in the Project's annotation. If the UID of the default user in the container is not within this range, OpenShift will pick a UID with the range, and run the container with the UID instead of the default user. As a result, OpenShift make sure all pods are running as non-root user.


Build Custome Images To Fix the Permission Issue
------------------------------------------------

By default, OpenShift Enterprise runs containers using an arbitrarily assigned user ID. This provides additional security against processes escaping the container due to a container engine vulnerability and thereby achieving escalated permissions on the host node.

For an image to support running as an arbitrary user, directories and files that may be written to by processes in the image should be owned by the root group and be read/writable by that group. Files to be executed should also have group execute permissions.

Adding the following to your Dockerfile sets the directory and file permissions to allow users in the root group to access them in the built image:

```
RUN chgrp -R 0 /some/directory \
  && chmod -R g+rwX /some/directory
```

Because the container user is always a member of the root group, the container user can read and write these files. The root group does not have any special permissions (unlike the root user) so there are no security concerns with this arrangement.


Beside, users are not allowed to listen on priviliged ports in OpenShift. Use port 8080 to listen on. When you expose the service for the web server outside of OpenShift the external route will use port 80 by default anyway and ensure traffic is routed through to port 8080 of your web server internally.


Therefore, we need to implement a customized image to run nginx on OpenShfit. Following is the example:

```
FROM nginx:1.14.2


RUN chgrp -R 0 /etc/nginx/ /var/cache/nginx /var/run /var/log/nginx  && \ 
  chmod -R g+rwX /etc/nginx/ /var/cache/nginx /var/run /var/log/nginx

# users are not allowed to listen on priviliged ports
RUN sed -i.bak 's/listen\(.*\)80;/listen 8081;/' /etc/nginx/conf.d/default.conf

# comment user directive as master process is run as user in OpenShift random UID
RUN sed -i.bak 's/^user/#user/' /etc/nginx/nginx.conf

EXPOSE 8080
```

Build image and push to our local Quay registry.

```
$ docker build -t quay-eu-uat.windmill.local/application-images/test:1 .
$ docker push quay-eu-uat.windmill.local/application-images/test:1
```


Use Images from Quay
----------------------

1) Create robot account


![](images/04_OCP_INTRO/quay_1.png)


2) Give permission to the robot account

![](images/04_OCP_INTRO/quay_2.png)

3) Get the pull secret

![](images/04_OCP_INTRO/quay_3.png)
![](images/04_OCP_INTRO/quay_4.png)

4) Apply the pull secret

```
$ oc apply -f application-images-ocp-secret.yml
```


5) Update the YAML to pull images from Quay

```YAML
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  labels:
    app: nginx
spec:
  containers:
  - name: nginx
    image: quay-eu-uat.windmill.local/application-images/test:1
    imagePullPolicy: Always
    ports:
    - containerPort: 8080
  imagePullSecrets:
    - name: application-images-ocp-pull-secret
```

5) Check the pods is running.

```
$ oc get pods
NAME        READY   STATUS    RESTARTS   AGE
nginx-pod   1/1     Running   0          26m
```




To fix the ssh/scp/sftp issue in containers running in OpenShift
---------------

According to RedHat's suggestion: https://access.redhat.com/solutions/4665281


Step 1: Modification of an image using Dockerfile

```
RUN chmod g=u /etc/passwd
```

Step 1: setting up an ENTRYPOINT that runs a script. For example:


ENTRYPOINT.sh
```
#!/bin/bash
if ! whoami &> /dev/null; then
  if [ -w /etc/passwd ]; then
    echo "${USER_NAME:-default}:x:$(id -u):0:${USER_NAME:-default} user:${HOME}:/sbin/nologin" >> /etc/passwd
  fi
fi
exec "$@"
```

Dockerfile
```
ADD ENTRYPOINT.sh /opt/ENTRYPOINT.sh

RUN chmod g=u /etc/passwd && \
    chgrp -R 0 /opt/ENTRYPOINT.sh && \
    chmod -R g+rx /opt/ENTRYPOINT.sh

ENTRYPOINT ["/opt/ENTRYPOINT.sh"]
```






Conclusion
-----------
There are some security consideration when run as root in the container, please check https://americanexpress.io/do-not-run-dockerized-applications-as-root/ for more details about this issue.

Openshift, as the Enterprise level platform, always consider the security overweight accessibility. Because many images on Dockerhub are run as root, which vialate the default Security Context Constraints rule in OpenShift, it is very important to learn how to customize your Image so that your can reduce privileges in a container and make it run on OpenShift without any issue. With the security scanning feature of Quay, you can also reduce the ``Vulnerability`` inside your container. Anyway, which company is willing to run their application/service in risk?


Reference
---------
- https://americanexpress.io/do-not-run-dockerized-applications-as-root/
- https://docs.openshift.com/container-platform/4.5/authentication/managing-security-context-constraints.html 




