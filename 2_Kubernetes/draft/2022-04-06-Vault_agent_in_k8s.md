---
title: "#6 HashiCorp Vault: Acquire Secrets in Kubernetes Pods"
author: Brandon Tsai
---


Create kv credentials
-----------------


Create Kubernete Auth
-------------------


Inject Vault to Pod in Same Cluster
------------------------------------


Inject Vault to Pod in Another Cluster
--------------------------------------
vault.hashicorp.com/agent-inject-template-config: |
          {{ with secret "secret/data/web" -}}
            export DB_CONNECTION=postgresql://{{ .Data.data.username }}:{{ .Data.data.password }}@postgres:5432/wizard
          {{- end }}



Vault Template
--------------


Getting Vault Credentials from Environment Varuable?

Update KV Credentials and renew the credentials in Pod
----------------------------------------------------------










### Injecting Secrets into Kubernetes Pods via Vault Agent Containers

https://learn.hashicorp.com/tutorials/vault/kubernetes-sidecar?in=vault/kubernetes







Kubernetes Auth Method
----------------------

https://www.vaultproject.io/docs/auth/kubernetes



Vault Agent with Kubernetes
-------------------------------

https://learn.hashicorp.com/tutorials/vault/agent-kubernetes?in=vault/kubernetes


Vault Agent



kubectl exec \
    $(kubectl get pod -l app=caddy -o jsonpath="{.items[0].metadata.name}") \
    --container caddy -- ls /vault/secrets