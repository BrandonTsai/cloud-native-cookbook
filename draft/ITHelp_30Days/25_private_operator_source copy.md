When OpenShift Container Platform is installed on restricted networks, also known as a disconnected cluster, Operator Lifecycle Manager (OLM) can no longer use the default OperatorHub sources because they require full Internet connectivity.


Cluster administrators can disable those default sources and create local mirrors so that OLM can install and manage Operators from the local sources instead.


1. Mirror target operator images to local Quay.
2. Build custom Operator Catalog image and push to local Quay.
3. Create new CatalogSource object


I will use Grafana Operator as the example


Steps
-----

1. Disable the default OperatorSources.

```bash
oc patch OperatorHub cluster --type json \
    -p '[{"op": "add", "path": "/spec/disableAllDefaultSources", "value": true}]'
```

2. Retrieve package lists.

```bash
$ curl https://quay.io/cnr/api/v1/packages?namespace=community-operators | jq . > packages.txt
```

3. find the package name and version from package lists


```json
{
    "channels": null,
    "created_at": "2019-07-26T13:43:46",
    "default": "3.5.0",
    "manifests": [
      "helm"
    ],
    "name": "community-operators/grafana-operator",
    "namespace": "community-operators",
    "releases": [
      "3.5.0",
      "2.0.0",
      "1.3.0"
    ],
    "updated_at": "2020-07-31T18:43:47",
    "visibility": "public"
  },
```

4. Pull Operator content

```bash
# curl -k https://quay.io/cnr/api/v1/packages/community-operators/grafana-operator/3.5.0 | jq .  
[
  {
    "content": {
      "digest": "cc44387393bbb233201a5a02de7697b38cfaa5bb89fbdf22b6b0cd78be3e96ef",
      "mediaType": "application/vnd.cnr.package.helm.v0.tar+gzip",
      "size": 8630,
      "urls": []
    },
    "created_at": "2020-07-31T18:43:47",
    "digest": "sha256:0aaed6bdaa093c3eb58378d45ae72f2d332f4632681e59507b8639c29a371b4c",
    "mediaType": "application/vnd.cnr.package-manifest.helm.v0.json",
    "metadata": null,
    "package": "community-operators/grafana-operator",
    "release": "3.5.0"
  }
]
```


5. prepare "manifests" folder

```
$ curl -k -XGET https://quay.io/cnr/api/v1/packages/community-operators/grafana-operator/blobs/sha256/cc44387393bbb233201a5a02de7697b38cfaa5bb89fbdf22b6b0cd78be3e96ef -o grafana-operator.tar.gz
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  8630  100  8630    0     0  13782      0 --:--:-- --:--:-- --:--:-- 13808

$ tar -xf grafana-operator.tar.gz
$ mkdir -p manifests/grafana-operator
$ cp -r grafana-operator-njs08tsk/* manifests/grafana-operator/
```

6. Set up mirror repository in Quay for the operator images in manifests

7. Replace the image url in manifests to your local quay mirror repository

```
$ sed -i "s;quay.io/integreatly/grafana-operator:v3.5.0;quay-uat/mirrors/grafana-operator:v3.5.0;g" ./3.5.0/grafana-operator.v3.5.0.clusterserviceversion.yaml
$ sed -i "s;quay.io/integreatly/grafana-operator:v2.0.0;quay-uat/mirrors/grafana-operator:v2.0.0;g" ./2.0.0/grafana-operator.v2.0.0.clusterserviceversion.yaml
```

8. Create an Operator catalog image.


Add following Dockerfile

```
FROM registry.redhat.io/openshift4/ose-operator-registry:v4.2.24 AS builder

COPY manifests manifests

RUN /bin/initializer -o ./bundles.db

FROM registry.access.redhat.com/ubi7/ubi

COPY --from=builder /registry/bundles.db /bundles.db
COPY --from=builder /usr/bin/registry-server /registry-server
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe

EXPOSE 50051

ENTRYPOINT ["/registry-server"]

CMD ["--database", "bundles.db"]
```

build images and push to local quay registry via podman command

```bash
$ podman build -t quay-uat/applications/operator-catalog-registry:0.1.0 .
$ podman push quay-uat/applications/operator-catalog-registry:0.1.0
```

9. Create a CatalogSource pointing to the new Operator catalog image


```yaml
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: brandon-operator-catalog
  namespace: openshift-marketplace
spec:
  displayName: Brandon Operator Catalog
  sourceType: grpc
  image: quay-uat/applications/operator-catalog-registry:0.1.0
```

apply above yaml file to OpenShift

10. Check OperatorHub from Web UI

![]()


11. Install Grafana Operator from Web UI

![]()



12 Update Grafana Operator for point image to local repository

```
                containers:
                  - args:
                      - >-
                        --grafana-image=quay-uk.windmill.local/gts-base-images/grafana
                      - >-
                        --grafana-plugins-init-container-image=quay-uk.windmill.local/gts-base-images/grafana_plugins_init
                    command:
                      - grafana-operator
```



Conclusion
-------

Above steps are suggested in the OpenShift 4.2 Document, there is another approach to do this for OpenShift 4.5, please refer [here](https://docs.openshift.com/container-platform/4.5/operators/admin/olm-restricted-networks.html).

Happy National Day!
