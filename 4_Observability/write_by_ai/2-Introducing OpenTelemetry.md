https://www.google.com/url?sa=i&url=https%3A%2F%2Fcribl.io%2Fblog%2Fobservability-vs-monitoring-vs-telemetry%2F&psig=AOvVaw2X_BzYHx7yVp8gReeHzHc_&ust=1693719667298000&source=images&cd=vfe&opi=89978449&ved=2ahUKEwi14ZvKm4uBAxVNbmwGHT3EBWIQjhx6BAgAEA0

Introducing OpenTelemetry
=========================

As we embark on our journey to delve deeper into observability and its practical applications, this chapter introduces us to OpenTelemetry, a transformative project that offers a standardized approach to capturing observability data. We'll explore the evolution and significance of OpenTelemetry, its role as a standardized observability framework, and the key components that make it a cornerstone of modern software observability.


Evolution and Background of OpenTelemetry
----------------------------------------------------------------

OpenTelemetry didn't emerge out of thin air; it's the result of a collaborative effort within the observability community. Born from the merger of two influential projects – OpenTracing and OpenCensus – OpenTelemetry unifies their strengths to create a comprehensive observability framework.

OpenTracing focused on providing instrumentation libraries to capture traces in various programming languages, enabling developers to gain insights into distributed systems. OpenCensus, on the other hand, aimed to collect metrics and distributed traces in a vendor-agnostic manner. The convergence of these projects resulted in OpenTelemetry, a project backed by the Cloud Native Computing Foundation (CNCF) and embraced by a thriving community.


OpenTelemetry as a Standardized Observability Framework
----------------------------------------------------------------

At its core, OpenTelemetry seeks to establish a standardized way of capturing observability data across different languages and frameworks. It's more than just a collection of libraries; it's a cohesive framework that provides guidance, APIs, and instrumentation tools to help developers seamlessly embed observability into their applications.

OpenTelemetry simplifies the process of instrumenting code to capture traces and metrics. By following its conventions, developers ensure consistent data collection regardless of the language or technology stack they're using. This standardization enables interoperability and consistency, making it easier to integrate observability data across the entire software ecosystem.


Signals:
----------------

**Logs**: 

**Metrics**: 

**Traces**: 

**Baggage**:

Components and Architecture of OpenTelemetry
----------------------------------------------------------------

OpenTelemetry comprises several key components that work together to capture and propagate observability data:

**SDKs and Instrumentation Libraries:** OpenTelemetry provides Software Development Kits (SDKs) and instrumentation libraries for various programming languages. These SDKs enable developers to instrument their code to capture traces, metrics, and contextual information.

**APIs and Conventions:** OpenTelemetry defines APIs and conventions for capturing and exporting observability data. These APIs ensure a consistent way of working with traces, metrics, and context propagation.

**Trace Exporters and Collectors**: OpenTelemetry offers trace exporters that allow captured trace data to be sent to various backends for storage and analysis. Collectors facilitate the aggregation and transformation of trace data before exporting.

**Metrics Exporters:** Similar to trace exporters, metrics exporters enable the transport of metrics data to backend systems for storage and visualization.

Context Propagation: OpenTelemetry ensures contextual information, such as trace and span IDs, is propagated across services to maintain trace continuity and correlation.

The architecture of OpenTelemetry is designed to be flexible and adaptable, allowing it to seamlessly integrate with different technology stacks and observability backends.