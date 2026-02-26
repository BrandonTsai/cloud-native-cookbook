# 🏗️ Red Hat Quay: VM → Operator on OpenShift Migration
## Senior Platform Engineer / SRE Demonstration Tasks

> **Migration Drivers:**
> 1. 💾 **Reduce storage usage** — enable quota management, garbage collection tuning, and blob deduplication
> 2. 💰 **Avoid license fee** — transition to Project Quay (community) or align with OpenShift Platform Plus entitlement
> 3. 🌏 **Architecture Decision Record** — determine whether to deploy 1 Global Quay (geo-replication) or 2 Independent Quay instances (EU + AP) with repository mirroring
>
> **Goal:** Every task is designed as *demonstrable evidence* for a senior salary uplift conversation.

---

## 📋 How to Use This File

- `[ ]` = Not started | `[x]` = Done | `[~]` = In progress
- Add notes/commands/outcomes after each task — this becomes your **evidence portfolio**
- Each phase maps to a named SRE/Platform Engineering competency

---

## ADR-001: Global Quay (Geo-Replication) vs Independent Regional Instances

> **Competency:** Architecture ownership, trade-off analysis, Principal-level technical decision-making
> **This ADR is your most senior deliverable — write it before any implementation begins.**

### Context
The business operates in EU and AP regions. We need a container registry that is available in both regions with acceptable image pull latency.

### Options Evaluated

| Factor | Option A: 1 Global Quay (Geo-Replication) | Option B: 2 Independent Quay Instances + Repository Mirroring |
|---|---|---|
| **Architecture** | Single `QuayRegistry` per cluster, shared PostgreSQL + Redis, region-local object storage | Two fully independent `QuayRegistry` deployments, no shared DB |
| **Storage usage** | Blobs replicated to ALL regions by default — higher total storage | Only mirrored repos replicated — lower storage, selective replication |
| **License / cost** | Single logical deployment — one entitlement | Two separate deployments — two entitlements (if licensed) |
| **Latency** | Pulls from nearest storage; push to preferred region, async replication | Pulls from local region; push only to primary, mirror syncs async |
| **Failover** | Requires Global Load Balancer + `/health/endtoend` monitoring; NO automatic DB failover | Independent — AP failure doesn't affect EU; simpler failure domains |
| **Complexity** | High: shared DB, shared Redis, cross-region network, TLS unmanaged | Lower: two independent stacks, mirroring config per repo |
| **Operational risk** | Single DB is SPOF; geo-replication loss = unreplicated blobs permanently lost | DB failure isolated to one region; no cross-region blast radius |
| **GDPR / data residency** | All blobs visible to all regions — potential compliance risk for EU data | EU blobs stay in EU unless explicitly mirrored — cleaner compliance posture |
| **Recommended for** | Teams with strong network infra, existing cross-region object storage, centralised ops | Teams prioritising resilience isolation, data residency compliance, or limited cross-region bandwidth |

### Constraints (fill in your actual values)
- [ ] **ADR-1** Document: Does cross-region network exist with sufficient bandwidth for async blob replication?
- [ ] **ADR-2** Document: Are EU data residency / GDPR constraints in scope? (affects geo-replication blob visibility)
- [ ] **ADR-3** Document: Is a Global Load Balancer available to route `/health/endtoend` traffic between regions?
- [ ] **ADR-4** Document: Is a single shared PostgreSQL instance with cross-region replication feasible?
- [ ] **ADR-5** Document: Is the existing object storage (S3/Ceph) accessible from BOTH regions simultaneously?

### Decision
- [ ] **ADR-6** Record chosen option with rationale in `adr-001-quay-global-vs-regional.md`
- [ ] **ADR-7** Record rejected option with explicit reasons
- [ ] **ADR-8** Identify consequences: operational runbooks required, monitoring requirements, cost implications

> **Recommendation to document:** If GDPR data residency applies OR cross-region object storage is not feasible → **Option B (Independent + Mirroring)**. If a Global LB and shared cross-region object storage exist → **Option A (Geo-Replication)**. Geo-replication is NOT recommended if you cannot meet all constraints listed in Section 16.2 of the Red Hat Quay documentation, as the risk of silent data loss from unreplicated blobs is significant.

---

## Phase 0: Discovery & Pre-Migration Inventory

> **Competency:** Platform ownership, risk assessment, baseline measurement
> **Salary Lever:** "I owned the discovery — I knew exactly what we had before touching anything."

