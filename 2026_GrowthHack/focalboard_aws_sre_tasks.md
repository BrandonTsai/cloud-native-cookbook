# 🚀 Focalboard → Production-Grade AWS Platform: Side Project Task Board

> **Goal:** Transform [Focalboard OSS](https://github.com/mattermost-community/focalboard) into a Production-Grade, AI-Augmented Platform on AWS.
> **Why:** Build an end-to-end, demonstrable project covering Platform Engineering + Principal SRE competencies for interviews.

---

## 📋 How to Use This File

- `[ ]` = Not started | `[x]` = Done | `[~]` = In progress
- Each task maps to an **AWS skill** and/or **SRE competency**
- Ordered by recommended execution sequence (dependencies first)

---

## Phase 0: Foundation & Architecture Design

> **SRE Skill:** System design, documentation, tenets  
> **AWS Skill:** Account setup, IAM, VPC, cost planning

- [ ] **P0-1** Create AWS account with MFA; set up AWS Organizations (root + dev account)
- [ ] **P0-2** Create IAM roles with least-privilege for CI/CD, ECS task roles, and developer access
- [ ] **P0-3** Design VPC: 2 AZs, public/private subnets, NAT Gateway, Security Groups skeleton
- [ ] **P0-4** Write Architecture Decision Record (ADR) — why ECS Fargate over EKS, why CodePipeline over GitHub Actions
- [ ] **P0-5** Define SLOs upfront (latency p99, availability %, error rate) — store in `slos.yaml`
- [ ] **P0-6** Set up AWS Cost Budget alert ($50/month threshold) and enable Cost Explorer

---

## Phase 1: AI Refactor Focalboard + Containerize

> **SRE Skill:** Production readiness, 12-factor app principles  
> **AWS Skill:** ECR, Dockerfile best practices, multi-stage builds

- [ ] **P1-1** Install OrbStack on MacBook for local development.
- [ ] **P1-2** Clone Focalboard repo; run locally to understand architecture (Go backend + React frontend)
- [ ] **P1-3** Use AI (Copilot/Claude/Cursor) to refactor: add structured logging (JSON), health check endpoints `/healthz` and `/readyz`
- [ ] **P1-4** Write multi-stage Dockerfile (builder → distroless/minimal runtime image)
- [ ] **P1-5** Add `docker-compose.yml` for local dev with Postgres + Focalboard
- [ ] **P1-6** Create Amazon ECR repository; push image manually (first time) to validate
- [ ] **P1-7** Document container image size before/after optimization (demonstrate cost-awareness)

---

## Phase 2: Infrastructure as Code (IaC)

> **SRE Skill:** Toil reduction, repeatable environments  
> **AWS Skill:** Terraform or AWS CDK, ECS Fargate, ALB, RDS, Secrets Manager

- [ ] **P2-1** Write Terraform modules: VPC, Security Groups, ALB, RDS (Postgres), ECS Cluster
- [ ] **P2-2** Write ECS Task Definition + Fargate Service (CPU/memory sizing, env vars from SSM)
- [ ] **P2-3** Store DB credentials in AWS Secrets Manager; inject into ECS task via `secrets:` block
- [ ] **P2-4** Configure Application Load Balancer with HTTPS (ACM certificate) and target group health checks
- [ ] **P2-5** Set up RDS Postgres with Multi-AZ for HA; enable automated backups (7-day retention)
- [ ] **P2-6** Tag all resources with `Environment`, `Project`, `Owner` — SRE cost attribution practice
- [ ] **P2-7** `terraform plan` + `terraform apply` — first working deployment ✅

---

## Phase 3: CI/CD with AWS CodePipeline (AI-Augmented)

> **SRE Skill:** Deployment frequency, change failure rate (DORA metrics)  
> **AWS Skill:** CodePipeline, CodeBuild, CodeDeploy, ECR integration

- [ ] **P3-1** Create CodeCommit repo (or connect GitHub) as source stage
- [ ] **P3-2** Write `buildspec.yml` for CodeBuild: lint → test → docker build → push to ECR
- [ ] **P3-3** Add AI-generated test coverage step: use AI to write unit tests for refactored Go code
- [ ] **P3-4** Configure CodePipeline stages: Source → Build → Deploy (Blue/Green to ECS via CodeDeploy)
- [ ] **P3-5** Implement Blue/Green deployment with CodeDeploy + ALB listener rule switching
- [ ] **P3-6** Add AI-augmented PR review step: GitHub Action or Lambda that runs AI code review and posts PR comments
- [ ] **P3-7** Implement pipeline notifications via SNS → Slack/email on failure
- [ ] **P3-8** Track and document DORA metrics: deployment frequency, lead time, MTTR — store in `dora-metrics.md`

---

## Phase 4: Observability Stack (SLOs & Error Budgets)

> **SRE Skill:** SLO definition, error budgets, alerting philosophy (alert on symptoms not causes)  
> **AWS Skill:** CloudWatch, CloudWatch Synthetics, X-Ray, Dashboards

- [ ] **P4-1** Instrument Focalboard Go backend with AWS X-Ray SDK (distributed tracing)
- [ ] **P4-2** Set up CloudWatch Container Insights for ECS Fargate (CPU, memory, network metrics)
- [ ] **P4-3** Create CloudWatch Synthetic canary to simulate user login + board creation every 5 min
- [ ] **P4-4** Define and implement SLIs in CloudWatch Metrics:
  - Availability SLI: `% of synthetic checks succeeding`
  - Latency SLI: `p99 ALB TargetResponseTime < 500ms`
  - Error Rate SLI: `5xx errors / total requests < 0.1%`
- [ ] **P4-5** Build Error Budget dashboard in CloudWatch: remaining budget (%), burn rate alarm
- [ ] **P4-6** Create CloudWatch Alarm for error budget burn rate > 5% in 1 hour (fast burn alert)
- [ ] **P4-7** Create CloudWatch Alarm for error budget burn rate > 2% over 6 hours (slow burn alert)
- [ ] **P4-8** Export dashboard as code (CDK/Terraform) — "Observability as Code"
- [ ] **P4-9** Write `runbook-alerting.md` explaining each alert, severity, and response steps

---

## Phase 5: Chaos Engineering

> **SRE Skill:** Failure mode analysis, production resilience validation  
> **AWS Skill:** AWS FIS (Fault Injection Simulator), ECS task termination

- [ ] **P5-1** Write `chaos-hypothesis.md` — define steady state, hypothesis, blast radius for each experiment
- [ ] **P5-2** **Experiment 1 — ECS Task Kill:** Use AWS FIS to terminate 50% of ECS tasks; verify auto-recovery within SLO
- [ ] **P5-3** **Experiment 2 — AZ Failure Simulation:** Block traffic to one AZ via Security Group; verify failover
- [ ] **P5-4** **Experiment 3 — RDS Failover:** Trigger RDS Multi-AZ failover; measure downtime against SLO
- [ ] **P5-5** **Experiment 4 — High CPU Stress:** FIS CPU stress on ECS task; verify auto-scaling triggers
- [ ] **P5-6** Record each experiment result in `chaos-results.md`: was hypothesis validated? SLO breached?
- [ ] **P5-7** Fix any discovered reliability gaps; re-run experiments to confirm

---

## Phase 6: AI Incident Response

> **SRE Skill:** Incident management, post-mortems, MTTR reduction  
> **AWS Skill:** Lambda, EventBridge, SNS, Bedrock (Claude/Titan)

- [ ] **P6-1** Build Lambda function triggered by CloudWatch Alarm → SNS
- [ ] **P6-2** Lambda calls AWS Bedrock (Claude) with alarm context + recent CloudWatch logs → generates incident summary
- [ ] **P6-3** Lambda posts AI-generated summary to Slack channel with: severity, likely cause, suggested remediation steps
- [ ] **P6-4** Implement automated rollback trigger: if error budget burn rate critical → CodeDeploy rollback via Lambda
- [ ] **P6-5** Write a simulated Post-Mortem document (`postmortem-template.md`) using Blameless format
- [ ] **P6-6** Track incident metrics: MTTD (Mean Time to Detect), MTTR (Mean Time to Resolve) in `incident-log.md`

---

## Phase 7: Resource Optimization

> **SRE Skill:** Capacity planning, cost attribution, efficiency toil  
> **AWS Skill:** Compute Optimizer, Cost Explorer, ECS right-sizing, Savings Plans

- [ ] **P7-1** Enable AWS Compute Optimizer; review ECS Fargate task right-sizing recommendations
- [ ] **P7-2** Implement ECS Service Auto Scaling (Target Tracking: CPU 60%, Request Count per target)
- [ ] **P7-3** Configure ECS Fargate Spot for non-critical tasks (reduce cost by ~70%)
- [ ] **P7-4** Implement S3 lifecycle policies for log/artifact storage (transition to Glacier after 30 days)
- [ ] **P7-5** Enable RDS storage autoscaling; evaluate Aurora Serverless v2 for dev environment
- [ ] **P7-6** Create monthly cost report in `cost-optimization.md`: baseline → optimized (show % saved)
- [ ] **P7-7** Set up CloudWatch cost anomaly detection alert

---

## Phase 8: Security Hardening

> **SRE Skill:** Security as reliability, supply chain security  
> **AWS Skill:** GuardDuty, Security Hub, WAF, ECR image scanning, KMS

- [ ] **P8-1** Enable Amazon GuardDuty on the AWS account
- [ ] **P8-2** Enable AWS Security Hub with CIS AWS Foundations Benchmark checks
- [ ] **P8-3** Add ECR image vulnerability scanning (on push + weekly scheduled scan)
- [ ] **P8-4** Implement AWS WAF on ALB: rate limiting, SQL injection, XSS rules
- [ ] **P8-5** Enable CloudTrail with S3 + integrity validation for all API audit logs
- [ ] **P8-6** Implement KMS encryption for: RDS at rest, S3 artifacts, Secrets Manager
- [ ] **P8-7** Run `tfsec` or `checkov` on all Terraform code; fix HIGH/CRITICAL findings
- [ ] **P8-8** Implement VPC Flow Logs → CloudWatch; create alert for unusual outbound traffic
- [ ] **P8-9** Document security posture in `security-baseline.md` (findings before vs after)

---

## Phase 9: Interview Demonstration Prep

> **Goal:** Package the project as a demonstrable interview portfolio

- [ ] **P9-1** Write `README.md` with architecture diagram (draw.io or Mermaid), tech stack, and key decisions
- [ ] **P9-2** Record a 5-min Loom walkthrough: architecture → CI/CD → SLO dashboard → chaos experiment
- [ ] **P9-3** Prepare SRE interview story using STAR format for each phase (incident response, chaos, SLOs)
- [ ] **P9-4** Prepare "what would you do differently at scale?" talking points for each phase
- [ ] **P9-5** Push all code + runbooks to public GitHub repo

---

## 📊 Skill Coverage Matrix

| Phase | AWS Services Practiced                   | SRE Competency          |
| ----- | ---------------------------------------- | ----------------------- |
| P0    | IAM, VPC, Organizations                  | System design, tenets   |
| P1    | ECR, Fargate                             | Production readiness    |
| P2    | Terraform/CDK, RDS, ALB, Secrets Manager | IaC, toil reduction     |
| P3    | CodePipeline, CodeBuild, CodeDeploy      | DORA metrics, CI/CD     |
| P4    | CloudWatch, X-Ray, Synthetics            | SLOs, error budgets     |
| P5    | FIS (Fault Injection Simulator)          | Chaos engineering       |
| P6    | Lambda, Bedrock, EventBridge             | Incident management     |
| P7    | Compute Optimizer, Auto Scaling          | Capacity planning       |
| P8    | GuardDuty, WAF, Security Hub, KMS        | Security as reliability |
| P9    | —                                        | Interview storytelling  |

---

## 📅 Suggested Weekly Schedule (2–3 hrs/week)

| Week  | Phase   | Milestone                                       |
| ----- | ------- | ----------------------------------------------- |
| 1–2   | P0 + P1 | AWS account live, Focalboard containerized      |
| 3–4   | P2      | Focalboard running on ECS Fargate via Terraform |
| 5–6   | P3      | Full CI/CD pipeline with Blue/Green deploy      |
| 7–8   | P4      | SLO dashboard + error budget alerts live        |
| 9–10  | P5      | 4 chaos experiments completed                   |
| 11–12 | P6      | AI incident bot posting to Slack                |
| 13–14 | P7 + P8 | Cost optimized, security hardened               |
| 15–16 | P9      | Interview-ready portfolio                       |

---

*Last updated: 2026-02-27 | Project: focalboard-aws-platform*
