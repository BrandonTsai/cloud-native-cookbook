---
title: "#2 Openshift: Service, Route and Deployment Strategies "
author: Brandon Tsai
---


Networking Concept: Service & Route
---------


In the traditional VM based environment, if you want to expose a new service, you need to assign an IP to your new VM and add it to the DNS server, set up Nginx to enable TLS certificate, and update the HAProxy configuration and reload the HAProxy service to handle the load balance. It is OK if you are in a small startup company, you do not need to settle the networking frequently. However, if you warking in a global organization, which many new services are launched every day, you might find that it is extremely annoying.

The networking model in Kubernetes can help you to reduce this pain by addressing the following concerns:

- Containers within a Pod can communicate via loopback (127.0.0.1).
- Different Pods can communicate via the `Service` resource.
- `Ingress Controller` and the related `Ingress` Resource can expose a service to be reachable from external nodes outside the cluster.


For the Kubernetes cluster, the cluster administrator has to set up the Ingress Controllers manually before using the Ingress resource.
There are many Ingress Controller open sources that can be used in the Kubernetes cluster. 
Cluster administrators need to investigate and test these open sources in advance.
For Openshift users, Openshift already provides a default built-in solution for the external traffic called `Route`.
The developers can use `Route` to expose a service directly without any complicated setup.

Following is the concept of the Openshift Networking Model:

![](19_images/network_1.png)



OpenShift Route
--------------------

Openshift Route is similar to Kubernetes Ingress, but it has additional capabilities such as splitting traffic between multiple backends, sticky sessions, etc. When a Route object is created on OpenShift, it gets picked up by the built-in HAProxy load balancer in order to expose the requested service and make it externally available with the given configuration.

![](19_images/network_ingress_route.png)


**Route without SSL Example**

```
$ oc expose svc Nginx
route.route.openshift.io/nginx exposed
```

```
apiVersion: v1
kind: Route
metadata:
  name: nginx
spec:
  to:
    kind: Service
    name: nginx
```

Apply YAML file
```
$ oc apply -f route.yml 
route.route.openshift.io/nginx created

$ oc get route
NAME    HOST/PORT                          PATH   SERVICES   PORT    TERMINATION   WILDCARD
nginx   nginx-myproject.apps-crc.testing          nginx      <all>                 None
```

If you do not specify the hostname, OpenShift will generate a host URL `<route-name>-<namespace>.<external-address>`
You can also specify a custom hostname, but make sure the hostname is in a subdomain of the `<external-address>`, for example:

```
apiVersion: v1
kind: Route
metadata:
  name: nginx
spec:
  host: nginx-uat.apps-crc.testing
  to:
    kind: Service
    name: nginx
```


Verify that we can access the host via Brower or curl

```
$ curl http://nginx-uat.apps-crc.testing
<html>
<head>
	<title>Test NGINX passed</title>
</head>
<body>
<h1>NGINX is working</h1>
</body>
</html>
```

**Route with SSL termination example**

The simplest way to create a SSL route with edge termination is

```
$ oc create route edge --service=nginx --hostname=nginx-uat.apps-crc.testing
route.route.openshift.io/nginx created
```

or via YAML

```
apiVersion: v1
kind: Route
metadata:
  name: nginx
spec:
  host: nginx-uat.apps-crc.testing
  to:
    kind: Service
    name: nginx
  tls:
    termination: edge
```


It will use a cluster default wildcard certificate.
If you want to use customize certificate, you can use command `oc create route edge --service=nginx --hostname=nginx-uat.apps-crc.testing  --key=nginx-uat.key --cert=nginx-uat.crt`


If you want to redirect HTTP connection to HTTPS connection automatically, you can add parameter `--insecure-policy=Redirect`

```
$ curl -k https://nginx-myproject.apps-crc.testing
<html>
<head>
	<title>Test NGINX passed</title>
</head>
<body>
<h1>NGINX is working</h1>
</body>
</html>
```


**Passthrough Termination**

Let's update the Nginx configuration in the container and enable the 8443 port in the deployment

The Nginx config:

```
server {
        listen       8443 ssl http2 default_server;
        listen       [::]:8443 ssl http2 default_server;
        server_name  _;
        root         /opt/app-root/src;

        ssl_certificate "/opt/app-root/tls/tls.crt";
        ssl_certificate_key "/opt/app-root/tls/tls.key";


        location / {
        }

        error_page 404 /404.html;
            location = /40x.html {
        }

        error_page 500 502 503 504 /50x.html;
            location = /50x.html {
        }
    }
```

Update the deployment YAML file.

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: quay.io/brandon_tsai/testlab:1
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
        - containerPort: 8443
```

Update and apply service config

```
apiVersion: v1
kind: Service
metadata:
  name: nginx
spec:
  selector:
    app: nginx
  ports:
    - name: port-http
      protocol: TCP
      port: 80
      targetPort: 8080
    - name: port-https
      protocol: TCP
      port: 443
      targetPort: 8443
