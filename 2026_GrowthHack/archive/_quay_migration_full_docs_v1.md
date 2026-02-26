# 📄 Quay Migration: Full Documentation Suite
## Restricted Network Environment Edition

> **Context:** This project runs in a **restricted/air-gapped network**. External images cannot be pulled directly by workloads. All external content must flow through a controlled ingestion pipeline into the internal Quay registry.
>
> **Purpose of this document:** Complete list of all Business Impact, Architecture Decision Records, and Platform Engineering documentation you should produce — both for the migration project and as ongoing platform ownership evidence.

---

## 📁 Suggested Repository Structure

```
quay-migration/
├── docs/
│   ├── business-impact/          # BID-xxx  (stakeholder/management audience)
│   ├── adr/                      # ADR-xxx  (technical decision records)
│   ├── platform/                 # Platform engineering reference docs
│   ├── runbooks/                 # Operational runbooks (incident response, Day-2 ops)
│   ├── security/                 # Security posture, compliance, network controls
│   └── architecture/             # Diagrams, network topology, data flows
```

---

---

# 🏢 SECTION 1: Business Impact Documents (BIDs)

> *Audience: Engineering Manager, CTO, Finance, Change Advisory Board*
> *These documents prove you think in business outcomes, not just technical tasks.*

---

### BID-001: Migration Business Case
- Executive summary: why migrate now?
- Current state problems (storage cost, licensing, VM operational toil, restricted network risks)
- Proposed solution: Quay Operator on OpenShift as the **single internal registry gateway**
- Cost analysis: current licensing + storage vs post-migration costs; net saving per year (`$____`)
- OpenShift Platform Plus entitlement value (Quay included at no extra cost — document this)
- Risk summary (reference BID-004)
- Recommendation and decision requested from stakeholders

---

### BID-002: Storage Reduction Report
- Before migration baseline: total blob storage, top 5 largest repos, untagged image count
- Storage reduction actions: GC runs, Time Machine tuning, org/repo cleanup
- After migration measurements: Day 0 / Day 14 / Day 30
- **Total reduction: `____ GB (____%)`** — use this number in salary review
- Quota policy implemented: limits per org, breach alerts configured

---

### BID-003: License Cost Avoidance Summary
- Current VM Quay licensing model and annual cost (`$____`)
- Post-migration model: OCP Platform Plus entitlement or Project Quay community
- Annual cost avoidance: `$____`
- Support tier trade-off documented (Red Hat enterprise support vs community)

---

### BID-004: Migration Risk Register

| Risk ID | Description | Likelihood | Impact | Mitigation |
|---|---|---|---|---|
| R-001 | Data loss during DB restore | Low | Critical | Full pg_dump + offsite backup; tested restore on staging |
| R-002 | Image blob not replicated post-migration | Medium | High | Blob count + digest validation before/after |
| R-003 | LDAP/OIDC auth failure after cutover | Medium | High | Auth tested in staging with real accounts |
| R-004 | CI/CD pipelines broken after DNS cutover | Medium | High | Consumer map; pipeline smoke tests post-cutover |
| R-005 | External image ingestion fails (firewall block) | High | Critical | Firewall rules tested before migration; fallback manual transfer |
| R-006 | Proxy/allowlist not updated before go-live | Medium | High | Network change request submitted at Phase 0 |
| R-007 | Quay Operator upgrade breaks config | Low | High | Manual approval strategy; upgrade tested in non-prod |
| R-008 | Storage quota breach disrupts CI/CD push | Low | High | 80% alert; grace period before hard reject |
| R-009 | Geo-replication lag causes stale AP pulls | Medium | Medium | Replication monitoring; GLB health check failover |
| R-010 | Corporate CA cert not trusted by Quay mirror | Medium | High | CA bundle injection tested in non-prod |

---

### BID-005: Maintenance Window & Cutover Plan
- Agreed maintenance window: `date / time / timezone`
- Go/No-Go criteria: what must be true before proceeding
- Step-by-step cutover checklist
- Communication plan: who to notify (teams, on-call, network/firewall team)
- Rollback decision point: if validation fails at `T + 30 min` → rollback
- Rollback procedure: DNS revert, VM Quay restart, estimated RTO

