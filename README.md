# devops-kubernetes

Submissions for the DevOps with Kubernetes course at the University of Helsinki

## Setup

#### 1. Install prerequisites

- Docker
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [k3d](https://github.com/k3d-io/k3d#get)

Verify the installation:

```bash
docker --version
kubectl version --client
k3d version
```

#### 2. Create a Kubernetes cluster

```bash
k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2
```

The manifests include a local `PersistentVolume` that is pinned to the first k3d agent node:

```yaml
k3d-k3s-default-agent-0
```

If your cluster/node names are different, update `manifests/persistentvolume.yaml` before applying it.

#### 3. Apply manifests

```bash
kubectl apply -f manifests/
kubectl apply -f todo_app/manifests/persistentvolumeclaim.yaml
kubectl apply -f log_output/manifests
kubectl apply -f pingpong/manifests
kubectl apply -f todo_app/manifests
```

Check that the storage and pods are ready:

```bash
kubectl get pv,pvc
kubectl get pods
```

The PVC should be `Bound` before the pods that mount it can start.

#### 4. Access the applications

```bash
kubectl get svc,ing
```

Open the apps through Ingress:

- Todo app: <http://localhost:8081/>
- Log output: <http://localhost:8081/logoutput>
- Ping-pong: <http://localhost:8081/pingpong>

## Course tasks

| Exercise | Tag | Description |
| :------: | :-: | :---------- |
| [1.1](https://github.com/ipersids/devops-kubernetes/tree/main/log_output) | [1.1](https://github.com/ipersids/devops-kubernetes/tree/1.1) | Deploy an app (Log output) that generates a random string on startup, stores this string into memory, and outputs it every 5 seconds with a timestamp into a Kubernetes cluster. |
| [1.2](https://github.com/ipersids/devops-kubernetes/tree/main/todo_app) | [1.2](https://github.com/ipersids/devops-kubernetes/tree/1.2) | Deploy a web server (Todo app) that outputs "Server started in port NNNN" when it is started. |
| [1.3](https://github.com/ipersids/devops-kubernetes/tree/main/log_output) | [1.3](https://github.com/ipersids/devops-kubernetes/tree/1.3) | Log output app: move the deployment into a declarative file. |
| [1.4](https://github.com/ipersids/devops-kubernetes/tree/main/todo_app) | [1.4](https://github.com/ipersids/devops-kubernetes/tree/1.4) | Todo app: move the deployment into a declarative file. |
| [1.5](https://github.com/ipersids/devops-kubernetes/tree/main/todo_app) | [1.5](https://github.com/ipersids/devops-kubernetes/tree/1.5) | Todo app: add GET / endpoint. |
| [1.6](https://github.com/ipersids/devops-kubernetes/tree/main/todo_app) | [1.6](https://github.com/ipersids/devops-kubernetes/tree/1.6) | Use a NodePort Service to enable access to the Todo app. |
| [1.7](https://github.com/ipersids/devops-kubernetes/tree/main/log_output) | [1.7](https://github.com/ipersids/devops-kubernetes/tree/1.7) | Add an endpoint to request the current status in the Log output app and an Ingress to access it with a browser. |
| [1.8](https://github.com/ipersids/devops-kubernetes/tree/main/todo_app) | [1.8](https://github.com/ipersids/devops-kubernetes/tree/1.8) | Switch to using Ingress instead of NodePort to access the Todo app. |
| [1.9](https://github.com/ipersids/devops-kubernetes/tree/main/pingpong) | [1.9](https://github.com/ipersids/devops-kubernetes/tree/1.9) | Develop a second application that simply responds with "pong 0" to a GET request and increases a counter. |
| [1.10](https://github.com/ipersids/devops-kubernetes/tree/main/log_output) | [1.10](https://github.com/ipersids/devops-kubernetes/tree/1.10) | Split the Log output app into two containers in one pod and add a shared volume (emptyDir). |
| [1.11](https://github.com/ipersids/devops-kubernetes/tree/main/pingpong) | [1.11](https://github.com/ipersids/devops-kubernetes/tree/1.11) | Share data between "Ping-pong" and "Log output" applications using persistent volumes. |
| [1.12](https://github.com/ipersids/devops-kubernetes/tree/main/todo_app) | [1.12](https://github.com/ipersids/devops-kubernetes/tree/1.12) | Get a random picture from Lorem Picsum and display it in Todo app, update picture every 10 minutes. |
| [1.13](https://github.com/ipersids/devops-kubernetes/tree/main/todo_app) | [1.13](https://github.com/ipersids/devops-kubernetes/tree/1.13) | Update a Todo app functionality. |
| [2.1](https://github.com/ipersids/devops-kubernetes/tree/main/log_output) | [2.1](https://github.com/ipersids/devops-kubernetes/tree/2.1) | Connect the Log output application and the Ping-pong application with HTTP. Remove the shared volume between those two applications. |
