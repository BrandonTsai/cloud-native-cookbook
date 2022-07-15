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



### Secret Engine

Secrets engines is used to store, generate or encrypt secrets/credentials.

Vault support different kind of Secrets Engines, such as `key/value`,  `Database` and `SSH key`..etc



### Authentication Method

Authentication in Vault is the process by which user or machine supplied information is verified against an internal or external system. Vault supports multiple auth methods including GitHub, LDAP, AppRole, and more. Each auth method has a specific use case.

Before a client can interact with Vault, it must authenticate against an auth method. Upon authentication, a token is generated. This token is conceptually similar to a session ID on a website. The token may have attached policy, which is mapped at authentication time. This process is described in detail in the policies concepts documentation.

Token



userpass


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




### Policies


Everything in Vault is path-based, and policies are no exception. Policies provide a declarative way to grant or forbid access to certain paths and operations in Vault. This section discusses policy workflows and syntaxes.

Policies are deny by default, so an empty policy grants no permission in the system.

![](https://www.vaultproject.io/_next/image?url=https%3A%2F%2Fcontent.hashicorp.com%2Fapi%2Fassets%3Fproduct%3Dvault%26version%3Drefs%252Fheads%252Frelease%252F1.11.x%26asset%3Dwebsite%252Fpublic%252Fimg%252Fvault-policy-workflow.svg%26width%3D669%26height%3D497&w=1920&q=75)

```
 vault policy list
default
root
```