---

### BID-006: Post-Migration Operational Handover Report
- Summary of what was migrated, when, and by whom
- New architecture overview (link to diagram)
- Key operational differences: VM service → Operator reconciliation model
- SLOs defined and dashboard location
- Known follow-up actions outstanding
- Stakeholder sign-off section

---

---

# 🧭 SECTION 2: Architecture Decision Records (ADRs)

> Standard format: **Title / Status / Context / Decision / Consequences / Alternatives Rejected**

---

### ADR-001: Global Quay (Geo-Replication) vs Independent Regional Instances (EU + AP)
- **Context:** Business requires registry availability in EU and AP. Geo-replication shares one PostgreSQL + Redis with region-local object storage; independent instances share nothing and use repository mirroring.
- **Key constraints:** Cross-region object storage accessibility, GDPR data residency, Global LB availability, cross-region PostgreSQL feasibility
- **Decision options:** Option A (geo-rep) vs Option B (independent + mirroring)
- **Recommended decision tree:** If GDPR data residency applies OR cross-region object storage is not feasible → Option B. If Global LB + shared cross-region DB + storage exist → Option A.
- **Consequences:** Document failure modes, operational overhead, compliance posture, cost per option

---

### ADR-002: Managed vs Unmanaged Operator Components
- **Context:** Quay Operator manages PostgreSQL, Redis, Clair, object storage, TLS, HPA, Route, monitoring independently.
- **Decision:** Document `managed: true/false` choice for each component with rationale
- **Consequences:** Unmanaged components = team owns lifecycle, upgrades, and HA

---

### ADR-003: Operator Update Channel & Approval Strategy
- **Context:** `Automatic` vs `Manual` InstallPlan approval. Automatic upgrades may break production config.
- **Decision:** `Manual` for production; `Automatic` acceptable for non-prod
- **Consequences:** Team must review release notes and approve InstallPlan before each upgrade

---

### ADR-004: Object Storage Backend Selection
- **Context:** Storage reduction is a migration driver. NooBaa/ODF vs existing S3/Ceph vs new Ceph cluster.
- **Decision:** Document chosen option with cost comparison per GB/month
- **Consequences:** Unmanaged S3/Ceph → team owns bucket lifecycle, versioning, cross-region replication

---

### ADR-005: Tag Retention & Time Machine Policy
- **Context:** `TIME_MACHINE_DURATION` drives how long deleted blobs persist before GC. Default 14 days.
- **Decision:** Production: `____ days`; Non-prod/dev: `____ days` (shorter — dev images are ephemeral)
- **Consequences:** Shorter duration = less storage, smaller recovery window. Teams must agree before change.

---

### ADR-006: Repository Mirroring Scope (if ADR-001 → Option B)
- **Context:** Not all repos need cross-region replication. Mirroring everything defeats storage reduction goal.
- **Decision:** Define tiers:
  - Tier 1 — Mirror everywhere: production base images, platform tooling
  - Tier 2 — Mirror on request: team app images needed cross-region
  - Tier 3 — No mirror: dev/test, CI intermediate images, single-region workloads
- **Consequences:** Tier 1/2 require robot accounts; sync interval defined per tier

---

### ADR-007: CI/CD Integration Strategy Post-Migration
- **Context:** All CI/CD pipelines, OCP pull-secrets, and developer tools reference the current Quay hostname.
- **Decision options:** Same DNS hostname (transparent), update all consumers, or parallel run hypercare period
- **Consequences:** Document chosen DNS cutover strategy, TTL, and pipeline smoke test plan

---

### ADR-008: External Image Ingestion Strategy (Restricted Network)
*This is the most operationally critical ADR in a restricted network environment.*

