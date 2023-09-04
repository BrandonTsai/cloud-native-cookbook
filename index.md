---
layout: home
title: "ToDO"
author: Brandon Tsai
---

https://brandontsai.github.io/cloud-native-cookbook/


Brain Dump - Future Plan
================================================

Container
-------

- Docker Concept
- Podman
  - Introduction and Basic Usage
  - Run containers in a pod
  - (Review) https://www.redhat.com/sysadmin/basic-security-principles-containers
- Buildah
- Skopeo
- Registry
  - Harbor
  - Redhat Quay


k8s & ocp
--------

### Primer: Setup and Management

- Kubernetes Introduction
- Set up
  - minikube
  - kops
  - Rancher labs
  - https://www.kubermatic.com/products/kubermatic/
  - k9s
  - ks3
  - kind?
- RBAC management
- etcd backuo and rollback
- Ingress Controller setup
- OCP Setup

### Basic Deployment

- Krew: cli tools for kubernetes
- From Pod to StatefulSets/Deployment/DaemonSet
- ConfigMap
- Secrets
- PVC, PV
- Health Checks
- Quota and LimitRange
- Autoscaling: HPA & CA
- Service, Ingress, and Deployment Strategies
- Assigning Pods to Nodes

### Advanced Usage

- Helm 3
- Operator
  - Postgresql
  - https://www.noobaa.io/
- Write your own Operator
- BACKUP AND RESTORE EKS USING VELERO
- Knative 
- odo 2.0: https://developers.redhat.com/blog/2020/10/06/kubernetes-integration-and-more-in-odo-2-0/?sc_cid=7013a00000264DlAAI

Observability
----------------------------------------------------------------

- Prometheus and Grafana
- OpenTelemetry
- Splunk
- Chaos Engineering

DevSecOps 
----------------------------------------------------------------

### CICD Pipeline

- Tools:
  * Tekton
  * Argo CD
- Spinnaker
- [Talisman](https://github.com/thoughtworks/talisman)
- SonarQube
- Trivy
- OPA
- KubeSec
- Seal Secrets

### K8S Security

- Istio/Service Mesh
- Kube-bench
- KubeScan
- Kyverno
- sealed-secrets or Vault
- Falco(https://ithelp.ithome.com.tw/articles/10248703)
- Open Policy Agent
- Service Mesh - istio: listen to the traffic, end to end encryption
- CIS.check for k8s
- Network policy: Calico
- HashiCorp Vault Test with Openshift Secrets
- Spiffe: Universal identity control plane for distributed systems