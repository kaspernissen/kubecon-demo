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

1. Boot the cluster: 

```sh
make cluster
```

2. Create the Dash0 secrets (optional - remove dash0 in the collector configurations)

```sh
kubectl create namespace opentelemetry
export DASH0_AUTH_TOKEN="..."
kubectl create secret generic dash0-secrets --from-literal=dash0-authorization-token=${DASH0_AUTH_TOKEN} --namespace opentelemetry
```

3. Install the prequisuites:

```sh
make install-all-prereqs
```

4. Let's have a look at the `Instrumentation` resource, and deploy it.

```sh
make instrumentation
```

5. let's have a look at the `OpenTelemetryCollector` resource, and deploy it.

```sh
make deploy-otel-collector
```

6. Let's build and deploy our services

```sh
make docker
make kind-load
make deploy-k8s
```

7. Let's deploy our dashboards

```sh
make deploy-perses-dashboards
```

