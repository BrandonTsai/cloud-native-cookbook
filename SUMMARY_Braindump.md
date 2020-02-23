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

- Kubernetes Introduction
- Set up
  - minikube/minishift
  - eksctl
  - (Review) [Rancher labs](02-Introduction/Review_Rancherlab.md)
  - Openshift 4 on AWS
- cli tools
  - Enabling shell autocompletion
  - kubectx
  - kubens
  - krew: grep plugin as example

Beginner
--------

- Pod and Deployment
- Service
- Ingress Controller
- Configmap
- Secrets
- PVC, PV
- StatefulSets
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
  - Ceph
  - GlusterFS
  - https://www.noobaa.io/
- Network policy: Calico
- Open Policy Agent
- Service Mesh - istio: listen to the traffic, end to end encryption
- CIS.check for k8s
- Spiffe
- Vault
- KubeFed && Razee
- AWS X-Ray

CLuster Admin
------------
- Managing TLS Certs
- etcd cluster


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


CI/CD
-----

- Jenkins
  - "In-process Script Approval" in Jenkins & how to automatically approve.
- Spinnaker