- **Context:** Workloads cannot pull images directly from `registry.redhat.io`, `registry.access.redhat.com`, `cdn.quay.io`, or any external registry. All external content must enter the network through a controlled, approved pipeline.
- **Decision options:**

  | Option | Description | Pros | Cons |
  |---|---|---|---|
  | A: Quay Repository Mirroring (pull-through) | Quay mirrors from external registry on schedule; internal consumers pull from Quay | Automated; versioned; audit trail | Requires outbound firewall allowlist to external registries |
  | B: Bastion/Jump Host Manual Transfer | Engineer pulls image on internet-connected bastion, saves as tar, transfers via secure file transfer to internal network, pushes to Quay | No persistent outbound connection needed | Manual toil; no automation; human error risk |
  | C: oc-mirror / oc adm catalog mirror | Use `oc-mirror` tool to create an image mirror set, export to portable media or file, import to internal Quay | Fully air-gapped; supports OCP operator catalogs | Complex toolchain; offline media management overhead |
  | D: Authenticated HTTP Forward Proxy | Route outbound registry traffic through a corp proxy with allowlist + SSL inspection | Centralised control; full audit log | Proxy becomes SPOF; SSL inspection may break image pull auth |

- **Recommended decision:** Document which option(s) are used per image class (RH base images, operator catalogs, third-party images)
- **Consequences:** Each option requires different network change requests, robot accounts, and monitoring

---

### ADR-009: Network & Firewall Allowlist for External Registry Access
*Critical for restricted network — defines exactly what outbound connections are permitted.*

- **Context:** Red Hat container registries changed their CDN infrastructure; image content is served from Quay.io CDN hosts. Without the correct allowlist, image pulls fail silently with `ImagePullBackOff`.
- **Required outbound hostnames (TCP 443):**
  ```
  registry.redhat.io          # Red Hat product images
  registry.access.redhat.com  # Red Hat community images
  cdn.quay.io                 # CDN for image blobs (required since 2023)
  cdn01.quay.io
  cdn02.quay.io
  cdn03.quay.io
  cdn04.quay.io
  cdn05.quay.io
  cdn06.quay.io
  quay.io                     # If mirroring from quay.io directly
  ```
- **Decision:** Which hosts are permitted outbound from which network zone (bastion only, Quay mirror pod only, or broader)?
- **Consequences:** Recommend hostname-based rules, NOT IP-based — Red Hat CDN IPs change and are not published as a static list.

---

### ADR-010: Corporate CA Certificate Trust for Internal Quay TLS
- **Context:** Internal Quay uses a TLS certificate signed by the corporate CA. All consumers (OCP clusters, CI/CD tools, developer workstations, Quay mirror jobs) must trust this CA to avoid `x509: certificate signed by unknown authority` errors.
- **Decision:** Document CA injection strategy:
  - OCP cluster: inject CA into `image.config.openshift.io/cluster` `additionalTrustedCA` field
  - Quay mirror robot: inject CA bundle into `configBundleSecret` as `extra_ca_cert_<name>.crt`
  - CI/CD agents: document how CA is distributed to build agents
  - Developer workstations: document OS-level CA trust procedure
- **Consequences:** Any new consumer must follow the CA trust onboarding checklist

---

---

# 🛠️ SECTION 3: Platform Engineering Documentation

> *These are the documents that demonstrate you OWN the platform — not just ran a migration task.*
> *This is what separates a Senior PE from a mid-level engineer in salary conversations.*

---

## 3A: Architecture & Design Documents

### ARCH-001: Platform Architecture Overview
- System context diagram: how Quay fits into the organisation's platform (OCP clusters, CI/CD, developers, external registries)
- Component diagram: Quay Operator internals (quay-app, postgres, redis, clair, mirror, HPA, Route)
- Network topology diagram:
  - Internal zones (OCP cluster network, storage network, management network)
  - External zone (internet-facing registries)
  - Ingestion path: external registry → firewall/proxy → Quay mirror → internal consumers
- Data flow diagram: image push path (CI/CD → Quay), image pull path (OCP workload → Quay), external ingestion path

