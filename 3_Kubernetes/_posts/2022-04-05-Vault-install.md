---
title: "#5 HashiCorp Vault: Instroduction"
author: Brandon Tsai
---

What is HashiCorp Vault
---------------


HashiCorp Vault is an open-source tool that helps you manage and store sensitive credentials easily.

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

https://www.vaultproject.io/docs/internals/architecture


### PATH
In Vault, everything is path based. This means that every operation that is performed in Vault is done through a path. The path is used to determine the location of the operation, as well as the permissions that are required to execute the operation.

A `path` specifies the storage location of your secret.

>> Can Secret and Auth use same path?
### Secret Engine

Secrets engines is used to store, generate or encrypt secrets/credentials.

Vault support different kind of Secrets Engines, such as `key/value`,  `Database` and `SSH key`..etc

```
$ vault secrets list
Path          Type         Accessor              Description
----          ----         --------              -----------
cubbyhole/    cubbyhole    cubbyhole_bf5c1325    per-token private secret storage
identity/     identity     identity_0c631d6d     identity store
sys/          system       system_5f907938       system endpoints used for control, policy and debugging
/
```



The following paths are supported by this backend. To view help for
any of the paths below, use the help command with any route matching
the path pattern. Note that depending on the policy of your auth token,
you may or may not be able to access certain paths.

`dynamic secrets`

### Authentication Method

Authentication in Vault is the process by which user or machine supplied information is verified against an internal or external system. Vault supports multiple auth methods including GitHub, LDAP, AppRole, and more. Each auth method has a specific use case.

Before a client can interact with Vault, it must authenticate against an auth method. Upon authentication, a token is generated. This token is conceptually similar to a session ID on a website. The token may have attached policy, which is mapped at authentication time. This process is described in detail in the policies concepts documentation.

```
$ vault auth list
Path      Type     Accessor               Description
----      ----     --------               -----------
token/    token    auth_token_2bae126c    token based credentials
```

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


Practice
---------


> Create a database secret, and allow web application auth via token to read the database secrets

> Allow a developer to authentication with user/password, and he can create/update.delete secrets


> Create another Vault admin user and rotate root token.


Conclusion
----------

When credential management leaved me wanting and weary, it became my catalyst for action.