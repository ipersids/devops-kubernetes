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
k3d cluster create -a 2
```

Verify that the cluster is running:

```bash
kubectl cluster-info
kubectl get nodes
```

> **Note:** `k3d cluster create` automatically configures `kubectl` to use the new cluster, so `kubectl config use-context k3d-k3s-default` is usually unnecessary unless you've switched to another context.

**3. Deploy the application**

```bash
kubectl create deployment hashgenerator-dep --image=ipersids/logoutput:ex1.1
```

Verify that the deployment and pod were created:

```bash
kubectl get deployments
kubectl get pods
```

Wait until the pod status is **Running**.

**4. Verify the application output**

Follow the logs:

```bash
kubectl logs -f deployment/hashgenerator-dep
```

Alternatively, using the pod name:

```bash
kubectl logs -f <pod-name>
```

Should see a UUID string printed approximately every 5 seconds.
