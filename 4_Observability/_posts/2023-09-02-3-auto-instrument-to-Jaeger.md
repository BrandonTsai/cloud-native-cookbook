---
title: "#6-2 Auto Instrumentation for Exporting Python Traces to Jaeger"
author: Brandon Tsai
---

https://www.digitalocean.com/community/tutorials/how-to-implement-distributed-tracing-with-jaeger-on-kubernetes


Preparing the Demo Environment 
--------------------------------

Let's create a another project directory 

1. Create a Directory: Open your terminal and create a new directory for your project.

```bash
mkdir -p opentelemetry/jaeger-auto
cd opentelemetry/jaeger-auto
```

2. Set Up a Virtual Environment: Create a virtual environment to isolate your project's dependencies.

```bash
python3 -m venv venv
source venv/bin/activate  # On Windows, use 'venv\Scripts\activate'
```

3. Install Required Packages: Install the necessary packages for our demo.

```bash
python3 -m pip install --upgrade pip
python3 -m pip install Flask requests
```

Build Sample Applications
--------------------


frontend Web app: `/whoami?uid=0` --> backend API app: `/get_users` 

Implement the frontend application

Implement the backend application


Deploy and Testing the Sample Applications
----------------------------------------------------------------



Install the OpenTelemetry Python Agent
------------------------------

```bash
pip install opentelemetry-distro opentelemetry-exporter-otlp
opentelemetry-bootstrap -a install
```

The opentelemetry-distro package installs the API, SDK, and the opentelemetry-bootstrap and opentelemetry-instrument tools.

The opentelemetry-bootstrap -a install command reads through the list of packages installed in your active site-packages folder, and installs the corresponding instrumentation libraries for these packages, if applicable. For example, if you already installed the flask package, running opentelemetry-bootstrap -a install will install opentelemetry-instrumentation-flask for you.


Run the OpenTelemetry Python Agent 
----------------------------------------------------------------


```bash
OTEL_SERVICE_NAME=whoami \
OTEL_TRACES_EXPORTER=jaeger \
OTEL_METRICS_EXPORTER=none \
OTEL_EXPORTER_JAEGER_GRPC_INSECURE=true \
opentelemetry-instrument python frontend/frontend.py
```

```bash
OTEL_SERVICE_NAME=whoami \
OTEL_TRACES_EXPORTER=jaeger \
OTEL_METRICS_EXPORTER=none \
OTEL_EXPORTER_JAEGER_GRPC_INSECURE=true \
opentelemetry-instrument python backend/backend.py
```


Run the client.py


Investigating Traces in Jaeger
--------------------------------

Can see the span, but it said "No service dependencies found." in System Architecture


References:
----------------

- https://opentelemetry-python.readthedocs.io/en/latest/sdk/environment_variables.html
- https://www.digitalocean.com/community/tutorials/how-to-implement-distributed-tracing-with-jaeger-on-kubernetes

