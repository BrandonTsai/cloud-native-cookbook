網路安全，人人有責： Network Policy
========

By default, all Pods are accessible from other Pods and network endpoints. To isolate one or more Pods in a project, you can create `NetworkPolicy` objects in that project to indicate the allowed incoming connections.

If a Pod is matched by selectors in one or more NetworkPolicy objects, then the Pod will accept only connections that are allowed by at least one of those NetworkPolicy objects. A Pod that is not selected by any NetworkPolicy objects is fully accessible.


Only accept connections from Pods within a project
-----------------------

If a pod in UAT project know the clusterIP of another pod in PROD project, actually it can connect to the pod via the clusterIP

```
## Go into UAT project
$ oc project uat
Now using project "uat" on server "https://api.crc.testing:6443".

## get the CLUSTER-IP of pods in PROD project
$ oc get svc -n prod
NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
nginx-prod   ClusterIP   172.25.45.167   <none>        80/TCP    6m23s

## query the PROD service inside a UAT pod
$ oc rsh nginx-uat-5fb55b7d79-hgzzl curl http://172.25.45.167
<html>
<head>
        <title>Test NGINX Prod passed</title>
</head>
<body>
<h1>NGINX Prod is working</h1>
</body>
</html>

```

To avoid this situation, you can apply following Network Policy on PROD project



```
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: allow-same-namespace
spec:
  podSelector:
  ingress:
  - from:
    - podSelector: {}
```

```
$ oc apply -f np-allow-same-project-only.yml -n prod
networkpolicy.networking.k8s.io/allow-same-namespace created

$ oc rsh nginx-uat-5fb55b7d79-hgzzl curl http://172.25.45.167
curl: (7) Failed connect to 172.25.45.167:80; Connection timed out
command terminated with exit code 7
```

Only accept connections from some Pods
----------------

For example, for some backend pods, we only want to allow the connection from nginx pods.

In this case, you can apply following network policy

```
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: allow-pod-and-namespace-both
spec:
  podSelector:
    matchLabels:
      app: api-gw-uat
  ingress:
    - from:
      - podSelector:
          matchLabels:
            app: nginx-uat
```


Before applying the network policy, every pods can access `api-gw-uat` service. 

```
$ oc rsh nginx-uat-5fb55b7d79-hgzzl curl -s -o /dev/null -w "%{http_code}\n"  http://api-gw-uat
200

$ oc rsh test-pod curl -s -o /dev/null -w "%{http_code}\n"  http://api-gw-uat
200

```

After applying the network policy, only nginx pods can access api-gateway pods

```
$ oc apply -f np-allow-pod.yml
networkpolicy.networking.k8s.io/allow-pod-to-api-gw created

$ oc rsh nginx-uat-5fb55b7d79-hgzzl curl -s -o /dev/null -w "%{http_code}\n"  http://api-gw-uat
200

$ oc rsh test-pod curl -s -o /dev/null -w "%{http_code}\n"  http://api-gw-uat
000
command terminated with exit code 7
```

Only allow HTTP and/or HTTPS traffic based on Pod labels
------------------


We can allow https connection only

```
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: allow-https-only-to-nginx
spec:
  podSelector:
    matchLabels:
      app: nginx-uat
  ingress:
  - ports:
    - protocol: TCP
      port: 8443
```

After apply this policy

```
$ oc apply -f np-allow-https-only.yml 
networkpolicy.networking.k8s.io/allow-https-only-to-nginx created

$ oc rsh test-pod curl -s -o /dev/null -w "%{http_code}\n" -k https://nginx-uat
200

## Other pods must use https port to connect to nginx pod
$ oc rsh test-pod curl -s -o /dev/null -w "%{http_code}\n" http://nginx-uat
000 
command terminated with exit code 7

```




Conclusion
------

In Kubernetes, you can use install [Calico](https://docs.projectcalico.org/getting-started/kubernetes/) to do the same thing.



Reference
----------

- https://docs.openshift.com/container-platform/4.5/networking/network_policy/about-network-policy.html