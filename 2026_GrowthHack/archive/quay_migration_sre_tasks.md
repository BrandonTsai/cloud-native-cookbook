# 🏗️ Red Hat Quay: VM → Operator on OpenShift Migration
## Senior Platform Engineer / SRE Demonstration Tasks

> **Objective:** Demonstrate Principal-level Platform Engineering and SRE competencies through this migration project to support a salary uplift request.
>
> **Strategy:** Every task is designed to be *evidence you can speak to in an interview or performance review* — not just "I migrated Quay" but "I designed a zero-downtime migration with documented runbooks, SLOs, and rollback procedures."

---

## 📋 How to Use This File

- `[ ]` = Not started | `[x]` = Done | `[~]` = In progress
- Each phase maps to a **Senior PE/SRE competency** you can directly reference in salary discussions
- Add notes after each task with actual commands, decisions, and outcomes — this becomes your **evidence portfolio**

---

## Phase 0: Discovery & Architecture Assessment

> **Competency Demonstrated:** Platform ownership, risk assessment, stakeholder communication
> **Salary Lever:** "I drove the discovery and architecture decisions, not just executed tasks."

- [ ] **P0-1** Inventory current VM-based Quay deployment: version, number of repos, image count, total storage used, active users, integrations (LDAP/OIDC, CI/CD clients)
- [ ] **P0-2** Document current `config.yaml`: identify all non-default settings, LDAP/OIDC config, custom TLS certs, storage backend (S3/Ceph/NFS), feature flags
- [ ] **P0-3** Map all upstream consumers of Quay: CI/CD pipelines, OpenShift clusters pulling images, developer workstations, external systems — store in `consumers-map.md`
- [ ] **P0-4** Identify managed vs unmanaged components decision for Operator deployment:
  - PostgreSQL: managed (Operator) vs external (existing DB)
  - Object storage: NooBaa/ODF managed vs unmanaged (existing S3/Ceph)
  - Redis: managed vs external
  - Clair (vulnerability scanning): managed vs unmanaged
  - TLS: Operator-managed vs custom cert
- [ ] **P0-5** Write Architecture Decision Record (ADR): `adr-001-quay-operator-migration.md`
  - Why Operator over VM? (lifecycle management, Day-2 ops, reconciliation loop)
  - Managed vs unmanaged component decisions with rationale
  - Risk register: data loss, auth downtime, DNS cutover window
- [ ] **P0-6** Define migration approach: **in-place operator migration** vs **parallel new deployment + cutover** — document chosen approach with rationale
- [ ] **P0-7** Agree on maintenance window with stakeholders; document in `migration-plan.md`

---

## Phase 1: Pre-Migration Backup & Validation

> **Competency Demonstrated:** Operational excellence, disaster recovery ownership
> **Salary Lever:** "I owned the backup and rollback strategy — no data loss risk."

- [ ] **P1-1** Back up Quay `config.yaml` from VM:
  ```bash
  mkdir /tmp/quay-backup
  cp /path/to/quay/config/config.yaml /tmp/quay-backup/
  ```
- [ ] **P1-2** Perform PostgreSQL database dump from current Quay DB:
  ```bash
  pg_dump -h DB_HOST -p 5432 -d QUAY_DB -U QUAY_USER -W -O > /tmp/quay-backup/quay-db-backup.sql
  ```
- [ ] **P1-3** Back up all image blobs from object storage:
  ```bash
  aws s3 sync --no-verify-ssl --endpoint-url https://STORAGE_ENDPOINT s3://QUAY_BUCKET/ /tmp/quay-backup/bucket-backup/
  ```
- [ ] **P1-4** Validate backup integrity:
  - Verify SQL dump file size is non-zero and parseable
  - Spot-check 3+ critical image blobs exist in bucket-backup
  - Record total blob count and storage size in `backup-validation.md`
