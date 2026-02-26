
## BID-007. Disaster Recovery Validation Report

**File:** `07-dr-validation.md`

* [ ] Defined RPO / RTO
* [ ] Backup strategy
* [ ] Restore test execution steps
* [ ] Measured recovery time
* [ ] Data integrity validation
* [ ] Lessons learned

---

### ADR-005: Backup and Restore Strategy

* Backup frequency decision
* Snapshot vs logical backup
* Restore testing cadence
* Storage of backup artifacts

---


## BID-008: Security & Compliance Impact

**File:** `08-security-impact.md`

* [ ] RBAC redesign
* [ ] OIDC/SSO integration
* [ ] Image vulnerability scanning improvements
* [ ] Immutable tag policy
* [ ] Supply chain risk mitigation
* [ ] Audit readiness improvements

---

### ADR-010: Security Hardening Model

* Authentication model decision
* RBAC scope control
* Vulnerability scanning enforcement
* Immutable tag policy decision


---

## BID-004. Performance Benchmark Report

**File:** `04-performance-benchmark.md`

* [ ] Baseline image pull latency
* [ ] Post-migration latency (p50 / p95)
* [ ] Push throughput testing results
* [ ] Load test methodology
* [ ] Scaling validation results
* [ ] Resource utilization comparison

---

### ADR-009: Global Quay Service vs Regional Quay Deployments


1. A single global Quay service serving both EU and AP regions?
OR
1. One independent Quay service deployed in each region?

Key concerns:
* Latency for image pulls
* Availability and blast radius
* Data sovereignty requirements
* Disaster recovery complexity
* Operational overhead
* Storage replication strategy

---

### ADR-008: Migration Strategy (Cutover Model)

* Blue/Green vs Big Bang
* Parallel deployment decision
* Rollback strategy
* Data sync method
* Downtime minimization approach



===================================================================================


## BID-001. Executive Summary – Platform Modernization

**File:** `01-executive-summary.md`

* [ ] Why migration was necessary
* [ ] Business risks of legacy VM-based Red Hat Quay
* [ ] Strategic alignment (cloud-native / OpenShift-first strategy)
* [ ] Summary of measurable improvements
* [ ] Before vs After comparison table

---

## BID-002. Risk Assessment & Mitigation Report

**File:** `02-risk-assessment.md`

* [ ] Identified single points of failure
* [ ] CI/CD dependency mapping
* [ ] Registry outage blast radius analysis
* [ ] Data integrity risks
* [ ] Security risks
* [ ] Mitigation strategies implemented
* [ ] Residual risk assessment

---

## BID-003. Reliability Improvement Report

**File:** `03-reliability-improvements.md`

* [ ] Defined SLOs and SLIs
* [ ] Availability baseline vs post-migration
* [ ] Error budget implementation
* [ ] MTTR before vs after
* [ ] Incident reduction metrics
* [ ] Alert noise reduction metrics


---

## BID-005. Cost Analysis Report

**File:** `05-cost-analysis.md`

* [ ] VM infrastructure cost (old model)
* [ ] OpenShift resource allocation cost
* [ ] Storage backend cost
* [ ] Operational effort cost (toil hours/month)
* [ ] Long-term scalability cost projection
* [ ] ROI estimation

---

## BID-006. Operational Efficiency Report

**File:** `06-operational-efficiency.md`

* [ ] Manual processes eliminated
* [ ] Automation introduced (GitOps, backups, monitoring)
* [ ] Weekly ops hours saved
* [ ] Deployment friction reduction
* [ ] Developer experience improvements


---

### ADR-002: High Availability Design

* Replica count decision
* Pod anti-affinity rules
* Route/Ingress redundancy
* Scaling model
* Trade-off between cost and redundancy

---

### ADR-003: Object Storage Backend Selection

* S3 vs ODF vs MinIO
* Durability guarantees
* Cost considerations
* Operational overhead
* Performance testing results

---

### ADR-004: Database Architecture Strategy

* External HA Postgres vs in-cluster DB
* Backup tooling decision
* Failover handling strategy
* RPO/RTO justification

---


### ADR-006: SLO Definition & Error Budget Policy

* Why 99.9% (or chosen number)
* Business impact justification
* Error budget burn alert design
* Trade-off between velocity and stability

---

### ADR-007: Observability Architecture

* Metrics collection approach
* Alert threshold design
* Log aggregation strategy
* Synthetic monitoring decision
* Alert noise control strategy


---

### ADR-001: Migration from VM-Based Quay to Operator Model

* Why move from VM deployment to Red Hat Quay Operator
* Operational limitations of legacy model
* Benefits of declarative reconciliation on OpenShift
* Trade-offs (complexity vs resilience)