### ARCH-002: Network Topology & Zones Document
- Network zone definitions: DMZ, internal cluster, management, storage
- Firewall rule inventory (reference ARCH-003)
- DNS entries required: internal Quay FQDN, mirror endpoint hostname
- Load balancer / OpenShift Route configuration

### ARCH-003: Firewall & Network Security Rules Register

| Rule ID | Direction | Source Zone | Destination | Port | Protocol | Purpose | Approved By | Review Date |
|---|---|---|---|---|---|---|---|---|
| FW-001 | Outbound | Quay mirror pod | registry.redhat.io | 443 | TCP/HTTPS | RH image ingestion | `____` | `____` |
| FW-002 | Outbound | Quay mirror pod | cdn.quay.io (01-06) | 443 | TCP/HTTPS | RH CDN blob pulls | `____` | `____` |
| FW-003 | Outbound | Bastion host | registry.redhat.io | 443 | TCP/HTTPS | Manual image transfer | `____` | `____` |
| FW-004 | Inbound | OCP cluster nodes | Internal Quay FQDN | 443 | TCP/HTTPS | Image pulls by workloads | `____` | `____` |
| FW-005 | Inbound | CI/CD agents | Internal Quay FQDN | 443 | TCP/HTTPS | Image push from pipelines | `____` | `____` |
| FW-006 | Internal | Quay pods | PostgreSQL | 5432 | TCP | DB connection | `____` | `____` |
| FW-007 | Internal | Quay pods | Redis | 6379 | TCP | Cache/queue | `____` | `____` |
| FW-008 | Internal | Clair pods | Clair DB | 5432 | TCP | Vulnerability DB | `____` | `____` |

---

## 3B: External Image Ingestion Process Documents

### PROC-001: External Image Ingestion Request & Approval Process
*This is your most important governance document in a restricted network.*

**Process flow:**
```
Developer/Team Request
        ↓
Submit Image Ingestion Request Form (PROC-001-form.md)
        ↓
Platform Engineer Review:
  - Is image from a trusted registry? (approved list)
  - Is image version pinned? (no :latest tags)
  - Has image passed Clair vulnerability scan?
        ↓
Security Review (if new external registry source):
  - Firewall change request (if new hostname needed)
  - Network team approval
        ↓
Platform Engineer: Configure Quay Repository Mirroring
        ↓
Verify: Mirror sync completes; image pullable from internal Quay
        ↓
Notify requestor: internal Quay path to use
        ↓
Update Approved Image Catalogue (PROC-003)
```

**Document contents:**
- Who can request external image ingestion (any engineer, approval required from PE team)
- Approved external registries list (registry.redhat.io, quay.io, docker.io, etc.)
- Prohibited sources list (unknown/unvetted registries)
- SLA for request fulfilment: `____ business days`
- Escalation path if urgent (production incident requires image immediately)
- Review cadence: quarterly audit of all active mirrors

### PROC-002: Quay Repository Mirroring Configuration Runbook
- How to create a mirror repository in Quay UI or via API
- Mirror configuration fields: source URL, robot account credentials, sync interval, tags filter, TLS verification
- How to set `QUAY_DISTRIBUTED_STORAGE_PREFERENCE` for geo-replicated mirrors
- How to trigger a manual sync: `POST /api/v1/repository/{repo}/mirror/sync-now`
- How to verify sync success: check `Last Sync` timestamp and last sync status in Quay UI
- How to troubleshoot failed mirror sync:
  - Check Quay mirror pod logs: `oc logs -l quay-component=quay-mirror -n quay-enterprise`
  - Check firewall connectivity from Quay mirror pod: `oc debug -n quay-enterprise` → `curl https://registry.redhat.io/v2/`
  - Check if corporate CA cert is trusted (common failure mode)
- How to decommission a mirror (remove sync config, delete robot account)

### PROC-003: Approved Internal Image Catalogue
*Single source of truth for what images are available internally and where they came from.*

