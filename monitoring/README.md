### 1. Install prerequisites

- [Helm CLI](https://helm.sh/docs/intro/install/) - the package manager for Kubernetes

Verify the installation:

```bash
helm version
```

### 2. Add the Prometheus and Grafana repositories

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
```

### 3. Create namespace for monitoring

Check if namespace `monitoring` doesn't exist and create it:

```bash
kubectl get namespaces | grep monitoring
kubectl create namespace monitoring
```

### 4. Install the charts

- Install [Prometheus](https://prometheus.io/):
```bash
helm upgrade --install prom prometheus-community/prometheus \
  --namespace monitoring \
  --create-namespace \
  --values manifests/prom-values.yaml
```

- Install [Loki](https://grafana.com/oss/loki/):
```bash
helm upgrade --install loki grafana/loki \
  --namespace monitoring \
  --values manifests/loki-values.yaml
```

- Install [Alloy](https://grafana.com/oss/alloy-opentelemetry-collector/):
```bash
helm upgrade --install k8smon grafana/k8s-monitoring \
  --namespace monitoring \
  --values manifests/k8smon-values.yaml
```

- Install [Grafana](https://grafana.com/oss/grafana/):
```bash
helm upgrade --install grafana grafana/grafana \
  --namespace monitoring \
  --values manifests/grafana-values.yaml
```

Ensure all pods are running:
```bash
helm list --namespace monitoring
kubectl get svc --namespace monitoring
kubectl get pods --namespace monitoring
```

###  5. Access Grafana

- Set up port forward for Grafana service:
```bash
kubectl port-forward --namespace monitoring svc/grafana 3000:80
```
- Open Grafana on `http://localhost:3000`.
-
