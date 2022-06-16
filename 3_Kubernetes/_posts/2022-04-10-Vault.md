Unseal Key 1: v/L0O+y9j8ij42OMJmJDZ1VvuIl27/Jwf/IZY764dAVh
Unseal Key 2: dXzXTa3qjZB7+nnj4TE+FhUY8aIx8UhPp+1I16aUF6Y3
#Unseal Key 3: 50mL9VKT3vbnRWEulvaOJqLM/sL5E6M0PE4e/gr9XHQQ
Unseal Key 4: KWPUG2FWG35TPDXqtNiiQe/v3EwS7ME4rqV5AhEQQKOd
Unseal Key 5: eab5WVORN/17cUP4SjxBAipn6SNgc8gKMwmFOK81qhmn

Initial Root Token: hvs.TAWW3IcWUlDHKq3IeZX33lcN



Vault Server Setup
=========

https://learn.hashicorp.com/tutorials/vault/getting-started-deploy?in=vault/getting-started


WebUI

https://learn.hashicorp.com/tutorials/vault/getting-started-ui?in=vault/getting-started





Secret Engine
============

kv
---

ssh
---


Policies
========

```
 vault policy list
default
root
```

Authentication Method
=====================


token
-----


userpass
-------

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

k8s
---

