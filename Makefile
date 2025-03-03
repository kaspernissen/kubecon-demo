VERSION ?= v1
CLUSTER_NAME ?=otel

cluster:
	kind create cluster --name=$(CLUSTER_NAME) --config ./kind/multi-node.yaml

docker:
	docker build -f ./services/todo-go/Dockerfile -t todo-go:$(VERSION) ./services/todo-go
	docker build -f ./services/todo-java/Dockerfile -t todo-java:$(VERSION) ./services/todo-java

kind-load: docker
	kind load docker-image --name $(CLUSTER_NAME) todo-go:$(VERSION)
	kind load docker-image --name $(CLUSTER_NAME) todo-java:$(VERSION)

deploy-k8s: 
	kubectl apply -f ./services/todo-go/manifests/
	kubectl apply -f ./services/todo-java/manifests/

deploy-otel-operator-prereqs:
	helm repo add jetstack https://charts.jetstack.io --force-update
	helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
	helm upgrade --install cert-manager jetstack/cert-manager --namespace cert-manager --create-namespace --set crds.enabled=true
	helm upgrade --install opentelemetry-operator open-telemetry/opentelemetry-operator --set manager.extraArgs="{--enable-go-instrumentation}" --set "manager.collectorImage.repository=otel/opentelemetry-collector-k8s" --namespace opentelemetry --create-namespace
	kubectl apply --server-side -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/v0.80.1/example/prometheus-operator-crd/monitoring.coreos.com_servicemonitors.yaml

deploy-otel-collector: deploy-otel-operator-prereqs 
	kubectl apply -f ./otel-collector/rbac-daemonset.yaml
	kubectl apply -f ./otel-collector/rbac-central.yaml
	kubectl apply -f ./otel-collector/rbac-target-allocator.yaml
	kubectl apply -f ./otel-collector/otel-collector-daemonset.yaml
	kubectl apply -f ./otel-collector/otel-collector-central.yaml

deploy-prometheus-jaeger:
	helm upgrade --install demo open-telemetry/opentelemetry-demo -f ./prometheus-jaeger/values.yaml

deploy-perses-operator:
	@echo "Installing Perses Operator CRDs and deploying the operator..."
	kustomize build github.com/perses/perses-operator/config/default | kubectl apply -f -
	kubectl apply -f ./perses/perses.yaml

deploy-perses-dashboards:
	kubectl apply -f ./perses/prometheus-datasource.yaml
	kubectl apply -f ./perses/go-dashboard.yaml

deploy-postgres:
	helm install pg oci://registry-1.docker.io/bitnamicharts/postgresql \
  		--set global.postgresql.auth.postgresPassword=password \
  		--set global.postgresql.auth.database=todo

instrumentation:
	kubectl apply -f ./instrumentations/instrumentation.yaml

