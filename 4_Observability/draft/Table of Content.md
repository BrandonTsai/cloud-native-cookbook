## Chapter 1: Introduction

- What is Observability?
- What is Distributed Trace?
- What is Instrumentation?
- CNCF Observability Landscape
- What is OpenTelemetry(Otel)?

## Chapter 2 : OpenTelemetry code instrumentation & SDKs 

- Introduction to OpenTelemetry code instrumentation.
- Prerequisites: install python and flask as examples
- Manually Instrumentation and exporter trace to Jarger
- Auto Instrumentation and exporter trace and metric to Grafana Cloud ( Prometheus + Loki + Tempo ) Introduction

## Chapter 3 : Otel Collectoers

- [] Introduction
- [] Install & Configure (exporter to Jaeger and Grafana)
- [] Update previous python SDK example exporter to Otel collector
- [] encryption or even sensitive data filtering

## Chapter 5: OpenTelemetry and Containers
- [] Monitor Local Docker Containers on VM with OpenTelemetry

## Chapter 6: OpenTelemetry and Kubernetes
- [] Install Kubernetes: by Kind, K3S?
- [] OpenTelemetry Operator
- [] OpenTelemetry Collector to Grafana Cloud

## Chapter 7: OpenTelemetry and Serverless

- [] 

Advanced Configurations
================================================================

## Chapter 8: Context Propagation and Correlation

Importance of Context in Distributed Systems
OpenTelemetry's Approach to Context Propagation
Enabling Correlation Across Services for Better Insights
Real-world Use Cases and Best Practices
VM -> K8S Operator
Cross-Correlation: Connecting Distributed Traces with Logs

## Chapter 9: Sampling and Performance Considerations

Understanding Sampling in Distributed Tracing
Different Sampling Strategies and Their Trade-offs
Balancing Observability and Performance
Analyzing Trace Data with Smart Sampling


## Chapter 10: Splunk

- SDK to splunk
- Container to Splunk
- k8s to splunk
  - Collect Prometheus metrics with Splunk-Otel
  - Collect etcd metrics from OCP node exporter to Splunk
  - Collect Journald Logs from OCP node exporter to Splunk
  - From Splunk Enterprise to Splunk Cloud


## Chapter 11: Advanced Topics in OpenTelemetry

Using OpenTelemetry for Error Monitoring and Anomaly Detection
Combining Metrics, Traces, and Logs for Rich Insights
Cross-Correlation: Connecting Distributed Traces with Logs
Best Practices for Configuring and Managing OpenTelemetry

## Chapter 12: Conclution

- How to balance the cost of Observability?
- Future of OpenTelemetry and Observability
  - Trends and Innovations in Observability Space
  - Ongoing Developments and Community Contributions to OpenTelemetry
  - The Role of OpenTelemetry in Shaping Future Observability Standards