- [ ] **P0-1** Inventory current VM-based Quay: version, repo count, image count, total blob storage size, active users, integrations (LDAP/OIDC, CI/CD clients, OCP pull secrets)
- [ ] **P0-2** **Measure current storage usage baseline:**
  ```bash
  # Total storage on VM
  du -sh /var/lib/quay/storage/
  # Or for S3-backed storage:
  aws s3 ls s3://QUAY_BUCKET --recursive --human-readable --summarize | tail -2
  ```
  Record in `storage-baseline.md`: total size, top 5 largest repos, number of untagged images
- [ ] **P0-3** Audit and identify storage waste:
  - Count untagged/dangling images per repository (`GET /api/v1/repository/{repo}/tag?onlyActiveTags=false`)
  - Identify repos with >30 tags that have no retention policy
  - Identify stale robot accounts not used in >90 days
- [ ] **P0-4** Document current `config.yaml`: all non-default settings, storage backend, LDAP/OIDC config, `TIME_MACHINE_DURATION`, `GARBAGE_COLLECTION` settings
- [ ] **P0-5** Map all upstream consumers: CI/CD pipelines, OCP clusters, developer tools — store in `consumers-map.md`
- [ ] **P0-6** Clarify licensing situation:
  - Is current VM running Red Hat Quay (paid) or Project Quay (community)?
  - Does the organisation have **Red Hat OpenShift Platform Plus** (includes Quay entitlement)?
  - Document in `license-analysis.md`: current cost vs post-migration cost
- [ ] **P0-7** Identify managed vs unmanaged Operator component decisions (PostgreSQL, Redis, Clair, object storage, TLS)
- [ ] **P0-8** Set migration approach: in-place Operator migration vs parallel new deployment + cutover

---

## Phase 1: Pre-Migration Backup & Rollback Plan

> **Competency:** Disaster recovery ownership, zero-risk migration design
> **Salary Lever:** "I defined and owned the rollback plan — no migration without a tested escape hatch."

- [ ] **P1-1** Back up Quay `config.yaml` from VM
- [ ] **P1-2** Perform PostgreSQL full dump:
  ```bash
  pg_dump -h DB_HOST -p 5432 -d QUAY_DB -U QUAY_USER -W -O > /tmp/quay-backup/quay-db-$(date +%Y%m%d).sql
  ```
- [ ] **P1-3** Back up all image blobs from current object storage:
  ```bash
  aws s3 sync --endpoint-url https://STORAGE_ENDPOINT s3://QUAY_BUCKET/ /tmp/quay-backup/blobs/
  ```
- [ ] **P1-4** Validate backup integrity: verify SQL size, spot-check 3+ blobs, record total count
- [ ] **P1-5** Store backup to secondary location (not the same cluster)
- [ ] **P1-6** Write `rollback-runbook.md`: DNS revert steps, Quay VM restart procedure, estimated RTO

---

## Phase 2: Storage Reduction — Pre-Migration Cleanup

> **Competency:** Storage optimisation, quota governance, platform cost management
> **Salary Lever:** "I reduced storage before migration — demonstrating platform cost ownership, not just migration execution."

- [ ] **P2-1** **Enable and configure Time Machine aggressively before migration:**
  ```yaml
  # In config.yaml — set shorter retention before migration to reduce blob carry-over
  DEFAULT_TAG_EXPIRATION: 2w         # Was default 2 weeks — review with teams
  TAG_EXPIRATION_OPTIONS:
    - 0s
    - 1d
    - 1w
    - 2w
    - 4w
  ```
- [ ] **P2-2** Communicate tag expiry policy change to development teams; agree on retention windows per repo
- [ ] **P2-3** Identify and delete clearly unused repositories and organisations (get approval first)
- [ ] **P2-4** Trigger garbage collection by toggling `GARBAGE_COLLECTION` off/on in `config.yaml` (only way to manually force GC):
  ```yaml
  FEATURE_GARBAGE_COLLECTION: false  # disable briefly
  # restart Quay, then re-enable
  FEATURE_GARBAGE_COLLECTION: true   # re-enable triggers GC workers
  ```
- [ ] **P2-5** Monitor GC metrics on VM Quay:
  - `quay_gc_storage_blobs_deleted_total` — number of blobs deleted
  - `quay_gc_repos_purged_total` — repos cleaned up
- [ ] **P2-6** Re-measure storage after GC run; record in `storage-reduction-log.md`:
  - Before: `____ GB`
  - After GC: `____ GB`
  - Reduction: `____ GB (____%)`
- [ ] **P2-7** Clean up S3 multipart upload remnants (if applicable):
  ```bash
  # Quay handles this automatically if FEATURE_CLEAN_BLOB_UPLOAD_FOLDER: true
  # Verify this is set in config.yaml
  ```

---

