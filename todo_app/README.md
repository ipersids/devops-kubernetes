# TODO App

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

**3. Deploy the application**

```bash
kubectl create deployment todo-dep --image=ipersids/todoapp:ex1.2
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
kubectl logs -f deployment/todo-dep
```

Alternatively, using the pod name:

```bash
kubectl logs -f <pod-name>
```

By default server listens on PORT 8080, set different one using command:

```bash 
kubectl set env deployment/todo-dep PORT=5454
```

Should see a message `Server started in port XXXX`.