| Internal Quay Path | Source Registry | Source Image | Sync Interval | Owner Team | Approved Date | Review Date |
|---|---|---|---|---|---|---|
| `quay.internal/rh-base/ubi9` | registry.redhat.io | `ubi9/ubi:latest` | Daily | Platform | `____` | `____` |
| `quay.internal/rh-base/ubi9-minimal` | registry.redhat.io | `ubi9/ubi-minimal` | Daily | Platform | `____` | `____` |
| `quay.internal/ocp/pause` | quay.io/openshift | `pause:latest` | Weekly | Platform | `____` | `____` |
| `quay.internal/app/nginx` | registry.redhat.io | `rhel9/nginx-120` | Weekly | App Team A | `____` | `____` |

### PROC-004: Bastion Host Manual Image Transfer Procedure
*(Used when Quay repository mirroring is not available or for one-off emergency ingestion)*

```bash
# Step 1: On internet-connected bastion host
podman login registry.redhat.io
podman pull registry.redhat.io/ubi9/ubi:9.4
podman save -o /tmp/ubi9-9.4.tar registry.redhat.io/ubi9/ubi:9.4

# Step 2: Secure transfer to internal network host
scp /tmp/ubi9-9.4.tar internal-host:/tmp/

# Step 3: On internal host — push to Quay
podman login quay.internal --tls-verify=true
podman load -i /tmp/ubi9-9.4.tar
podman tag registry.redhat.io/ubi9/ubi:9.4 quay.internal/rh-base/ubi9:9.4
podman push quay.internal/rh-base/ubi9:9.4
```
- Document: who has bastion host access (restricted — not all engineers)
- Document: file transfer mechanism (SFTP? internal file transfer system?)
- Log all manual transfers in `PROC-003` (Approved Image Catalogue)

---

## 3C: Security & Compliance Documents

### SEC-001: Supply Chain Security Policy
- All workloads must pull images from internal Quay only — no direct external pulls
- All external images must go through the PROC-001 approval process
- All images scanned by Clair before promotion to production namespaces
- CRITICAL CVE policy: images with CRITICAL CVEs blocked from production (CI/CD gate)
- Image signing policy (cosign or Quay robot account restrictions)

### SEC-002: Robot Account Registry & Audit Log
- Inventory of all Quay robot accounts: name, organisation, permissions, owning team, last used date
- Review cadence: quarterly — deactivate accounts not used in >90 days
- Naming convention: `team-purpose-env` (e.g. `ci-push-prod`, `ocp-pull-staging`)
- Credential rotation policy: rotate annually or on team member departure

### SEC-003: TLS Certificate Management
- Internal Quay FQDN and certificate details: issuing CA, expiry date, renewal owner
- Renewal runbook: how to update cert in `configBundleSecret` and trigger Operator reconciliation
- Alert: CloudWatch/Prometheus alert on cert expiry < 30 days
- Corporate CA distribution: how CA is pushed to OCP clusters, CI/CD agents, developer workstations

### SEC-004: LDAP/OIDC Integration Configuration Record
- Identity provider details: LDAP server hostname, bind DN, search base, group sync config
- OIDC provider: client ID, scopes, redirect URIs
- How to test auth: `ldapsearch` from Quay pod, test OIDC flow via browser
- Troubleshooting auth failures runbook

---

## 3D: Operational Runbooks

### RUN-001: Quay Day-2 Operations Runbook
Common tasks any PE team member should be able to perform:
- Add a new organisation and set storage quota
- Create a robot account and assign permissions
- Reset a user password / unlock an account
- Approve Quay Operator upgrade (review InstallPlan)
- Adjust `TIME_MACHINE_DURATION` via `configBundleSecret` update
- Trigger manual GC monitoring check
- Review and action Clair vulnerability scan results

### RUN-002: Incident Response — Quay Platform