```


Create route via command

```
$ oc create route passthrough --service=nginx --hostname=nginx-uat.apps-crc.testing --port=port-https
route.route.openshift.io/nginx created
```

or via a YAML file

```
apiVersion: route.openshift.io/v1
kind: Route
metadata:
  name: nginx
spec:
  host: nginx-uat.apps-crc.testing
  port:
    targetPort: port-https
  tls:
    termination: passthrough
  to:
    kind: Service
    name: nginx
status: {}
```



**Re-encryption Termination**

The router terminates TLS with a certificate and then re-encrypts the connection with another different certificate. As a result, the full path of the connection is encrypted. The method is more secure because the user can not get the internal certificate from the browser directly.

Using the same deployment and service above, but this time we change it to re-encrypts.
The communication between external clients and Openshift via the cluster default certificate.
The connection inside OpenShift is encrypted by our custom certificate.

```
oc create route reencrypt --service=nginx --hostname=nginx-uat.apps-crc.testing --port=port-https  --dest-ca-cert=rootCA.crt
```

```
apiVersion: route.openshift.io/v1
kind: Route
metadata:
  name: nginx
spec:
  host: nginx-uat.apps-crc.testing
  port:
    targetPort: port-https
  tls:
    destinationCACertificate: |
      -----BEGIN CERTIFICATE-----
      MIIFRDCCAywCCQCe2OOaDBEvLjANBgkqhkiG9w0BAQsFADBkMQswCQYDVQQGEwJB
      VTEMMAoGA1UECAwDTlNXMQ8wDQYDVQQHDAZTeWRuZXkxEjAQBgNVBAoMCURhcnVt
      YXRpYzELMAkGA1UECwwCSVQxFTATBgNVBAMMDGJyYW5kb24udGVzdDAeFw0yMDA5
      MjcwNzU5NDBaFw0yMzA3MTgwNzU5NDBaMGQxCzAJBgNVBAYTAkFVMQwwCgYDVQQI
      DANOU1cxDzANBgNVBAcMBlN5ZG5leTESMBAGA1UECgwJRGFydW1hdGljMQswCQYD
      VQQLDAJJVDEVMBMGA1UEAwwMYnJhbmRvbi50ZXN0MIICIjANBgkqhkiG9w0BAQEF
      AAOCAg8AMIICCgKCAgEAxl7AtuZa/kXDQKNsgIYbHCvDhOUXW7Jvz8WVMAL94/Fe
      lcktvieClHzIkBYk599G3INpsEBEiersKGyPjlIBPqmrfDJmlZSnpZwnWFhrBbIs
      /EouQe4t6LsqUg+Jj9WpTPSFGAzxqn8OZrMUoMOLj8xRxp8p85ziV9t6CZtfwET6
      laj+Cv7MznsNn8R+cgK2YW+516W7YQgg1szoucBlldoKRR4Xya7h4VcfNa3s4uKx
      RPoBUmLnV5Edes/BCUjwFtC7lenzNjc+mO9El75XGPJxZY+NtTonQ0v5L4rgzUsW
      3Nz2nR36NwOwR+buq/tfodRwR29ZqlJ4mHBDrJntmWmfqnR+WAu6Dsbwt0YTFj1i
      XDRrbXSTHz5Efu+2IQv6wiUdczHD958MZBBzNTpCr7Ss+4gvSBgiVM2yvwiZQJeg
      2I3147d+hz/57J0mLc07tbDQ0pGnTWyAHXvEm7KlO0yZIaTLY6SReTVQsqq16uJT
      flhZWEz13fn3axQpD/OTSXJIRbRyusVJrJKglbFpuUeGjbR/I3K32sZZti4fBpDZ
      ldSnCpuR/z27iBoTpHHg5Aa6SBRO5TjI91yUBdNxw2NHtVzSoY3Z0fGMNfG8hSrK
      DYD2tXJx4k1upiT68HMMYA2kdGVblisBPygL3VM9eZXtoQmcy9BjixhaTWOCyfMC
      AwEAATANBgkqhkiG9w0BAQsFAAOCAgEAtxnHufjkm8ZsQ3aftdsrm8sHL5XUzlYf
      RyH60QLK0Gjl19FkG/sS/XminoZZO0PFFb/Z78L+KVezMj6bd6FNc3ULiKmssQA0
      9Pvzr4c6dyXRapMRWArGCrfYns8vPy8TAJ9DDGV+VNHI2L0VTPk9h/a1qp8qAXmp
      XM0tVCZQrVFc7e6DeCfYYZ/ukAj2n70jUm5iuDTkM5OgbE9XQrbgJeJnGEhm5XrY
      mmE9+G+VXxoAkV2EaNVAzHTg3AeywrLlWArKPL4vl/pG15u5xDDRG074Tkb8gyRS
      UZ+Nc9lNDs2Rw0dsP5E7njnNYkQU81XxgsIu96XSbws0Z5GTGpeHk+CNVycQ9wOV
      Kdyc7aosKxzGUQi69gTa4xmn+EGsboUOappo4fTkP3TeetUYk/79q7AZxxgOTHkr
      fLgrNcrjiUiW91U4ma6PtHbnlzNCl7MYZyy+sxLogR3NFnO8xOtK/1Xdkrtf9/YI
      NmzFsSaumCtxbVsYrTvZMt7eVkJUKL3Kx4K1Vs51emEwtsPB/HJh5ozY2fk3rjlw
      GQ/TM3lH3dViUkhh5DJiGnYU05lOP7aZKR2yWxlqMkdMpnrUq6tF392s1585YdSO
      ohRM5gMVGrL95F/vXln3e4vX4mA4bgr9LLj5xkELJR7UEIXHV6nCeCSutohzejzs
      CwdaB0xYUqI=
      -----END CERTIFICATE-----
    termination: reencrypt
  to:
    kind: Service
    name: nginx
  ```


Route-based Deployment Strategies
----------------------------------

Except for RollingUpdate and Recreate, you can have a more flexible strategy to deploy your application based on the OpenSift Route setting - `Blue-Green` Strategies and `A/B` Strategies.

Let's start by deploying an initial deployment and service

version 1

```
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-1
  labels:
    app: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx-1
  template:
    metadata:
      labels:
        app: nginx-1
    spec:
      containers:
      - name: nginx
        image: quay.io/brandon_tsai/testlab:1
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-1
spec:
  selector:
    app: nginx-1
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
```

Apply the basic route and make sure it works

```
apiVersion: v1
kind: Route
metadata:
  name: nginx
