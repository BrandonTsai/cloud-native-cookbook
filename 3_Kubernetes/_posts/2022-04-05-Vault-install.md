---
title: "#5 HashiCorp Vault: Instroduction"
author: Brandon Tsai
---

What is HashiCorp Vault
---------------




Vault Installation via Helm
-------------------

Follow this [document](https://learn.hashicorp.com/tutorials/vault/kubernetes-minikube-consul) to install Vault on local Kubernetes cluster via Helm


The Vault server can be reached via the CLI and the web UI outside of the Kubernetes cluster if the Vault service running on port 8200 is forwarded.

In another terminal, port forward all requests made to http://localhost:8200 to the vault-0 pod on port 8200.

```
$  kubectl port-forward vault-0 8200:8200
```


Vault Components
--------------

![](https://mktg-content-api-hashicorp.vercel.app/api/assets?product=tutorials&version=main&asset=public%2Fimg%2Fvault%2Fvault-triangle.png)



Secret Engine
------------

Secrets engines are Vault components which store, generate or encrypt secrets.


### kv


### ssh



Policies
---------

```
 vault policy list
default
root
```

Authentication Method
-------------------


### token



### userpass


```bash
vault auth enable userpass
```

1. create 

```
$ vault write auth/userpass/users/brandon \
    password=test123 \
    policies=root
```

2. list users

```
$ vault list auth/userpass/users/ 
Keys
----
brandon
```

3. Get user details

```
$ vault read auth/userpass/users/brandon
Key                        Value
---                        -----
policies                   [default]
token_bound_cidrs          []
token_explicit_max_ttl     0s
token_max_ttl              0s
token_no_default_policy    false
token_num_uses             0
token_period               0s
token_policies             [default]
token_ttl                  0s
token_type                 default
```


4. delete

```
$ vault delete auth/userpass/users/brandon
```



