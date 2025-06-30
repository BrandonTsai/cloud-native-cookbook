2025 Q1
===========


OTel

- [] Observability with OpenTelemetry.
- [] Observability without Compromise.

Splunk-Otel-Collector

- [] Collect Prometheus metrics with Splunk-Otel
- [] Collect etcd metrics from OCP node exporter to Splunk
- [] Collect Journald Logs from OCP node exporter to Splunk
- [] **From Splunk Enterprise to Splunk Cloud**
- [] https://splunk.github.io/observability-workshop/v5.94/en/index.html


eBPF:

- [] eBPF Network Collector - https://github.com/open-telemetry/opentelemetry-network




Observability with OpenTelemetry
===================================

## Chapter 1: Introduction

- What is Observability?
- What is Distributed Trace?
- What is Instrumentation?
- CNCF Observability Landscape
- What is OpenTelemetry(Otel)?

## Chapter 2 : OpenTelemetry code instrumentation & SDKs 

- Introduction to OpenTelemetry code instrumentation.
- Prerequisites: install python and flask as examples
- Manually Instrumentation: exporter python trace to Jarger
- Auto Instrumentation: exporter python trace to Jarger


## Chapter 3 : Backend Provider - Grafana Lab/Stack

- Grafana Cloud as Example
- https://grafana.com/oss/


## Chapter 4 : Otel Collectoers

- [] Introduction
- [] Update previous python SDK example exporter to Otel collector
- [] encryption or even sensitive data filtering

## Chapter 5: OpenTelemetry on Hosted VMs
- [] Monitor Local Docker Containers on VM with OpenTelemetry

## Chapter 6: OpenTelemetry on Kubernetes
- [] Install Kubernetes: by Kind, K3S?
- [] OpenTelemetry Operator
- [] OpenTelemetry Collector to Grafana Cloud

## Chapter 7: OpenTelemetry and ...

- [] Serverless (AWS Lambda)
- [] Database

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


## Chapter 10: Advanced Topics in OpenTelemetry

Using OpenTelemetry for Error Monitoring and Anomaly Detection
Combining Metrics, Traces, and Logs for Rich Insights
Cross-Correlation: Connecting Distributed Traces with Logs
Best Practices for Configuring and Managing OpenTelemetry


## Chapter 11: Explore Other Backend Provider - Splunk

- Splunk-Otel-Collector
  - Collect Prometheus metrics with Splunk-Otel
  - Collect etcd metrics from OCP node exporter to Splunk
  - Collect Journald Logs from OCP node exporter to Splunk
  - From Splunk Enterprise to Splunk Cloud


## Chapter 12: eBPF

- eBPF Network Collector
  https://github.com/open-telemetry/opentelemetry-network
- eBPF based automatic instrumentation for HTTP/gRPC App
  Grafana Beyla: https://grafana.com/oss/beyla-ebpf/
- Cilium 
https://cilium.io/use-cases/metrics-export/



## Chapter 12: Conclution

- How to balance the cost of Observability?
- Future of OpenTelemetry and Observability
  - Trends and Innovations in Observability Space
  - Ongoing Developments and Community Contributions to OpenTelemetry
  - The Role of OpenTelemetry in Shaping Future Observability Standards



Observability without Compromise.
==============================

Comprehensive Data Coverage
------------------------------

**No Sampling:**
Traditional monitoring often relies on sampling data, which can lead to missed issues or inaccurate insights. Observability without compromise aims for complete data ingestion, ensuring all relevant information is captured. 

**All Pillars of Observability:**
This means collecting and analyzing logs, metrics, and traces, which are the three key components of observability. 

**Extended Retention:**
Instead of limiting data retention periods, this approach advocates for keeping data long enough to analyze historical trends and troubleshoot past issues effectively. 


Cost-Effective Implementation
------------------------------

**Data Optimization:**
Organizations can reduce observability costs by optimizing the data they collect, focusing on relevant information and reducing the amount of irrelevant data being shipped. 

**Data Tiering:**
Implementing tiered storage solutions can help manage costs by storing less frequently accessed data in cheaper storage options while maintaining high-performance access to frequently used data. 

**Efficient Tooling:**
Leveraging tools and platforms designed for efficient data processing and analysis can minimize costs associated with data storage and querying. 


Performance and Agility:
--------------------------


**Fast Querying and Analysis:**
Observability without compromise necessitates fast and efficient querying of observability data, enabling quick responses to incidents and faster troubleshooting. 
**No-Index Architecture:**
Some approaches, like those using a no-index architecture, prioritize speed and performance by avoiding the overhead of indexing large datasets. 
**Automated Insights:**
Tools and platforms that automate data analysis, such as those using AI and machine learning, can accelerate the identification of issues and their root causes. 


Ease of Use and Collaboration:
----------------------------------

**Unified Data Platform:**
A unified observability platform that provides a single source of truth for all relevant data can streamline troubleshooting and collaboration across teams. 
**Minimal Code Changes:**
The ability to gain observability without requiring extensive code changes or redeployment can reduce the risk and complexity of implementing observability solutions. 
**Enhanced Collaboration:**
By providing a common understanding of system behavior, observability without compromise can improve collaboration between development, operations, and other teams. 
In essence, observability without compromise is about having the ability to deeply understand your systems, identify issues quickly, and resolve them effectively, all without making sacrifices on the completeness, cost, or performance of your observability implementation. 


Observability with eBPF
===================================

The same use case for eBPF makes projects like Falco (security), Pixie (APM for apps on Kubernetes), and Cilium (networking monitoring) possible.




