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




# curl -k -XGET https://quay.io/cnr/api/v1/packages/community-operators/grafana-operator/blobs/sha256/cc44387393bbb233201a5a02de7697b38cfaa5bb89fbdf22b6b0cd78be3e96ef -o grafana-operator.tar.gz
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  8630  100  8630    0     0  13782      0 --:--:-- --:--:-- --:--:-- 13808

# tar -xf grafana-operator.tar.gz
# mkdir -p manifests/grafana-operator
# cp -r grafana-operator-njs08tsk/* manifests/grafana-operator/

# sed -i "s;quay.io/integreatly/grafana-operator:v3.5.0;quay-uat/mirrors/grafana-operator:v3.5.0;g" ./3.5.0/grafana-operator.v3.5.0.clusterserviceversion.yaml
# sed -i "s;quay.io/integreatly/grafana-operator:v2.0.0;quay-uat/mirrors/grafana-operator:v2.0.0;g" ./2.0.0/grafana-operator.v2.0.0.clusterserviceversion.yaml


Dockerfile
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

podman build -t quay-uat/applications/operator-catalog-registry:0.1.0 .
podman push quay-uat/applications/operator-catalog-registry:0.1.0

```
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

