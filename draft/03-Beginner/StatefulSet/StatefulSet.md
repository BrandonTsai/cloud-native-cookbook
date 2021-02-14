Deployment and StatefulSet compare
=================================

| Deployment | StatefulSet |
| ---------- | ----------- |
| Stateless Pods  |  Stateful Pods  |
| **PersistentVolumeClaim**: <BR> - shared by all pod replicas <BR> - ReadWriteMany / ReadOnlyMany | **volumeClaimTemplates**: <BR> - no shared volume, each replica pod gets a unique PersistentVolumeClaim associated with it <BR> - ReadWriteOnce |