| Scenario | First Response | Escalation |
|---|---|---|
| Quay pods CrashLoopBackOff | Check OCP events, pod logs; check `configBundleSecret` valid | Engage Platform Lead |
| Storage quota breach (push rejected) | Identify offending org; raise temporary quota; notify team to clean up | Storage team |
| Clair not scanning new images | Check clair-app pod logs; check vulnerability DB sync status | Platform Lead |
| External image pull blocked (firewall) | Check mirror pod logs; verify firewall rule active; emergency bastion transfer if urgent | Network/Security team |
| LDAP/OIDC auth failure | Test LDAP connectivity from pod; check bind account not locked | Identity team |
| Geo-replication lag > threshold | Check storage replication queue; verify cross-region network; failover GLB if needed | Platform Lead |
| Quay certificate expired | Follow SEC-003 cert renewal runbook | Platform Lead |

### RUN-003: Quay Operator Upgrade Procedure
- Check Red Hat Quay release notes for breaking changes
- Test upgrade in non-prod cluster first; validate all pods healthy
- Review and approve `InstallPlan` in production:
  ```bash
  oc get installplan -n openshift-operators
  oc patch installplan <name> --type merge --patch '{"spec":{"approved":true}}'
  ```
- Monitor reconciliation: `oc get pods -n quay-enterprise -w`
- Rollback: if upgrade fails, previous CSV is retained — document rollback steps

### RUN-004: New Developer Onboarding to Quay
- How to request Quay access (link to PROC-001 form or internal request system)
- How to log in (LDAP/OIDC SSO)
- How to find approved images in the internal catalogue (PROC-003)
- How to request a robot account for a CI/CD pipeline
- How to request a new external image be mirrored (PROC-001)
- What NOT to do: no direct pulls from external registries; no pushing of non-approved base images

---

## 3E: Capacity & SRE Documents

### SRE-001: Quay SLO Definition (`quay-slos.yaml`)
```yaml
slos:
  - name: availability
    description: "99.9% of image pull requests succeed (30-day window)"
    sli: "1 - (rate(quay_pull_errors_total[30d]) / rate(quay_pull_requests_total[30d]))"
    target: 99.9%

  - name: pull_latency
    description: "p95 image pull response time < 2 seconds"
    sli: "histogram_quantile(0.95, quay_pull_duration_seconds)"
    target: "< 2s"

  - name: push_success
    description: "99.5% of image push requests succeed"
    target: 99.5%

  - name: clair_scan_freshness
    description: "95% of new images scanned within 10 minutes of push"
    target: 95%

  - name: mirror_sync_lag
    description: "External image mirrors complete sync within 1 hour of schedule"
    target: 95%

  - name: storage_quota_headroom
    description: "No organisation exceeds 80% of quota (alert threshold)"
    target: "< 80% utilisation"
```

### SRE-002: Error Budget & Burn Rate Policy
- Error budget calculation per SLO (30-day rolling window)
- Fast-burn alert: budget consumed > 5% in 1 hour → page on-call immediately
- Slow-burn alert: budget consumed > 2% over 6 hours → ticket created
- Error budget freeze policy: if budget < 10% remaining → no non-critical changes to Quay platform

### SRE-003: Capacity Planning Report (Quarterly)
- Current storage utilisation trend (by org, by total)
- Projected storage exhaustion date at current growth rate
- CPU/memory utilisation of Quay pods (HPA trigger history)
- PostgreSQL storage growth rate and projected capacity
- Recommendations: quota adjustments, GC tuning, hardware procurement

---

## 💰 Salary Leverage: Document Ownership Summary

| Document | What It Proves |
|---|---|
| ADR-008 External Ingestion Strategy | "I designed the gateway between the internet and our internal platform." |
| ADR-009 Firewall Allowlist | "I defined and owned the network security boundary for the registry." |
| PROC-001 Approval Process | "I established governance — not just a technical setup." |
| PROC-003 Image Catalogue | "I created the single source of truth for what's running in production." |
| SEC-001 Supply Chain Security Policy | "I own the supply chain security posture." |
| ARCH-003 Firewall Rules Register | "Every network exception is documented, approved, and reviewed by me." |
| SRE-001 SLO Definition | "I defined the reliability contract for the most critical platform component." |
| Full runbook suite | "I made this platform operable by the whole team, not just me." |

---

*Last updated: 2026-02-27 | Project: quay-vm-to-operator-migration | Environment: Restricted Network*
