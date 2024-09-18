
<img width="1512" alt="Screenshot 2024-09-19 at 00 46 37" src="https://github.com/user-attachments/assets/6b6226ce-529e-4e1a-8801-5a00f0fff070">

This project demonstrates a basic architecture for setting up distributed tracing with OpenTelemetry using a Golang application, OpenTelemetry Agent, Gateway, and Jaeger for visualization. The traces collected are exported to Elasticsearch for storage and further analysis.

Architecture Overview

**Components:**
Golang Application with OpenTelemetry Instrumentation:
A Golang application is instrumented with OpenTelemetry SDK to generate trace data.

**OpenTelemetry Agent (Running on the Host):**
The OpenTelemetry Agent is deployed as a local collector on the same host as the Golang application.
It receives telemetry data (traces) directly from the Golang application.
The agent can perform some processing, filtering, or sampling on the trace data.
The processed trace data is then forwarded to the OpenTelemetry Gateway.

**OpenTelemetry Gateway:**
The gateway acts as a centralized telemetry processing point.
It collects traces from the OpenTelemetry Agent.
The collected traces are processed and forwarded to Jaeger for visualization and Elasticsearch for long-term storage.

**Jaeger:**
Jaeger is the tool used for visualizing trace data.
The traces sent from the OpenTelemetry Gateway are displayed in Jaeger, allowing for detailed monitoring and tracing of distributed applications.

**Elasticsearch:**
Elasticsearch is used as the storage backend for the trace data.
Jaeger stores the traces in Elasticsearch, making them searchable and available for long-term analysis.

**Data Flow:**
The Golang application generates telemetry data (traces) via OpenTelemetry instrumentation.
The OpenTelemetry Agent running on the same host collects this telemetry and forwards it to the OpenTelemetry Gateway.
The Gateway aggregates telemetry data from multiple agents, processes it, and exports the traces to both Jaeger (for visualization) and Elasticsearch (for long-term storage).
