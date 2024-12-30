---
layout: home
title: "ToDO"
author: Brandon Tsai
---

https://brandontsai.github.io/cloud-native-cookbook/


Brain Dump - 2025 RoadMap Q1
================================================

Observability
----------------------------------------------------------------

- **OpenTelemetry** --> `**OpenTelemetry Certified Associate**`
  - OpenTelemetry Concept and components
- **Splunk**
  - Collect Prometheus metrics with Splunk-Otel
  - Collect etcd metrics from OCP node exporter to Splunk
  - Collect Journald Logs from OCP node exporter to Splunk
  - From Splunk Enterprise to Splunk Cloud

- **eBPF and Cilium**
  - Monitor Network metrics with Cilium?


Brain Dump - Catalog
===================

Container
-------

- Docker Concept
- Podman
  - Introduction and Basic Usage
  - Run containers in a pod
- Buildah
- Skopeo
- Registry
  - Harbor
  - Redhat Quay
  - **[ Reduce 50% Storage capacity ]**
  - **[ Monitor Container Registry with AquaSec ]**
  - **[ Secure your container images ]**
    - (Review) https://www.redhat.com/sysadmin/basic-security-principles-containers


k8s & ocp
--------

### Primer: Cluster Setup and App Deployment

- Set up
  - minikube
  - k9s
  - ks3
- Krew: cli tools for kubernetes
- From Pod to StatefulSets/Deployment/DaemonSet
- ConfigMap
- Secrets
- **[ Autoscaling with HPA, VPA, CA (Cluster Autoscaler) ]**
- **[ Autoscaling with KEDA]**
- Service, Ingress, and Deployment Strategies
- Assigning Pods to Nodes
- RBAC management
  - Role binding with CLusterRole
  - **[ ServiceAccount Token Rotation ]**

### Advanced Usage

- Operator LCM Concept
  - https://www.noobaa.io/
  - **[ OperatorGroup Introduction ]**
- Write your own Operator
  - odo 2.0: https://developers.redhat.com/blog/2020/10/06/kubernetes-integration-and-more-in-odo-2-0/?sc_cid=7013a00000264DlAAI
- Knative
- BACKUP AND RESTORE EKS USING VELERO


IaC - Infra as Code
-----------------------------------

- Helm3 vs Kustomize
- GitOps --> `**GitOps Certified Associate**`
- ArgoCD --> `**Certified Argo Project Associate**`
  - **[How to set up Helm proxy in ArgoCD]** 



Security
----------------------------------------------------------------

### CICD Pipeline

- Spinnaker
- [Talisman](https://github.com/thoughtworks/talisman)
- SonarQube
- Trivy
- OPA
- KubeSec
- Seal Secrets

### K8S Security

- Service Mesh - istio: listen to the traffic, end to end encryption
- Kube-bench
- KubeScan
- Kyverno
- sealed-secrets
- Falco(https://ithelp.ithome.com.tw/articles/10248703)
- Open Policy Agent
- CIS.check for k8s
- Network policy: Calico
- Spiffe: Universal identity control plane for distributed systems
- eBPF & Cilium