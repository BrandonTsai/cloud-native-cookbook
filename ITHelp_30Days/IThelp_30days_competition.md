From Docker Container to RedHat OpenShift
=========================================



V 16 Docker Concept and installation
V 17 Develop with Docker
V 18 Podman: Concept
V 19 Podman: Run a Pod
V 20 Quay: Introduction and installation on single VM for Test

V 21 Quay: Set up Mirror Repository from Dockerhub and check the security scan result
V 22 K8S Introduction and minikube installation, kubectl introduction.
V 23 Krew
24 Openshift introduction and "oc cluster up" to run OpenShift 3.x locally.
25 Install Openshift 4.5 (on AWS)
26 OCP: Deployment - How to use images from Quay via Robot Account.
27 OCP: Privillage Permission and Run as non-root

28 OCP: Config Map
29 OCP: Secrets (Mention Conjur & Vault in Conclution)
30 OCP: Storage Class, PV and PVC (How to share PV with others)
1 OCP: DaemonSet & Selector
2 OCP: Service & Route (& Network policy)
3 OCP: resource request, limits, and HPA
4 OCP: resource management: Multi-Project Quota, Quota & LimitRange
    - https://learnk8s.io/setting-cpu-memory-limits-requests

5 Operator Introduction
7 OpenShift: container-security-operator with Quay
8 OCP: Buildin Prometheus Operator: Setup and Test
9 OCP: Buildin & Set up another Grafana Operator: Setup and Test
10 OCP: Buildin Alert Manager: Setup and Test
11 Helm 3: Introduction
12 Keep and use Helm Charts on Quay
13 Splunk-connect: introduction (version 1.4.3) and installation by Helm
15 Backup & Restore
16 Conclution of 30 days challenge ( What have coverd and what did not covered )



(21 Image Signing with Quay)




ReplicaSet
-----------

A ReplicaSet's purpose is to maintain a stable set of replica Pods running at any given time. As such, it is often used to guarantee the availability of a specified number of identical Pods.


A ReplicaSet ensures that a specified number of pod replicas are running at any given time. However, a Deployment is a higher-level concept that manages ReplicaSets and provides declarative updates to Pods along with a lot of other useful features. Therefore, we recommend using Deployments instead of directly using ReplicaSets, unless you require custom update orchestration or don't require updates at all.