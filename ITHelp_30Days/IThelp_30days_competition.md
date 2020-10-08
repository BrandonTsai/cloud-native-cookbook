From Docker Container to RedHat OpenShift
=========================================



16 Docker Concept and installation
17 Develop with Docker
18 Podman: Concept
19 Podman: Run a Pod
20 Quay: Introduction and installation on single VM for Test

21 Quay: Set up Mirror Repository from Dockerhub and check the security scan result
22 K8S Introduction and minikube installation
23 Krew
24 Openshift introduction & install OCP 4.x on MAC
  - set up oc completion
25 OCP:
  - First Pod
  - Customize image and push to Quay.
  - Use images from Quay via Robot Account.
26 OCP: S2I
27 OCP: DeploymentConfig vs Deployment


28 OCP: Config Map & Secrets (Mention Conjur & Vault in Conclusion)
29 OCP: Storage Class, PV and PVC (How to share PV with others)
30 OCP: Service & Route
~

1. OCP: route-based deployment strategies
2. OCP: Network Policy
3. OCP: resource management: Multi-Project Quota, Quota & LimitRange
    - https://learnk8s.io/setting-cpu-memory-limits-requests
4. new-app, template

5. helm overview
6. Operator Introduction
7. OpenShift: container-security-operator with Quay
8. OCP: Buildin Prometheus Operator Introduction
9. OCP: How to monitor external node via Prometheus on OpenShift?
10. OCP: Buildin Grafana & how to set up dashboard.
11. Loki Trial: https://grafana.com/docs/loki/latest/installation/

12. Splunk-connect: introduction (version 1.4.3) and installation by Helm.
13. Falco? kube-bench?
14. Backup/Recovery
15. Conclusion of 30 days challenge ( What have coverd and what did not covered )




Several third-party tools support Helm chart creation such as Draft(https://draft.sh/). Local Helm development is also supported by garden.io and/or skaffold. Check your favorite tool for native Helm support.




Keep Helm Chart on chartmuseum/harbor/Quay/github/gitlab/bitbucket?
-----------------------------

Most charts on https://hub.helm.sh/ are not suitable for OCP because the permission issue.
Beside, you might want to do some customize of the helm chart, such as replace "Ingress" with "Route", Adding extra environment variable.. etc.

In these cases, we have to fork the helm chart and release it to our local private helm repo

https://www.goodwith.tech/blog/hosting-helm-chart-private-repository-in-github-and-gitlab
chartmuseum based
- harbor


OCI-compatible registry
- harbor
- Quay