- [ ] **P1-5** Store backup in a secondary location (separate storage, not the same cluster)
- [ ] **P1-6** Document rollback procedure in `rollback-runbook.md`:
  - How to restore VM-based Quay if operator migration fails
  - DNS/load balancer revert steps
  - Estimated RTO (Recovery Time Objective) for rollback

---

## Phase 2: OpenShift Environment Preparation

> **Competency Demonstrated:** OpenShift platform administration, capacity planning, security posture
> **Salary Lever:** "I set up the platform correctly, not just ran a wizard."

- [ ] **P2-1** Validate OpenShift cluster prerequisites:
  - OCP version ≥ 4.5
  - Cluster-admin access confirmed
  - Default StorageClass configured (required for Quay + Clair PostgreSQL PVCs)
  - Minimum node capacity: 8Gi RAM + 2 vCPUs available for Quay pods
- [ ] **P2-2** Create dedicated namespace:
  ```bash
  oc new-project quay-enterprise
  ```
- [ ] **P2-3** (Optional but senior-level) Schedule Quay pods on dedicated infrastructure nodes:
  - Label infra nodes: `oc label node <node> node-role.kubernetes.io/infra=`
  - Taint infra nodes: `oc adm taint nodes -l node-role.kubernetes.io/infra node-role.kubernetes.io/infra=reserved:NoSchedule`
  - Annotate namespace with node-selector and tolerations
  - Document in `infra-node-scheduling.md` — why isolation matters for production workloads
- [ ] **P2-4** Configure object storage for Operator:
  - If using NooBaa/ODF: verify `ObjectBucketClaim` API is available
  - If using external S3/Ceph: prepare `DISTRIBUTED_STORAGE_CONFIG` credentials
  - Document choice and configuration in `storage-config.md`
- [ ] **P2-5** Prepare `configBundleSecret` from backed-up config:
  ```bash
  # Carry over SECRET_KEY, LDAP config, OIDC config, storage config, feature flags
  cat /tmp/quay-backup/config.yaml | grep SECRET_KEY > /tmp/quay-backup/config-bundle.yaml
  # Manually add LDAP, OIDC, custom TLS, and all non-default settings
  oc create secret generic quay-config-bundle     --from-file=config.yaml=/tmp/quay-backup/config-bundle.yaml     -n quay-enterprise
  ```
- [ ] **P2-6** Verify LDAP/OIDC connectivity from OpenShift cluster network (firewall rules, service endpoints reachable)
- [ ] **P2-7** Prepare TLS certificates (custom certs or confirm Operator-managed TLS via OpenShift Routes)

---

## Phase 3: Quay Operator Installation & Initial Deployment

> **Competency Demonstrated:** Operator lifecycle management, GitOps/declarative configuration
> **Salary Lever:** "I deployed Quay as a managed OpenShift-native application, not a VM service."

- [ ] **P3-1** Install Red Hat Quay Operator from OperatorHub:
  - Choose cluster-wide installation (enables monitoring component)
  - Set update channel (e.g. `stable-3.x`)
  - Set approval strategy: **Manual** for production (review upgrades before applying)
  - Document channel selection and approval strategy rationale
- [ ] **P3-2** Create `QuayRegistry` Custom Resource with correct managed/unmanaged component flags:
  ```yaml
  apiVersion: quay.redhat.com/v1
  kind: QuayRegistry
  metadata:
    name: quay-registry
    namespace: quay-enterprise
  spec:
    configBundleSecret: quay-config-bundle
    components:
      - kind: quay
        managed: true
      - kind: postgres
        managed: true        # or false if using external DB
      - kind: clair
        managed: true
      - kind: redis
        managed: true
      - kind: horizontalpodautoscaler
        managed: true
      - kind: objectstorage
        managed: false       # if using external S3/Ceph
      - kind: route
        managed: true
      - kind: mirror
        managed: true
      - kind: monitoring
        managed: true
      - kind: tls
        managed: true
      - kind: clairpostgres
        managed: true
  ```
