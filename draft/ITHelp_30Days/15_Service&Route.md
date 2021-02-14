
免除 Ingress Controller 煩惱，擁抱 Route 新世界
==============================

Networking Concept: Service & Route
================


In the traditional VM based environment, if you want to expose a new service, you need to assign a IP to your new VM and add it the the DNS server, set up Nginx to enable TLS certificate, and update the HAProxy configuration and reload HAProxy service to handle the load balance. It is OK if you are in a small startup company, you does not need to settle the networking frequently. However, if you warking in a global organization, which many new services are lauched every day, you might soonly found that it is quick annoying.

The networking model in Kubernetes can help you to reduce this pain by addressing following concerns:

- Containers within a Pod use networking to communicate via loopback (127.0.0.1).
- The `Service` Resource provides communication ability between different Pods.
- The `Ingress` Resource lets you expose a Service of a Pod to be reachable from outside your cluster through `Ingress Controller`.


In OpenShift, it is similar to Kubernetes. However,
Unlike Kubernetes, you have to set up your Ingress Controllers manually before using Ingress Object, Openshift provide a default built-in solution for the external traffic call `Route`.

Following is the concept of OpenShift Networking Model:

![](images/04_OCP_INTRO/network_1.png)



Service Concept Review
----------

An abstract way to expose an application running on a set of Pods as a network service.
With Kubernetes you don't need to modify your application to use an unfamiliar service discovery mechanism. Kubernetes gives Pods their own IP addresses and a single DNS name for a set of Pods, and can load-balance across them.



**Example of Service:**

```
apiVersion: v1
kind: Service
metadata:
  name: nginx
spec:
  selector:
    app: nginx
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
```

This specification creates a new Service object named "nginx-service", which targets TCP port 8080 on any Pod with the app=MyApp label.

Kubernetes assigns this Service an IP address (sometimes called the "cluster IP"), which is used by the Service proxies (see Virtual IPs and service proxies below).

The controller for the Service selector continuously scans for Pods that match its selector, and then POSTs any updates to an Endpoint object also named "my-service".

The default protocol for Services is TCP; you can also use any other supported protocol.


**Discovering services via Built-in DNS**

 If the service is deleted and recreated, a new IP address can be assigned to the service, and requires the frontend pods to be recreated in order to pick up the updated values for the service IP environment variable. Additionally, the backend service has to be created before any of the frontend pods to ensure that the service IP is generated properly, and that it can be provided to the frontend pods as an environment variable.

OpenShift Container Platform has a built-in DNS so that the services can be reached by the service DNS as well as the service IP/port. 
In other pods, they can access the service via this build-in DNS service, for example:

```
$ curl http://nginx
<html>
<head>
	<title>Test NGINX passed</title>
</head>
<body>
<h1>NGINX is working</h1>
</body>
</html>
```


For more detail of service, please refer https://kubernetes.io/docs/concepts/services-networking/service/


OpenShift Route
----------------------------------------------------------

Kubernetes Ingress exposes HTTP and HTTPS routes from outside the cluster to services within the cluster. Traffic routing is controlled by rules defined on the Ingress resource.

OpenShift created a concept called `Route` for the same purpose of Kubernetes Ingress, but it has additional capabilities such as splitting traffic between multiple backends, sticky sessions, etc. When a Route object is created on OpenShift, it gets picked up by the built-in HAProxy load balancer in order to expose the requested service and make it externally available with the given configuration.

![](images/04_OCP_INTRO/network_ingress_route.png)


**Route without ssl Example**


```
$ oc expose svc nginx
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

If you does not specify the hostname, OpenShift will generate a host url `<route-name>-<namespace>.<external-address>`
You can also specify custome hostname, but make sure the hostname in a subdomain of the `external-address`, for example:

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

**Route with ssl termination example**

The simplest way to create a ssl route with edge termination is

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


It will use cluster default widecard certificate.

if you want to use customize certificate, you can use command `oc create route edge --service=nginx --hostname=nginx-uat.apps-crc.testing  --key=nginx-uat.key --cert=nginx-uat.crt`


If you want to redirect http connection to https connection automatically, you can add paramente `--insecure-policy=Redirect`

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

your deployment must enable 8443 port and service must enable 443 port.

add nginx configuration for https

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

generate ConfigMap manually

```
$ oc create configmap nginx-config --from-file=nginx/https.conf -n myproject --dry-run=client -o yaml | oc apply -f -
```

generate secret for tls certificate via command

```
$ oc create secret tls nginx-tls --key="nginx-uat.key" --cert="nginx-uat.crt" -n myproject --dry-run=client -o yaml | oc apply -f -
```

Update deployment 

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
        volumeMounts:
        - name: nginx-tls
          mountPath: /opt/app-root/tls
        - name: nginx-config
          mountPath: /opt/app-root/etc/nginx.d/
      volumes:
        - name: nginx-tls
          secret:
            secretName: nginx-tls
        - name: nginx-config
          configMap:
            name: nginx-config
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

or via YAML file

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

Re-encryption is a variation on edge termination where the router terminates TLS with a certificate, then re-encrypts its connection to the endpoint which may have a different certificate. Therefore the full path of the connection is encrypted, even over the internal network. The router uses health checks to determine the authenticity of the host.



Using the same deployment and service above, but this time we change it to reencrypt.
The communication between external client and Openshift via the cluster default certificate.
The connection inside OpenShift is encripted by our custom certificate.

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


Conclusion
----------

I have introduce the basic concept of Service and Route in OpenShift. I also give 3 examples for different types of TLS termination supported by OpenShift Route. next bloger I will introduce some route-based deployment strategies. 


Reference
---------
- https://docs.openshift.com/container-platform/3.11/dev_guide/routes.html
https://docs.openshift.com/container-platform/3.11/architecture/networking/routes.html#route-types
- https://www.openshift.com/blog/kubernetes-ingress-vs-openshift-route