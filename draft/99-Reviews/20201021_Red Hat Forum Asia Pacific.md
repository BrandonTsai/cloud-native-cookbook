Keynotes:

Prescriptive Self-service Adoption of Native AWS Services to Compliment Application Workloads in OpenShift on AWS 
---------------

why OCP on AWS?
- scalability
- availability
- security
- elasticity

OCP: stable and secure
AWS: wide range of service.


Self Service for Dev on AWS
- broad choices on AWS.
- allow dev deploy on ocp with the integration of these AWS servicebroker
  - github/awslab/aws-servicebroker
- Dev can create new app of AWS service from Catalog
  - s3 for object storage.
  - database
  - lambda for serverless.

What is the benefit to using OCP Service Catalog for using aws service insteading of managing by other tools?
- without leaving to another tool
- limit access
- share portfolios
- constraint of using products

Run ECS/EKS with service broker in OCP?

Product -> Portfolio -> Contraint -> Products list -> Provisioned Products



Hybrid Cluster Management
-------------
Red Hat Insights
Red Hat Satellite
Red Hat Advance Cluster Management


Ansible Automation Platform | Strategy and Vision
-----------

Automation Hub
Automation Analytics
Automation services catalog

-> Use ansible to bring content to OpenShift



How to be Cloud Native Architecture
-------------

Cloud Native Architecture:
- CI/CD DevOps Process
- Portable Application Platform
- Loosely Coupled Architecture



Selecting your application service architecture.
- It is not about all or nothing, it is about the best fit


Microservice -> Decomposition
- Rapid infra provision
- Basic monitoring and management
- Rapid develop automatically

Deploy Strategies:
- Recreate
- Rolling
- blue/green
- canary
- a/b testing
- shadow


How to measure your Cloud Native Architecture?
- the mean time of recovery



OCP: Host and Manage Services
----------------------

ocp support hybrid cloud service.

k8s done right is hard
- install, deploy, harden, operate

building a platform is not your business focus.




Integration - The Missing Link in your Serverless Adoption Strategy
---------------------------



3 Ways to Develop for OpenShift | A Deep Dive into Tools and Strategies 
----------------------------

- Learning Subscription (Very Expensive, for java engineer ) 
- Red Hat Administrator Certificated


### Remote developers/Consultants often takes 5-10 days to set up on the internal system due to VPN or VDI.
-> Red HAt CodeReady Workspaces
  -> Make developers develop container-based app on OCP esay.
  -> easy latency of network connetivity
-> OCP VSCode plug-in
-> Devfile: Developer Environment as Code.


### Local Development:
- CodeReady Containers: OCP on your laptop
  - replace minishift, CDK and "oc cluster up"
- CodeReady Studio (IDE): 
  - Based on Eclipse, support Java, node, spring boot boot
- IDE Plugin
- odo
  - OCP dev-focused cli
  - what is the different with new-app?

Tekton Pipelines

Challenge:
- Admin has less control on dev env


### Local with Remote Target.



DevOps for Grown Ups 
----------

Self Service from public image repo
-> 97% of Vulnerabilities from public image repo

DevOps fone wrong:
- non compliant with security policy
- large cost and shadow IT


Put Governance into DevOps:
- Preventative
- Corrective
- Detective

might slow the process

### Grown Up

left-shifted Governanace:


DevSecOps Mapped Out

Design and Implementation
- Platform
  - Manage platform risk
  - evaluate cost for platform


Workflow management
- what is input, 
- whst action to handle the input?
- what are possible risks of input and action
- Control risk


CI:
- Compliance as Code: Falco, Compliance Operator (ACM), Open Policy Agent
- GitOps (Part 1)
- Revision Control

CD: Faster change Automation
- CAB and eCAB 2.0
- GitOps (part 2)


Audit:
- immutable Record: tools such as Frafeas
- cut down audit from month to days











