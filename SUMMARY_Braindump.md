Brain Dump - Future Plan
========================

Container
----------

- Docker Concept
- Podman
- Buildah
- Skopeo
- Registry
  - Harbor
  - Redhat Quay
- Other tips
  - systemd issue with centos/redhat base image

Introduction
------------

- Kubernetes Introduction
- Set up
  - minikube/minishift
  - eksctl
  - Rancher labs
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
- Quota and LimitRange
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
- Other Security settings
- KubeFed && Razee
- AWS X-Ray

Monitoring
----------

- Prometheus and Grafana
  - Scale with M3
- New Relic
- Splunk-connect-for-k8s
CI/CD
-----

- Jenkins
  - "In-process Script Approval" in Jenkins & how to automatic approve.
- Spinnaker
