---
title: "#2 Understanding the Power of Distributed Tracing"
author: Brandon Tsai
---

In order to make a system observable, the application must emit signals such as `traces`, `metrics`, and `logs` to backend observability platform.
And a trace, records the paths taken by requests as they propagate through multi-service architectures, like microservice and serverless applications.
In this article, we will explore the base of distributed tracing, its benefits, and how it plays a crucial role in maintaining the health and performance of software systems.


What is Distributed Tracing?
----------------------------------------------------------------

In the world of modern software development, applications have become increasingly complex. They traverse multiple services, databases, and cloud providers, making it challenging to monitor, diagnose, and troubleshoot issues efficiently. This is where distributed tracing comes to the rescue. Distributed tracing is a powerful technique that allows developers and operations teams to monitor and gain insights of the flow of requests as they traverse through a distributed system. It provides a detailed view of the interactions between various components of an application, such as microservices, databases, and external services. This information is captured in the form of **traces**, which are sequences of timed events representing the path of a request.


Traces and Spans
----------------------------------------------------------------

In the realm of distributed tracing, a trace is a view into the journey of a request as it moves through a a complex, distributed system.
A trace is made of one or more "spans". A span is a unit of work that represents an individual operations of the workflow.
The first span represents the root span. Each root span represents a request from start to finish. The spans underneath the parent provide a more in-depth context of what occurs during a request (or what steps make up a request).


Span Composition
----------------------------------------------------------------

A span contains name, time-related data, structured log messages, and other metadata to provide information about the specific activity, such as a function call or an HTTP request, along with its duration and associated metadata.

A standard span in distributed tracing includes:

**An operation/service name:** 
A title of the work performed

**Timestamps:**
A reference from the start to the end of the system process

**Span Tags (Attributes)**
Essentially, span tags allow users to define customized annotations that facilitate querying, filtering, and other functions involving trace data. Examples of span tags include db.instances that identify a data host, serverID, userID, and HTTP response code.

Developers may apply standard tags across common scenarios, including db.type (string tag), which refers to database type and peer.service (integer tag) that references a remote port. Key:value pairs provide spans with additional contexts, such as the specific operation it tracks.

Tags provide developers with the specific information necessary for monitoring multi-dimensional queries that analyze a trace. For instance, with span tags, developers can quickly home in on the digital users facing errors or determine the API endpoints with the slowest performance.

Developers should consider maintaining a simple naming convention for span tags to fulfill operations with ease and minimal confusion.

A span may also have zero or more key/value attributes. Attributes allow you to create metadata about the span. For example, you might create attributes that hold a customer ID, or information about the environment that the request is operating in, or an app’s release. Attributes don’t reflect any time-based event (events in OpenTelemetry). The OpenTelemetry spec defines several standard attributes. You can also implement your own attributes.

**Span Logs**
Key:value span logs enable users to capture span-specific messages and other data input from an application. Users refer to span logs to document exact events and timelines in a trace. While tags apply to the whole span, logs refer to a “snapshot” of the trace.

**Span Context**
includes IDs that identify and monitor spans across multiple process boundaries and baggage items such as key:value pairs that cross process boundaries 

The Span Context carries data across various points/boundaries in a process. Logically, a SpanContext divides into two major components: user-level baggage and implementation-specific fields that provide context for the associated span instance.

Essentially, baggage items are key:value pairs that cross process boundaries across distributed systems. Each instance of a baggage item contains valuable data that users may access throughout a trace. Developers can conveniently refer to the SpanContext for contextual metrics (e.g. service requests and duration) to facilitate troubleshooting and debugging processes.



Benefits of Distributed Tracing
--------------------------------

1. Root Cause Analysis:
Distributed tracing simplifies the process of identifying the root cause of performance issues or errors in a distributed system. By examining the traces, developers and operators can pinpoint the exact service or component responsible for a problem, reducing the mean time to resolution (MTTR).

1. Performance Optimization:
Traces provide valuable insights into the performance of an application. By analyzing the duration of spans and understanding where time is spent, teams can optimize bottlenecks and improve overall system performance.

1. Dependency Mapping:
Distributed tracing also aids in creating dependency maps, showing how different services interact with each other. This information is crucial for understanding the architecture of complex applications and ensuring that changes do not introduce unexpected issues.

1. Proactive Monitoring:
Rather than waiting for issues to arise, distributed tracing enables proactive monitoring. By setting thresholds and alerts based on trace data, teams can detect and address potential problems before they impact users.

Tools for Distributed Tracing
----------------------------------------------------------------

Several tools and frameworks are available for implementing distributed tracing in your applications. Some popular options include:

**OpenTelemetry**: An open-source project that provides a set of APIs, libraries, agents, and instrumentation to enable distributed tracing and observability in various programming languages.

**Jaeger**: An open-source, end-to-end distributed tracing system that is particularly well-suited for microservices architectures.

**Zipkin**: Another open-source distributed tracing system that helps track requests as they flow through various services.

**Commercial Solutions**: Many cloud providers and observability platforms offer distributed tracing as part of their monitoring and debugging tools.


Conclusion
---------
In the ever-evolving landscape of software development, understanding and optimizing the performance of distributed applications is paramount. Distributed tracing has emerged as a vital tool for achieving this goal. By providing detailed insights into how requests flow through complex systems, distributed tracing enables teams to troubleshoot issues, optimize performance, and ensure the reliability of their applications. As software systems continue to grow in complexity, distributed tracing will remain a critical component of modern application monitoring and observability strategies.