## Phase 3: OpenShift Environment Preparation

> **Competency:** OpenShift platform administration, capacity planning, infra node isolation

- [ ] **P3-1** Validate OCP prerequisites: version ≥ 4.6, cluster-admin access, default StorageClass, minimum 8Gi RAM + 2 vCPUs available
- [ ] **P3-2** Create dedicated namespace: `oc new-project quay-enterprise`
- [ ] **P3-3** (Senior-level) Schedule Quay on dedicated infra nodes:
  ```bash
  oc label node <node> node-role.kubernetes.io/infra=
  oc adm taint nodes -l node-role.kubernetes.io/infra node-role.kubernetes.io/infra=reserved:NoSchedule
  ```
  Annotate the `quay-enterprise` namespace with matching node-selector and tolerations
- [ ] **P3-4** Prepare object storage backend:
  - If NooBaa/ODF: verify `ObjectBucketClaim` API available
  - If external S3/Ceph: prepare `DISTRIBUTED_STORAGE_CONFIG` credentials
- [ ] **P3-5** Prepare `configBundleSecret` from backed-up config (carry over SECRET_KEY, LDAP/OIDC, storage, feature flags)
- [ ] **P3-6** Verify LDAP/OIDC endpoints are reachable from OCP cluster network

---

## Phase 4: Quay Operator Installation & Deployment

> **Competency:** Operator lifecycle management, GitOps, declarative platform management

- [ ] **P4-1** Install Red Hat Quay Operator from OperatorHub:
  - Update channel: `stable-3.x`
  - Approval strategy: **Manual** (review upgrades before applying in production)
- [ ] **P4-2** Create `QuayRegistry` CR with correct managed/unmanaged flags per ADR-001 decision
- [ ] **P4-3** Monitor operator reconciliation: `oc describe quayregistry -n quay-enterprise`
- [ ] **P4-4** Verify all pods reach `Running` status
- [ ] **P4-5** Test fresh instance pre-data migration: login, push/pull a test image

---

## Phase 5: Data Migration (VM → Operator)

> **Competency:** Zero-downtime migration execution, database operations, data integrity

- [ ] **P5-1** Scale down Quay Operator and app deployments:
  ```bash
  oc scale --replicas=0 deployment quay-operator.<version> -n openshift-operators
  oc scale --replicas=0 deployment <quay-app> <quay-mirror> -n quay-enterprise
  ```
- [ ] **P5-2** Copy and restore database backup into Operator-managed PostgreSQL pod
- [ ] **P5-3** Sync image blobs to new storage backend (NooBaa or external S3/Ceph)
- [ ] **P5-4** Patch `QuayRegistry` CR with updated `configBundleSecret`
- [ ] **P5-5** Scale Quay Operator and app pods back up
- [ ] **P5-6** Verify data integrity: login, pull pre-existing images, verify image digests match

---

## Phase 6: Post-Migration Storage Optimisation (Operator)

> **Competency:** Platform cost management, quota governance, Day-2 storage operations
> **Salary Lever:** "Storage reduction was a migration driver — I implemented and measured it."

- [ ] **P6-1** **Enable Quota Management** (Quay Operator ≥ 3.7):
  ```yaml
  # In configBundleSecret config.yaml
  FEATURE_QUOTA_MANAGEMENT: true
  DEFAULT_SYSTEM_REJECT_QUOTA_BYTES: 10737418240  # 10 GiB per namespace default
  ```
- [ ] **P6-2** Set per-organisation storage quotas via Quay API or UI:
  ```bash
  # Set quota for each org — start with generous limit, tighten over time
  curl -X PUT https://QUAY_HOST/api/v1/organization/ORG_NAME/quota     -H "Authorization: Bearer TOKEN"     -d '{"limit_bytes": 5368709120}'  # 5 GiB example
  ```
- [ ] **P6-3** Configure `TAG_EXPIRATION_OPTIONS` and `DEFAULT_TAG_EXPIRATION` per team agreement
- [ ] **P6-4** Verify continuous garbage collection is active on Operator deployment:
  - Check GC worker is running: `oc logs -l quay-component=quay-app -n quay-enterprise | grep gcworker`
  - Verify `quay_gc_storage_blobs_deleted_total` metric is incrementing
- [ ] **P6-5** Configure `FEATURE_CLEAN_BLOB_UPLOAD_FOLDER: true` to auto-clean stale multipart uploads
- [ ] **P6-6** Set `PUSH_TEMP_TAG_EXPIRATION_SEC: 3600` (1 hour) — prevents premature blob deletion during pushes
- [ ] **P6-7** **Measure storage after 2 weeks of operation:**
  - Record in `storage-reduction-log.md`:
    - VM baseline: `____ GB`
    - Post-migration (day 0): `____ GB`
    - Post-GC + quota (day 14): `____ GB`
    - Total reduction: `____ GB (____%)` ← **use this number in salary review**

