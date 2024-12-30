---
title: "#4 Navigating the CNCF Observability Landscape"
author: Brandon Tsai
---

What is CNCF?
----------------------------------------------------------------

CNCF stands for the "Cloud Native Computing Foundation", which is a well-known organization in the world of cloud computing and software development. CNCF serves as a home for open-source projects and initiatives that are focused on advancing the adoption and innovation of cloud-native technologies.

Key aspects of CNCF include:

- **Open Source**: CNCF hosts a variety of open-source projects related to cloud-native computing. These projects are developed collaboratively by communities of developers and are freely available for anyone to use, modify, and contribute to.

- **Cloud-Native Technologies**: CNCF's primary focus is on cloud-native technologies. These technologies are designed to leverage the scalability, agility, and flexibility of cloud computing environments. Cloud-native applications are typically containerized, dynamically orchestrated, and designed for microservices architecture.

- **Ecosystem Growth**: CNCF aims to grow the cloud-native ecosystem by supporting projects, promoting standards, and fostering collaboration within the community. This includes providing resources and support for both emerging and mature projects.

- **Interoperability and Compatibility**: CNCF encourages interoperability and compatibility among cloud-native tools and platforms. This helps organizations avoid vendor lock-in and ensures that cloud-native solutions work well together.

- **Education and Outreach**: CNCF provides educational resources, training, and events to help individuals and organizations learn about and adopt cloud-native technologies. These resources include webinars, conferences, and certification programs.

Overall, CNCF plays a vital role in shaping the future of cloud-native computing by fostering innovation, collaboration, and the development of open-source technologies that enable organizations to build and operate modern, scalable, and resilient applications in cloud environments.


CNCF Observability Projects
---------------------

As applications become more distributed and dynamic, traditional observability approaches often fall short. This is where cloud-native observability platforms come into play, redefining how we gain insights into our systems

The Cloud Native Computing Foundation (CNCF) hosts a variety of projects that form the backbone of cloud-native observability. 
Let's introduce some key CNCF projects:

### CNCF Graduated Projects:

**Prometheus** an open-source monitoring and alerting toolkit for metrics. It's designed for reliability and scalability in dynamic, cloud-native environments. Prometheus excels at time-series data collection and querying. It discovers targets automatically and scrapes their metrics, making it well-suited for containerized applications.

**Fluentd**: an open-source data collector that unifies data collection and consumption for logs. It's cloud-native by design, supporting various input and output plugins, making it an excellent choice for aggregating logs from different services and shipping them to centralized storage.

**Jaeger**: an end-to-end distributed tracing system. It helps you monitor and troubleshoot transaction latency in complex, microservices-based architectures. Jaeger provides insights into how requests propagate through services, making it vital for understanding service dependencies and performance bottlenecks.

### CNCF Incubating Projects 

**OpenTelemetry**: an essential initiative to create and supervise telemetry data such as traces, metrics, and logs. It offers APIs, libraries, agents, and instrumentation to enable distributed tracing and observability across different languages and environments. OpenTelemetry provides a standardized way to instrument applications for tracing, ensuring consistency in observability practices.

### CNCF Member Projects

**Grafana**: often paired with Prometheus. It provides a powerful platform for visualizing and analyzing metrics. Grafana's rich set of plugins and dashboards makes it a popular choice for building customizable observability dashboards.

**Grafana Loki**: focuses on log aggregation and storage. It's optimized for cloud-native environments, using labels to filter and query logs efficiently. Loki pairs well with Grafana for visualization.

**Grafana Tempo**: a high-volume, minimal-dependency, and open-source distributed tracing backend. it serves as a straightforward, cost-effective, and scalable alternative to Jaeger. Unlike many other tracing backends, which often necessitate the use of data stores like Cassandra or ElasticSearch, Grafana Tempo operates efficiently with just object storage, such as Google Cloud Storage or Amazon S3. This unique approach empowers software teams to effortlessly accumulate and retain a greater volume of traces from their distributed applications, all without the requirement for sampling.


Conclusion
----------

Observability is a fundamental aspect of managing modern applications in dynamic and distributed environments. CNCF owned projects like Prometheus, Fluentd, Jaeger, and OpenTelemetry are pivotal in providing the tools and standards needed to navigate this complex landscape effectively. Many famouse observability platform service providers, such as Splunk, DataDog and Rew Relic, are also highly intergrated with these CNCF projects.
Embracing and leveraging these projects can help organizations gain deeper insights, improve reliability, and optimize their cloud-native applications for peak performance and resilience.