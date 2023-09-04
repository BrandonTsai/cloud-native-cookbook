

Set the expiration of the token of a service account
----------------------------

Yes, you can set the expiration of the token of a service account in OpenShift. Setting an expiration time for service account tokens is a good security practice as it helps to limit the window of opportunity for potential attackers to use the token.

To set the expiration of a service account token, you can create or modify the service account's secret by setting the expirationTimestamp field in the secret's data. The expirationTimestamp field should be set to the time when the token should expire, represented in the format YYYY-MM-DDTHH:MM:SSZ.

Here's an example of how you can set the expiration time of a service account token:

Get the name of the secret associated with the service account:

```bash
oc describe sa <service_account_name> -n <namespace>
```

Get the current token of the service account:

```bash
oc get secret <secret_name> -n <namespace> -o jsonpath='{.data.token}' | base64 -d
```

Modify the secret associated with the service account to set the expirationTimestamp field:

```csharp
oc get secret <secret_name> -n <namespace> -o json | \
jq '.data |= with_entries(select(.key != "token")) | \
.data.expirationTimestamp = "<expiry_time>" | \
{"apiVersion": "v1", "kind": "Secret", "metadata": .metadata, "type": "kubernetes.io/service-account-token", "data": .data}' | \
oc apply -f -
```

In this command, replace `<secret_name>` with the name of the secret associated with the service account, `<namespace>` with the namespace where the service account is located, and `<expiry_time>` with the time when the token should expire, represented in the format YYYY-MM-DDTHH:MM:SSZ.

After setting the expirationTimestamp field, the service account token will be valid until the specified time. Once the token has expired, you will need to rotate the token to generate a new token for the service account.


Rotate the token of a service account manually
---------------------------------------

Yes, you can rotate the token of a service account in OpenShift. Rotating the token is a good security practice and helps to ensure that the service account is secure.

To rotate the token of a service account, you can use the following command:

```bash
oc rotate secret <secret_name> --namespace=<namespace>
```
In this command, replace `<secret_name>` with the name of the secret associated with the service account that you want to rotate the token for, and `<namespace>` with the namespace where the service account is located.

This command will generate a new token and update the secret associated with the service account with the new token. Once you have rotated the token, you will need to update any applications or processes that use the old token with the new token.


Rotatle the token automatically
-------

Yes, it is possible to rotate the token of a service account automatically in OpenShift. OpenShift provides a feature called "Service Account Token Rotation" which allows you to automate the token rotation process.

When the Service Account Token Rotation feature is enabled, OpenShift automatically rotates the tokens for all service accounts in your cluster based on the defined rotation interval. This ensures that your service accounts always have fresh tokens and helps prevent any unauthorized access to your cluster.

To enable the Service Account Token Rotation feature, you need to create a service-account-controller ConfigMap in the openshift-config namespace with the desired rotation interval. The rotation interval is specified in seconds using the service-account-token-cleanup-interval key. For example, to set the rotation interval to 1 hour, you can create the ConfigMap with the following YAML:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: service-account-controller
  namespace: openshift-config
data:
  service-account-token-cleanup-interval: "3600"
```

Once you have created the service-account-controller ConfigMap, OpenShift will automatically rotate the tokens for all service accounts in your cluster based on the defined rotation interval.

Note that token rotation can cause disruption to services that rely on service account tokens, so it's important to test this feature in a non-production environment first before enabling it in a production environment.


-------

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: job-create-only-role
rules:
- apiGroups: ["batch", "extensions"]
  resources: ["jobs"]
  verbs: ["get", "list", "create"]
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: external-sa-rolebinding
  namespace: gts-lab-dev
subjects:
- kind: ServiceAccount
  name: external
  namespace: gts-lab-dev
roleRef:
  kind: Role
  name: job-create-only-role
  apiGroup: rbac.authorization.k8s.io
```