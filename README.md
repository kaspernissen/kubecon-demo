# OpenTelemetry Operator, Collector, Auto-Instrumentation, Dashboards as Code with Perses

This repository contains the setup for my demo at KubeCon+CloudNativeCon Europe 2025.

The demo is comprised of multiple different open source projects, amongst others;

* OpenTelemetry for telemetry collection and auto-instrumentation
* Prometheus for metric storage
* Jaeger for tracing storage
* Perses for dashboards as code

Further this project makes use of `kind` to run a local cluster, and `helm` to deploy tools. 

The demo also includes to applications written in `golang` and `java` to showcase, the auto-instrumentatio features of OpenTelemetry Operator.

## Get started

Boot the cluster: 

```
make cluster
```