---

## Phase 7: Regional Architecture Implementation (per ADR-001)

> **Competency:** Multi-region platform design, geo-replication or mirroring, data residency awareness

### Option A: Geo-Replication (if ADR-001 selects Global Quay)

- [ ] **P7A-1** Deploy external PostgreSQL (required — Operator-managed DB cannot be shared across clusters)
  ```sql
  CREATE DATABASE quay;
  \c quay; CREATE EXTENSION IF NOT EXISTS pg_trgm;
  ```
- [ ] **P7A-2** Deploy shared Redis instance accessible from both EU and AP clusters (port 6379 open)
- [ ] **P7A-3** Create two object storage buckets: `quay-storage-eu` and `quay-storage-ap`
- [ ] **P7A-4** Create shared `configBundleSecret` with both storage backends defined:
  ```yaml
  DISTRIBUTED_STORAGE_CONFIG:
    eu-storage:
      - RHOCSStorage
      - access_key: EU_KEY
        secret_key: EU_SECRET
        bucket_name: quay-storage-eu
        hostname: EU_STORAGE_ENDPOINT
    ap-storage:
      - RHOCSStorage
      - access_key: AP_KEY
        secret_key: AP_SECRET
        bucket_name: quay-storage-ap
        hostname: AP_STORAGE_ENDPOINT
  DISTRIBUTED_STORAGE_DEFAULT_LOCATIONS:
    - eu-storage
    - ap-storage
  DISTRIBUTED_STORAGE_PREFERENCE: eu-storage   # override per cluster with env var
  FEATURE_STORAGE_REPLICATION: true
  ```
- [ ] **P7A-5** Deploy EU `QuayRegistry` CR with `QUAY_DISTRIBUTED_STORAGE_PREFERENCE=eu-storage` env override; set `tls: managed: false`, supply custom TLS cert in config bundle
- [ ] **P7A-6** Deploy AP `QuayRegistry` CR with `QUAY_DISTRIBUTED_STORAGE_PREFERENCE=ap-storage` env override
- [ ] **P7A-7** Configure Global Load Balancer (GLB) with `/health/endtoend` health check on both clusters
- [ ] **P7A-8** Run initial blob backfill replication:
  ```bash
  oc rsh <quay-app-pod> -n quay-enterprise
  python -m util.backfillreplication
  ```
- [ ] **P7A-9** Verify geo-replication: push image to EU → confirm blob appears in AP storage bucket
- [ ] **P7A-10** Document geo-replication failure mode in `runbook-georep-failover.md`:
  - AP storage failure → shut down AP Quay deployment → GLB redirects to EU
  - **Important:** There is NO automatic failover — this is a manual procedure

### Option B: Independent Instances + Repository Mirroring (if ADR-001 selects Regional)

- [ ] **P7B-1** Deploy two fully independent `QuayRegistry` instances: one in EU OCP cluster, one in AP OCP cluster
- [ ] **P7B-2** Ensure both have independent PostgreSQL + Redis + object storage (no shared components)
- [ ] **P7B-3** Identify repositories that need cross-region availability (not all repos need mirroring)
- [ ] **P7B-4** Configure Repository Mirroring for critical repos (EU → AP):
  - Go to Quay UI → Repository → Settings → Mirroring
  - Set: External Registry URL (EU endpoint), robot account credentials, sync interval
  - Enable: "Verify TLS" and "Force sync" for initial population
- [ ] **P7B-5** Verify mirror sync: push image to EU Quay → wait for sync interval → pull from AP Quay
- [ ] **P7B-6** Set up mirroring robot accounts with read-only permissions on source registry
- [ ] **P7B-7** Document which repos are mirrored and why in `mirror-registry.md`
- [ ] **P7B-8** Document EU data residency posture: AP-only repos stay in AP; EU-only repos stay in EU unless explicitly mirrored — record in `data-residency.md`

---

## Phase 8: Observability & SLOs

> **Competency:** SRE observability practice, SLO ownership, error budget management

- [ ] **P8-1** Enable monitoring component in `QuayRegistry` CR: `monitoring: managed: true`
- [ ] **P8-2** Define SLOs in `quay-slos.yaml`:
  ```yaml
  # Availability: 99.9% of image pull requests succeed (30-day window)
  # Latency: p95 pull response time < 2s
  # Push Success: 99.5% of pushes succeed
  # Clair Scan: 95% of images scanned within 10 min of push
  # Storage: Quota utilisation alert at 80% per namespace
  ```
