Brain Dump - Future Plan
========================

Container
----------

- [Docker Concept](01-Container/Docker/Docker_concept.md)
- Podman
  - [Introduction and Basic Usage](01-Container/Podman/01_Concept.md)
  - [Run containers in a pod](01-Container/Podman/02_Run-container-in-pods.md)
  - (Review) https://www.redhat.com/sysadmin/basic-security-principles-containers
- Buildah
- Skopeo
- Registry
  - Harbor
  - Redhat Quay
- Other tips
  - [systemd issue with centos/redhat base image](01-Container/systemd-issue.md)

Introduction
------------
- odo 2.0: https://developers.redhat.com/blog/2020/10/06/kubernetes-integration-and-more-in-odo-2-0/?sc_cid=7013a00000264DlAAI
- Kubernetes Introduction
- Set up
  - minikube
  - kops
  - (Review) [Rancher labs](02-Introduction/Review_Rancherlab.md)
  - https://www.kubermatic.com/products/kubermatic/
- [cli tools for kubernetes](../blogs/05_Improve_Kubectl_Command_with_Krew.md)
- k9s
- ks3


Beginner
--------

- Ingress Controller
- ConfigMap
- Secrets
- PVC, PV
- StatefulSets vs Deployment
- Health Checks
- [Quota and LimitRange](03-Beginner/ResourceManage/quota.md)
- Autoscaling: HPA & CA
- RBAC
- Assigning Pods to Nodes
- Dashboard

Intermediate
------------

- helm 3
- Operator
  - Postgresql
  - https://www.noobaa.io/
- Network policy: Calico
- Open Policy Agent
- Service Mesh - istio: listen to the traffic, end to end encryption
- CIS.check for k8s
- Spiffe
- Vault
- KubeFed && Razee
- AWS X-Ray
- BACKUP AND RESTORE EKS USING VELERO
- Knative 
- HashiCorp Vault Test with Openshift Secrets

Cluster Admin
------------
- Managing TLS Certs
- etcd cluster
- kubeaudit


Monitoring
----------
- Prometheus and Grafana
  - Scale with M3
- New Relic
- [Splunk-connect-for-k8s](05-Monitoring/splunk-connect/01-introduction.md)
  - [Logging setting](05-Monitoring/splunk-connect/02-logging-setting.md)
  - [Metrics setting](05-Monitoring/splunk-connect/03-metrics-setting.md)
  - [Objects setting](05-Monitoring/splunk-connect/04-objects-setting.md)
  - [Splunk dashboard example](05-Monitoring/splunk-connect/05-splunk-dashboard.md)
- Falco(https://ithelp.ithome.com.tw/articles/10248703)

CI/CD
-----
- Jenkins
  - "In-process Script Approval" in Jenkins & how to automatically approve.
- Spinnaker

DevSecOps CICD Pipeline
--------------------------------
* Tekton
* Argo CD

- [Talisman](https://github.com/thoughtworks/talisman)
- SonarQube
- Trivy
- OPA
- KubeSec



K8S Security
----------------------------

- Istio/Service Mesh
- Kube-bench
- Falco
- KubeScan
- Kyverno
- sealed-secrets or Vault


Monitoring
--------------------------------

- OpenTelemetry
- Pixie