- [ ] **P3-3** Monitor operator reconciliation:
  ```bash
  oc describe quayregistry quay-registry -n quay-enterprise
  oc get pods -n quay-enterprise -w
  ```
- [ ] **P3-4** Verify all pods reach `Running` status (quay-app, quay-database, clair-app, clair-postgres, quay-redis, quay-mirror)
- [ ] **P3-5** Confirm `status.registryEndpoint` is set on `QuayRegistry` CR
- [ ] **P3-6** Test fresh Quay instance (pre-data migration): login, push/pull a test image

---

## Phase 4: Data Migration (VM → Operator)

> **Competency Demonstrated:** Zero-downtime migration execution, database operations, data integrity verification
> **Salary Lever:** "I executed a live production data migration with no data loss and documented every step."

- [ ] **P4-1** Scale down Quay Operator temporarily to prevent reconciliation during DB restore:
  ```bash
  oc scale --replicas=0 deployment quay-operator.<version> -n openshift-operators
  ```
- [ ] **P4-2** Scale down Quay app and mirror deployments:
  ```bash
  oc scale --replicas=0 deployment <QUAY_APP_DEPLOYMENT> <QUAY_MIRROR_DEPLOYMENT> -n quay-enterprise
  ```
- [ ] **P4-3** Copy database backup into Operator-managed PostgreSQL pod:
  ```bash
  oc cp /tmp/quay-backup/quay-db-backup.sql quay-enterprise/<quay-postgres-pod>:/var/lib/pgsql/data/userdata/
  ```
- [ ] **P4-4** Restore database from backup:
  ```bash
  # Exec into postgres pod
  oc exec -it <quay-postgres-pod> -n quay-enterprise -- /bin/bash
  # Drop and recreate DB, restore from SQL dump
  psql -h localhost -d QUAY_DB -U QUAY_USER -W < /var/lib/pgsql/data/userdata/quay-db-backup.sql
  ```
- [ ] **P4-5** Sync image blobs to new storage backend:
  ```bash
  aws s3 sync --no-verify-ssl     --endpoint-url https://NOOBAA_OR_NEW_STORAGE_ENDPOINT     /tmp/quay-backup/bucket-backup/*     s3://QUAY_DATASTORE_BUCKET_NAME
  ```
- [ ] **P4-6** Patch `QuayRegistry` CR with updated `configBundleSecret`:
  ```bash
  oc patch quayregistry quay-registry --type=merge     -p '{"spec":{"configBundleSecret":"quay-config-bundle"}}'     -n quay-enterprise
  ```
- [ ] **P4-7** Scale Quay Operator back up:
  ```bash
  oc scale --replicas=1 deployment quay-operator.<version> -n openshift-operators
  ```
- [ ] **P4-8** Scale Quay app and mirror pods back up:
  ```bash
  oc scale --replicas=1 deployment quayregistry-quay-app quayregistry-quay-mirror -n quay-enterprise
  ```

---

## Phase 5: Post-Migration Validation

> **Competency Demonstrated:** Quality assurance, SRE acceptance testing, operational handover
> **Salary Lever:** "I defined and executed a formal acceptance test plan, not just 'it seemed to work'."

- [ ] **P5-1** Verify authentication: login with LDAP/OIDC accounts that existed pre-migration
- [ ] **P5-2** Verify all organisations, repositories, and teams are present
- [ ] **P5-3** Pull 5+ critical images that existed pre-migration — verify digests match originals
- [ ] **P5-4** Push a new test image; verify it appears in Quay UI and is pullable
- [ ] **P5-5** Verify Clair vulnerability scanning is functional (trigger a scan, check results)
- [ ] **P5-6** Verify repository mirroring is operational (if used)
- [ ] **P5-7** Test robot accounts and service account tokens used by CI/CD pipelines
- [ ] **P5-8** Verify OpenShift integrated registry pull-secret still works (cluster pulls from Quay)
- [ ] **P5-9** Document all test results in `acceptance-test-results.md` with pass/fail status
- [ ] **P5-10** Perform DNS/load balancer cutover to new Quay endpoint (or update Route hostname)