- [ ] **P8-3** Build Grafana/OpenShift dashboard: push/pull rate, error rate, GC metrics (`quay_gc_storage_blobs_deleted_total`), quota utilisation
- [ ] **P8-4** Create alerts:
  - Quay pod CrashLoopBackOff
  - PostgreSQL storage > 80%
  - Error budget burn rate > 5% in 1 hour (fast burn)
  - Clair DB staleness (vulnerability DB not updated in > 24 hours)
  - **Namespace quota utilisation > 80%** (storage-specific to migration driver)
- [ ] **P8-5** For geo-replication (Option A): alert on replication lag — monitor blob replication queue depth
- [ ] **P8-6** Export all alert rules and dashboards as YAML (Observability as Code)
- [ ] **P8-7** Write `runbook-quay-alerts.md` — one runbook entry per alert

---

## Phase 9: Security Hardening

> **Competency:** Supply chain security, registry hardening, least-privilege RBAC

- [ ] **P9-1** Enable Clair vulnerability scanning; verify CVE database auto-updates
- [ ] **P9-2** Integrate Clair scan results into CI/CD: block images with CRITICAL CVEs from promotion
- [ ] **P9-3** Audit and harden robot accounts: remove unused accounts, enforce least-privilege (read-only vs read-write)
- [ ] **P9-4** Implement `FEATURE_REQUIRE_TEAM_INVITE: true` to prevent unsanctioned org access
- [ ] **P9-5** Verify TLS enforced end-to-end; confirm cert expiry monitoring alert
- [ ] **P9-6** Run OCP Security Context Constraints audit on Quay pods
- [ ] **P9-7** Document security baseline in `security-baseline.md`: findings before vs after hardening

---

## Phase 10: Day-2 Operations & Documentation

> **Competency:** Platform lifecycle ownership, team enablement, engineering leadership

- [ ] **P10-1** Write `operations-runbook.md`: add org, reset robot account, approve operator upgrade, trigger GC monitoring, adjust quota
- [ ] **P10-2** Write `upgrade-procedure.md`: how to review and approve Quay Operator `InstallPlan` upgrades (Manual approval strategy)
- [ ] **P10-3** Implement automated daily backup CronJob (pg_dump + S3 sync); alert if job fails
- [ ] **P10-4** Write `incident-response-quay.md`: Quay down, storage full (quota breach), Clair not scanning, geo-rep lag, auth failure
- [ ] **P10-5** Write `onboarding-guide.md`: how a new developer requests access, creates robot account, integrates with CI/CD
- [ ] **P10-6** Conduct team knowledge transfer session; record outcome in `kt-log.md`

---

## 💰 Salary Uplift Evidence Map

| Deliverable | Competency Signal | Talking Point |
|---|---|---|
| ADR-001 (Global vs Regional) | Principal architecture ownership | "I drove the multi-region design decision with documented trade-offs." |
| `storage-reduction-log.md` | Platform cost management | "I reduced storage by \_\_\_% — this was the migration's primary business driver." |
| Storage quota + GC config | Operational excellence | "I implemented quota governance so storage never creeps unchecked again." |
| Geo-replication OR mirroring | Multi-region platform design | "I designed and implemented the cross-region registry architecture." |
| SLOs + burn rate alerts | SRE practice | "I own the reliability contract for the container registry." |
| `data-residency.md` | Compliance-aware design | "I addressed GDPR data residency constraints in the architecture." |
| Operator GitOps + Manual upgrades | Platform engineering maturity | "All config changes are auditable and go through review before applying." |
| Full runbook suite | Engineering leadership | "I made this platform maintainable for the whole team, not just me." |

---

## 📊 Competency Coverage

| Skill Area | Phases | Key Evidence |
|---|---|---|
| Architecture Decision Making | ADR-001 | ADR document with trade-off analysis |
| Storage Optimisation | P2, P6 | Measured before/after storage reduction % |
| OpenShift Operator Management | P3, P4, P7 | QuayRegistry CR, lifecycle management |
| Multi-region Design | P7A or P7B | Geo-rep or mirroring implementation |
| SRE Practice | P8 | SLOs, error budgets, burn-rate alerts |
| Data Residency / Compliance | ADR-001, P7B | `data-residency.md` |
| Supply Chain Security | P9 | Clair, RBAC, CI/CD gate |
| Engineering Leadership | P10 | Runbooks, KT, onboarding guide |

---

*Last updated: 2026-02-27 | Project: quay-vm-to-operator-migration-v2*
