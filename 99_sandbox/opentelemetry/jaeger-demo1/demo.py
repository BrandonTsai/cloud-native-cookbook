import time
from opentelemetry import trace
from opentelemetry.exporter.jaeger.thrift import JaegerExporter
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor


# Create a tracer provider
tracer_provider = TracerProvider(
    resource=Resource.create({"service.name": "demo-app"}))

# Create a Jaeger exporter
jaeger_exporter = JaegerExporter(
    agent_host_name="localhost",
    agent_port=6831,
)

# Create a batch span processor and add the Jaeger exporter
span_processor = BatchSpanProcessor(jaeger_exporter)
tracer_provider.add_span_processor(span_processor)

# Set the tracer provider as the global tracer provider
trace.set_tracer_provider(tracer_provider)

# Get a tracer
tracer = trace.get_tracer(__name__)

# Define a function to create and record a span
def perform_operation():
    with tracer.start_as_current_span("demo-operation"):
        print("Performing some operation...")
        time.sleep(1)


# Call our function
perform_operation()

# Gracefully shutdown the exporter and span processor
span_processor.shutdown()
