# Log Output

**1. Install prerequisites**

- Docker
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [k3d](https://github.com/k3d-io/k3d#get)

Verify the installation:

```bash
docker --version
kubectl version --client
k3d version
```

**2. Create a Kubernetes cluster**

```bash
k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2
```

Verify that the cluster is running:

```bash
kubectl cluster-info
kubectl get nodes
```

> **Note:** `k3d cluster create` automatically configures `kubectl` to use the new cluster, so `kubectl config use-context k3d-k3s-default` is usually unnecessary unless you've switched to another context.

**3. Apply manifests**

```bash
kubectl apply -f manifests/
```

Verify that the deployment and pod were created:

```bash
kubectl get deployments
kubectl get pods
```

**4. Verify the application output**

Follow the logs:

```bash
kubectl logs -f deployment/logoutput-dep
```

Alternatively, using the pod name:

```bash
kubectl logs -f <pod-name>
```

Should see a UUID string printed approximately every 5 seconds.

**5. Access the application in browser**  

Ensure Ingress is listening on port 80:
```bash
kubectl get svc,ing
```

Access the application on `http://localhost:8081`.
