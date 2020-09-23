https://aws.amazon.com/blogs/opensource/managing-eks-clusters-rancher/

Advantage:
- Friendly UI to set up EKS with default public subnet, eaiser to manage by Non-Tech Operator
- Default Alerts
- tools to forward logs to splunk or elasticsearch
- tools for set up istio

Weakness:
- Manually set up cluster, can not configure by yaml or integrate with ansible.
- Can not scaling nodes from Rancher lab
- Require another server for UI? --> Run in fargate?