---

## Phase 6: Observability & SLO Implementation

> **Competency Demonstrated:** SRE observability practice, SLO ownership, proactive alerting
> **Salary Lever:** "I defined SLOs for the registry and built the monitoring stack — this is Principal SRE work."

- [ ] **P6-1** Enable Quay Operator monitoring component (Grafana dashboard + metrics):
  - Confirm `monitoring: managed: true` in `QuayRegistry` CR
  - Access Quay Grafana dashboard in OpenShift monitoring stack
- [ ] **P6-2** Define SLOs for Quay registry — store in `quay-slos.yaml`:
  ```yaml
  # Availability SLO: 99.9% of image pull requests succeed (30-day window)
  # Latency SLO: p95 image pull response time < 2 seconds
  # Push Success SLO: 99.5% of image push requests succeed
  # Clair Scan SLO: 95% of new images scanned within 10 minutes of push
  ```
- [ ] **P6-3** Create Prometheus alerting rules for SLO burn rate (fast-burn + slow-burn alerts)
- [ ] **P6-4** Build OpenShift/Grafana dashboard showing:
  - Quay pod health and restart count
  - Push/pull request rate and error rate
  - Storage utilisation trend
  - Clair scan queue depth
- [ ] **P6-5** Set up alerting for:
  - Quay pod CrashLoopBackOff
  - PostgreSQL storage > 80% capacity
  - Error budget burn rate exceeding threshold
  - Clair database out of sync (vulnerability DB staleness)
- [ ] **P6-6** Export all alerting rules and dashboards as YAML/JSON (Observability as Code)
- [ ] **P6-7** Document alert runbooks in `runbook-quay-alerts.md` (one runbook per alert)

---

## Phase 7: Day-2 Operations & Operator Lifecycle Management

> **Competency Demonstrated:** Platform lifecycle ownership, upgrade strategy, operator pattern mastery
> **Salary Lever:** "I own the Quay platform ongoing — upgrades, scaling, and config changes are GitOps-driven."

- [ ] **P7-1** Document Quay Operator upgrade procedure:
  - How to review release notes before approving upgrade (Manual approval strategy)
  - How to upgrade: approve `InstallPlan` in OCP console or via `oc`
  - How to verify successful upgrade (`oc get csv -n openshift-operators`)
- [ ] **P7-2** Implement HorizontalPodAutoscaler (HPA) for Quay app pods:
  - Confirm `horizontalpodautoscaler: managed: true` in QuayRegistry CR
  - Document min/max replica settings and CPU trigger threshold
- [ ] **P7-3** Implement GitOps for `configBundleSecret` changes:
  - Store `config.yaml` changes in Git (secrets managed via Sealed Secrets or Vault)
  - Any config change goes through PR review → `oc apply` → Operator reconciles
- [ ] **P7-4** Test config change via Operator reconciliation (add a feature flag, verify Quay pods restart cleanly)
- [ ] **P7-5** Implement backup automation:
  - CronJob or pipeline to run `pg_dump` + S3 sync on schedule (daily)
  - Alert if backup job fails
  - Document RTO/RPO targets in `dr-plan.md`
- [ ] **P7-6** Implement Quay garbage collection schedule (remove untagged blobs):
  - Document `CLEAN_BLOB_UPLOAD_FOLDER` and GC worker configuration
  - Monitor storage reclaim after first GC run

---

## Phase 8: Security Hardening

> **Competency Demonstrated:** Security as reliability, supply chain security, least-privilege
> **Salary Lever:** "I hardened the container registry — the most critical piece of supply chain security."

