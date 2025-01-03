---
title: "#2-2 ConfigMap and Secret Usage"
author: Brandon Tsai
---

Overview
--------
Many applications require configuration using some combination of configuration files and environment variables.
These configuration files and environment variables should be decoupled from docker image content  in order to keep containerized applications portable.
In Kubernetes based platform, we could use `ConfigMap` and `Secret` object to setting configuration data separately from application code.

A ConfigMap is an object used to store non-confidential data in key-value pairs, it does not provide secrecy or encryption.
If the data you want to store are confidential, You can use a Secret rather than a ConfigMap. However, a Secret is justed encoded by Base64 in Kubernetes, it is not encrypted. The better way is using additional (third party) tools to keep your data private.

| features       | ConfigMap | Secret |
| -------------- | --------- | ------ |
| Binary file    | V         | X      |
| Base64 Encoded | X         | V      |
| Encrypted      | X         | X      |



Create ConfigMap From Directories/file
---------------------------------------

Hint:

- file can be binary
- can not create configmap from nested floders/subfolders
- Use "--dry-run -o yaml " in pipeline, for example, `oc create configmap nginx-config --from-file=configs/nginx -n gts-lab-dev --dry-run -o yaml | oc apply -f -`

```bash
# Create a new configmap from file in folder
oc create configmap nginx-config --from-file=configs/nginx/nginx-app.conf -n gts-lab-dev --dry-run -o yaml | oc apply -f -

oc create configmap nginx-html --from-file=configs/html -n gts-lab-dev --dry-run -o yaml | oc apply -f -

oc create configmap nginx-html-icons --from-file=configs/html/icons -n gts-lab-dev --dry-run -o yaml | oc apply -f -

oc create configmap nginx-app1-env --from-env-file=configs/app1.env -n gts-lab-dev --dry-run -o yaml | oc apply -f -

# Create a new configmap named my-ssh-key with specified keys instead of names on disk
oc create configmap my-ssh-config --from-file=ssh-config=~/.ssh/config --from-file=ssh-publickey=~/.ssh/id_rsa.pub --dry-run -o yaml | oc apply -f -

# Create a new configMap named my-config with key1=config1 and key2=config2
oc create configmap proxy-env --from-literal=HTTP_PROXY=https://myproxy.example.local --from-literal=NO_PROXY="*.example.local" --dry-run -o yaml | oc apply -f -

```


Create Secrets From Directories/file
---------------------------------------

The Secret object type provides a mechanism to hold sensitive information such as passwords

Hint:

- can not create secrets from nested floders/subfolders
- Use "--dry-run -o yaml " in pipeline

```bash
# Create a new secret from file
oc create secret generic nginx-ssl-key --from-file=configs/ssl.key -n gts-lab-dev --dry-run -o yaml | oc apply -f -

# Create a new secret from file with specified keys
oc create secret generic nginx-ssl-key --from-file=nginx-key=configs/ssl.key -n gts-lab-dev --dry-run -o yaml | oc apply -f -

# Create a new secret from raw data
oc create secret generic pgsql-secret --from-literal pgsql_user=brandon --from-literal pgsql_key=testing123 --dry-run -o yaml | oc apply -f -
```



Consuming ConfigMap & Secrets in Pods
-----------------------------

ConfigMap and Secret can be used to:

- Populate environment variable values in containers
- Populate configuration files in a volume

Notice that a ConfigMap and Secret must be created before its contents can be consumed in Pods.

Following is the example: 


```yaml
kind: Deployment
apiVersion: apps/v1
metadata:
  name: nginx-app1
  namespace: gts-lab-dev
  labels:
    app: nginx-app1
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx-app1
  template:
    metadata:
      labels:
        app: nginx-app1
    spec:
      containers:
      - name: nginx-app1
        image: registry.redhat.io/rhscl/nginx-114-rhel7:1
        args:
        - nginx
        - -g
        - 'daemon off;'
        ports:
        - containerPort: 8001
          protocol: TCP
        env:
          - name: OS_TYPE
            valueFrom:
              configMapKeyRef:
                name: example-env
                key: os.type
        envFrom:
          - configMapRef:
              name: nginx-app1-env
          - secretRef:
              name: pgsql-secret
        volumeMounts:
        - name: nginx-ssl-key
          mountPath: /opt/app-root/etc/nginx.d/ssl
        - name: nginx-config
          mountPath: /opt/app-root/etc/nginx.d/nginx-app.conf
          subPath: nginx-app.conf
        - name: nginx-html
          mountPath: /opt/app-root/src
        - name: nginx-html-icons
          mountPath: /opt/app-root/src/icons
      volumes:
        - name: nginx-ssl-key
          secret:
            secretName: nginx-ssl-key
        - name: nginx-config
          configMap:
            name: nginx-config
            items:
            - key: nginx-app.conf
              path: nginx-app.conf
        - name: nginx-html
          configMap:
             name: nginx-html
        - name: nginx-html-icons
          configMap:
            name: nginx-html-icons
```

Risks of using secrets in Kubernetes
----------

- In the API server, secret data is stored in etcd; therefore:
  - Administrators should enable encryption at rest for cluster data (requires v1.13 or later).
  - Administrators should limit access to etcd to admin users.
  - Administrators may want to wipe/shred disks used by etcd when no longer in use.
  - If running etcd in a cluster, administrators should make sure to use SSL/TLS for etcd peer-to-peer communication.
- If you configure the secret through a manifest (JSON or YAML) file which has the secret data encoded as base64, sharing this file or checking it in to a source repository means the secret is compromised. Base64 encoding is not an encryption method and is considered the same as plain text.
- Applications still need to protect the value of secret after reading it from the volume, such as not accidentally logging it or transmitting it to an untrusted party.
- A user who can create a Pod that uses a secret can also see the value of that secret. Even if the API server policy does not allow that user to read the Secret, the user could run a Pod which exposes the secret.
- Currently, anyone with root permission on any node can read any secret from the API server, by impersonating the kubelet. It is a planned feature to only send secrets to nodes that actually require them, to restrict the impact of a root exploit on a single node.


How OpenShift handle the secrets issue.
------------------------------------------

From the case description, I understand that you are looking to leverage resources encryption, in particular secret encryption at the datastore layer. By default, etcd data is not encrypted in OpenShift Container Platform. You can enable etcd encryption for your cluster to provide an additional layer of data security.

When you enable etcd encryption, the following OpenShift API server and Kubernetes API server resources are encrypted:

    Secrets
    ConfigMaps
    Routes
    OAuth access tokens
    OAuth authorize tokens

You can get more insights on the same after going through this [documentation](https://docs.openshift.com/container-platform/4.5/security/encrypting-etcd.html#enabling-etcd-encryption_encrypting-etcd). 

https://docs.openshift.com/container-platform/4.5/security/encrypting-etcd.html#enabling-etcd-encryption_encrypting-etcd

K8S
https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/

