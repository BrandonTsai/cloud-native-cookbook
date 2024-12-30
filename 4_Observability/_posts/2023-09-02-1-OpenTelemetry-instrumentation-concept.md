---
title: "#6 OpenTelemetry Code Instrumentation Introduction"
author: Brandon Tsai
---

As we mentioned earlier, code instrumentation is the keystone of observability. It empowers you to unveil the hidden intricacies of your software system by incorporating monitoring and tracing functionalities into your codebase.

OpenTelemetry stands as a leading force in the world of observability, providing a standardized and extensible framework for instrumenting your applications. It is supported for many popular programming languages and frameworks. Whether you're building microservices, monolithic applications, or anything in between, OpenTelemetry equips you with the tools needed to capture crucial data about your software's performance, dependencies, and behavior.



Automatic & Manual instrumentation
---------------------

In OpenTelemetry, there are two methods to do code instrumentation to help you capture observability data like traces and metrics.

### Automatic Instrumentation:

**Pros:**

- Ease of Use: Automatic instrumentation involves using OpenTelemetry's pre-built instrumentation libraries and integrations. These libraries are designed to automatically capture data without requiring you to make extensive code changes.
- Consistency: It ensures uniform observability practices across different parts of your application or services because the instrumentation process is standardized.

**Cons:**

- Limited Customization: Automatic instrumentation might not allow for fine-grained customization. You might have less control over which specific parts of your code are instrumented and how data is captured.
- Language and Framework Support: The availability of automatic instrumentation libraries can vary depending on the programming language and framework you are using. Some languages and frameworks may have more extensive support than others.



### Manual Instrumentation:

**Pros:**
- Full Control: Manual instrumentation provides you with complete control over what gets instrumented and how. You can tailor instrumentation to specific code sections, allowing for fine-grained observability.
- Customization: It allows you to add additional context, tags, or metadata to traces and metrics, which can be valuable for more detailed analysis.

**Cons:**
- Complexity: Manual instrumentation can be more complex and time-consuming. You need to write the code to capture observability data explicitly.
- Consistency Challenges: Ensuring consistent instrumentation practices across your codebase might require extra diligence because it's up to the developer to apply instrumentation consistently.


In practice, you might use a combination of both automatic and manual instrumentation based on your needs. Automatic instrumentation can be a quick way to get basic observability data from your applications, while manual instrumentation allows you to fine-tune and customize data capture for specific, critical sections of your code. The choice between them depends on your observability requirements and the level of control and customization you need in your application.



Signal supported
-----------------

OpenTelemetry code instrumentation is supported for many popular programming languages
However, not all languages support all traces, metrics and logs signals. Please refer the official [website](https://opentelemetry.io/docs/instrumentation/) to get the latest status of the major functional components for OpenTelemetry.


| Language      | Traces | Metrics             | Logs                | Automatic Instrumentation |
| ------------- | ------ | ------------------- | ------------------- | ------------------------- |
| C++           | Stable | Stable              | Experimental        | N/A                       |
| C#/.NET       | Stable | Stable              | Mixed*              | Supported                 |
| Erlang/Elixir | Stable | Experimental        | Experimental        | N/A                       |
| Go            | Stable | Mixed*              | Not yet implemented | N/A                       |
| Java          | Stable | Stable              | Stable              | Supported                 |
| JavaScript    | Stable | Stable              | Development         | Supported                 |
| PHP           | Beta   | Beta                | Alpha               | N/A                       |
| Python        | Stable | Stable              | Experimental        | Supported                 |
| Ruby          | Stable | Not yet implemented | Not yet implemented | Supported                 |
| Rust          | Beta   | Alpha               | Alpha               | N/A                       |
| Swift         | Stable | Experimental        | In development      | N/A                       |