spec:
  host: nginx-uat.apps-crc.testing
  to:
    kind: Service
    name: nginx
```

```
$ oc apply -f route.yml
route.route.openshift.io/nginx created

$ oc get route
NAME    HOST/PORT                    PATH   SERVICES   PORT    TERMINATION   WILDCARD
nginx   nginx-uat.apps-crc.testing          nginx-1    <all>                 None

$ curl http://nginx-uat.apps-crc.testing
<html>
<head>
	<title>Test NGINX passed</title>
</head>
<body>
<h1>NGINX is working</h1>
</body>
</html>
```

Then let's deploy another deployment and service with a new index.html.

```
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx2-index
data:
  index.html: |
    <html>
    <head>
            <title>Test NGINX 2 passed</title>
    </head>
    <body>
    <h1>NGINX 2 is working</h1>
    </body>
    </html>
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-2
  labels:
    app: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx-2
  template:
    metadata:
      labels:
        app: nginx-2
    spec:
      containers:
      - name: nginx
        image: quay.io/brandon_tsai/testlab:1
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
        volumeMounts:
        - name: nginx2-index
          mountPath: /opt/app-root/src
      volumes:
        - name: nginx2-index
          configMap:
            name: nginx2-index
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-2
spec:
  selector:
    app: nginx-2
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080

```

### Blue-Green Deployments

The blue-green deployment strategy can easily switch between two versions of running stateless application stacks. As a consequence, it is very simple and fast to perform a rollback if there is issue in the new version of application.

We have two deployments and services (nginx-1 and nginx-2) in the Openshift cluster.
It is very easy to switch the route to nginx-2 service by patching the target service name.

```
$ oc patch route/nginx -p '{"spec":{"to":{"name":"nginx-2"}}}'
route.route.openshift.io/nginx patched

$ curl http://nginx-uat.apps-crc.testing
<html>
<head>
        <title>Test NGINX 2 passed</title>
</head>
<body>
<h1>NGINX 2 is working</h1>
</body>
</html>
```

We can roll back to nginx-1 service very easily by patching the target service name to nginx-1.

```
$ oc patch route/nginx -p '{"spec":{"to":{"name":"nginx-1"}}}'
route.route.openshift.io/nginx patched

$ curl http://nginx-uat.apps-crc.testing
<html>
<head>
	<title>Test NGINX passed</title>
</head>
<body>
<h1>NGINX is working</h1>
</body>
</html>

```

A/B Deployments
------------

A/B deployments strategy can let a part of clients connect to the new version of deployment to test the new application features and roll back to your initial application (Version 1) or proceed with your new application ( Version 2) fully after the test.

For example:

```
$ oc annotate route/nginx haproxy.router.openshift.io/balance=roundrobin
$ oc set route-backends nginx nginx-1=50 nginx-2=50
```


```
$ for i in {1..10}; do curl -s http://nginx-uat.apps-crc.testing | grep "<h1>" ; done
<h1>NGINX is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX is working</h1>
<h1>NGINX 2 is working</h1>
```

And you can adjust the access ratio between nginx-1 and nginx-2 service by simply one command.

```
$ oc set route-backends nginx nginx-1=20 nginx-2=80
route.route.openshift.io/nginx backends updated

$ for i in {1..10}; do curl -s http://nginx-uat.apps-crc.testing | grep "<h1>" ; done
<h1>NGINX 2 is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX 2 is working</h1>
<h1>NGINX is working</h1>
```


Conclusion
-----------
Compare to Kubernetes, it is super easy to achieve the "Blue/Green" and "A/B" Test in the Openshift Platform. Which give OpenShift more suitable for enterprise requirements.
