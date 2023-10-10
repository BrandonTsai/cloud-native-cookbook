---
title: "#5 OpenTelemetry Introduction"
author: Brandon Tsai
---

What is OpenTelemetry
=========================

OpenTelemetry (also referred to as OTel) is an open source observability framework made up of a collection of tools, APIs, and SDKs. 
It serves as a beacon of innovation and standardization in the observability landscape. 
OpenTelemetry is more than just a set of tools. It's a unified framework that equips developers with the means to capture and collect observability data across different programming languages and frameworks.

OpenTelemetry transcends the boundaries of a single technology stack. Whether you're building applications in Java, Python, Go, or any other language, whether you're leveraging Kubernetes, serverless architectures, or traditional VMs, OpenTelemetry provides the tools to bring observability into your software's DNA. By standardizing data formats and instrumentation practices, OpenTelemetry opens the doors to a new era of interoperability and integration in the observability ecosystem.


Evolution and Background of OpenTelemetry
----------------------------------------------------------------

OpenTelemetry didn't emerge out of thin air; it's the result of a collaborative effort within the observability community. 
Born from the collaboration of the Cloud Native Computing Foundation, the OpenTracing project, and the OpenCensus project, 

- `OpenTracing` focused on providing instrumentation libraries to capture traces in various programming languages, enabling developers to gain insights into distributed systems.
- `OpenCensus`, on the other hand, aimed to collect metrics and distributed traces in a vendor-agnostic manner. The convergence of these projects resulted in OpenTelemetry, a project backed by the Cloud Native Computing Foundation (CNCF) and embraced by a thriving community.

OpenTelemetry unifies their strengths to create a comprehensive observability framework.


Signals
----------------

OpenTelemetry is a robust and versatile observability framework designed to provide deep insights into the behavior and performance of modern software applications. Within OpenTelemetry, "signals" (also known as "Telemetry data") play a pivotal role. Signals represent distinct categories of telemetry data that are essential for understanding and monitoring your applications. Following are the currently supported signals in OpenTelemetry.

### Logs
Logs are a fundamental source of information for diagnosing issues, auditing application behavior, and maintaining a historical record of events. In OpenTelemetry, logging is enhanced with structured data, making it easier to filter and analyze logs efficiently. With OpenTelemetry's log signals, you can:

Capture structured log messages that provide context and detail about events.
Enhance debugging and troubleshooting by correlating logs with traces.
Integrate with various logging backends and tools, ensuring seamless log management.

### Metrics
Metrics are essential for quantifying the performance and health of your applications. They provide a quantitative view of various aspects, such as resource utilization, error rates, and business-specific KPIs. OpenTelemetry's metrics signals empower you to:

Collect and aggregate metrics data across your application's components.
Create custom metrics to track specific business metrics or application behavior.
Export metrics to various monitoring systems, enabling real-time visualization and alerting.

### Traces
Traces are a cornerstone of observability, offering a detailed view of how requests and transactions flow through a distributed system. OpenTelemetry's tracing signals provide the means to:

Trace the path of requests as they traverse various services and components.
Monitor and analyze the performance of individual requests and operations.
Correlate traces with logs and metrics to gain a comprehensive understanding of application behavior.

### Baggage
Baggage is a relatively unique concept in the world of observability. It enables the propagation of contextual information across services within a distributed system. With OpenTelemetry's baggage signals, you can:

Attach key-value pairs to requests, allowing you to carry contextual information throughout a request's journey.
Enhance trace context with metadata that provides valuable insights into request context.
Facilitate debugging and troubleshooting by preserving context as requests traverse multiple services.


Choose observability backend for OpenTelemetry
--------------------------------------

There are two major steps involved in setting up observability for your application:

- **Instrumentation**: Collecting relevant data such as trace, metrics and logs that indicates the application health and behaviors.
- **Observability backend**: Storing, managing, and visualizing the collected data to take quick actions

While OpenTelemetry pay attentions to addresses the first step, you still need to set up a tool like Jaeger, Grafa or Splunk for storing and visualizing telemetry data.


Components and Architecture of OpenTelemetry
----------------------------------------------------------------

OpenTelemetry comprises several key components that work together to capture and propagate observability data:

**SDKs and Instrumentation Libraries:** OpenTelemetry provides Software Development Kits (SDKs) and instrumentation libraries for various programming languages. These SDKs enable developers to instrument their code to capture traces, metrics, and contextual information.

**APIs and Conventions:** OpenTelemetry defines APIs and conventions for capturing and exporting observability data. These APIs ensure a consistent way of working with traces, metrics, and context propagation.

**Trace Exporters and Collectors**: OpenTelemetry offers trace exporters that allow captured trace data to be sent to various backends for storage and analysis. Collectors facilitate the aggregation and transformation of trace data before exporting.

**Metrics Exporters:** Similar to trace exporters, metrics exporters enable the transport of metrics data to backend systems for storage and visualization.

Context Propagation: OpenTelemetry ensures contextual information, such as trace and span IDs, is propagated across services to maintain trace continuity and correlation.

The architecture of OpenTelemetry is designed to be flexible and adaptable, allowing it to seamlessly integrate with different technology stacks and observability backends.



Conclusion: OpenTelemetry as a Standardized Observability Framework
----------------------------------------------------------------

At its core, OpenTelemetry seeks to establish a standardized way of capturing observability data across different languages and frameworks. It's more than just a collection of libraries; it's a cohesive framework that provides guidance, APIs, and instrumentation tools to help developers seamlessly embed observability into their applications.

OpenTelemetry simplifies the process of instrumenting code to capture traces and metrics. By following its conventions, developers ensure consistent data collection regardless of the language or technology stack they're using. This standardization enables interoperability and consistency, making it easier to integrate observability data across the entire software ecosystem.