- [ ] **P8-1** Enforce image signing policy: configure Quay to require signed images for production namespaces (cosign integration or Quay robot account restrictions)
- [ ] **P8-2** Verify Clair is scanning all pushed images; confirm CVE database is auto-updated
- [ ] **P8-3** Implement Quay quota enforcement per organisation/repository (storage limits)
- [ ] **P8-4** Review and harden RBAC:
  - Audit robot accounts — remove unused accounts
  - Enforce principle of least privilege for robot account permissions (read-only vs read-write)
  - Document robot account inventory in `robot-accounts-registry.md`
- [ ] **P8-5** Enable and test repository mirroring with signature verification (if applicable)
- [ ] **P8-6** Verify TLS is enforced end-to-end (no HTTP fallback); validate cert expiry monitoring
- [ ] **P8-7** Run OpenShift Security Context Constraints (SCC) audit on Quay pods
- [ ] **P8-8** Integrate Quay vulnerability scan results into CI/CD pipeline (block images with CRITICAL CVEs)
- [ ] **P8-9** Document security baseline in `security-baseline.md`: findings before vs after hardening

---

## Phase 9: Documentation & Knowledge Transfer

> **Competency Demonstrated:** Platform ownership maturity, team enablement, engineering leadership
> **Salary Lever:** "I didn't just migrate it — I made it maintainable for the whole team."

- [ ] **P9-1** Write `architecture-overview.md`: current state architecture diagram (Mermaid or draw.io), component responsibilities, data flows
- [ ] **P9-2** Write `operations-runbook.md`: common Day-2 tasks (add org, reset robot account, trigger GC, approve operator upgrade)
- [ ] **P9-3** Write `incident-response-quay.md`: Quay-specific incident playbook (Quay down, storage full, Clair not scanning, auth failure)
- [ ] **P9-4** Write `onboarding-guide.md`: how a new developer gets access to Quay, creates a robot account, integrates with CI/CD
- [ ] **P9-5** Conduct knowledge transfer session with team; record it if possible
- [ ] **P9-6** Add Quay runbooks to team wiki/Confluence/internal docs platform

---

## 💰 Salary Uplift Evidence Map

Use these talking points directly in your salary review or interview:

| What You Did | Senior PE/SRE Competency | Salary Argument |
|---|---|---|
| Phase 0: Discovery & ADR | Platform architecture ownership | "I led the technical design, not just executed." |
| Phase 1: Backup + Rollback Plan | Operational excellence, DR | "Zero-risk migration with documented RTO/RPO." |
| Phase 2: Infra node scheduling | OpenShift platform admin | "Production-grade scheduling and isolation." |
| Phase 4: Live DB + blob migration | Complex migration execution | "I migrated production data without data loss." |
| Phase 5: Formal acceptance tests | Quality gate ownership | "I defined what 'done' means, not just 'it works'." |
| Phase 6: SLOs + Error Budgets | Principal SRE practice | "I own the reliability contract for this platform." |
| Phase 7: Operator lifecycle GitOps | Platform engineering maturity | "Day-2 ops are automated and auditable." |
| Phase 8: Supply chain security | Security engineering | "I hardened the most critical artifact in our supply chain." |
| Phase 9: Runbooks + KT | Engineering leadership | "I made the platform maintainable for the whole org." |

---

## 📊 Competency Coverage

| Skill Area | Tasks Covered | Interview Evidence |
|---|---|---|
| OpenShift Platform Admin | P2, P3, P7 | QuayRegistry CR, infra node scheduling, HPA |
| Container Registry Operations | P3, P4, P5 | Full migration execution, acceptance testing |
| SRE Practice (SLOs/Alerts) | P6 | `quay-slos.yaml`, burn rate alerts, dashboards |
| Database Operations | P1, P4 | pg_dump, restore, integrity validation |
| Security Engineering | P8 | Clair, RBAC, TLS, supply chain hardening |
| Incident Management | P6, P9 | Runbooks, alert playbooks, on-call docs |
| Documentation & Leadership | P9 | ADRs, architecture docs, knowledge transfer |

---

*Last updated: 2026-02-27 | Project: quay-vm-to-operator-migration*
