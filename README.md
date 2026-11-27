# devops-kubernetes

Submissions for the DevOps with Kubernetes course at the University of Helsinki

## Setup

#### 1. Install prerequisites

- Docker
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [k3d](https://github.com/k3d-io/k3d#get)
- [kubectx + kubens](https://github.com/ahmetb/kubectx) (optional)
- [SOPS: Secrets OPerationS](https://github.com/getsops/sops)
- [age: File encryption](https://github.com/FiloSottile/age) (optional)

Verify the installation:

```bash
docker --version
kubectl version --client
k3d version
sops --version --check-for-updates
```

#### 2. Create a Kubernetes cluster

```bash
k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2
# create folder for PersistentVolume
docker exec k3d-k3s-default-agent-0 mkdir -p /tmp/kube
```

The manifests include a local `PersistentVolume` that is pinned to the first k3d agent node:

```yaml
k3d-k3s-default-agent-0
```

If your cluster/node names are different, update `manifests/persistentvolume.yaml` before applying it.

#### 3. Apply manifests

```bash
kubectl apply -f manifests/

export SOPS_AGE_KEY_FILE=$(pwd)/secrets/key.txt
sops --decrypt pingpong/manifests/enc/secret.enc.yaml | kubectl apply -f -
sops --decrypt todo_backend/manifests/enc/secret.enc.yaml | kubectl apply -f -

kubectl apply -f todo_app/manifests/todoapp-config.yaml
kubectl apply -f todo_app/manifests/persistentvolumeclaim.yaml
kubectl apply -f todo_app/manifests

kubectl apply -f todo_backend/manifests/todobackend-config.yaml
kubectl apply -f todo_backend/manifests/database.yaml
kubectl apply -f todo_backend/manifests

kubectl apply -f log_output/manifests/logoutput-config.yaml
kubectl apply -f log_output/manifests

kubectl apply -f pingpong/manifests/pingpong-config.yaml
kubectl apply -f pingpong/manifests/database.yaml
kubectl apply -f pingpong/manifests

kubectl apply -f todo_cron/manifests/todo-cron.yaml
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

#### 5. Enable monitoring using Grafana

Follow the instructions in the [monitoring README](https://github.com/ipersids/devops-kubernetes/tree/main/monitoring/README.md).

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
| [1.12](https://github.com/ipersids/devops-kubernetes/tree/main/todo_app) | [1.12](https://github.com/ipersids/devops-kubernetes/tree/1.12) | Get a random picture from Lorem Picsum, display it in the Todo app, and update it every 10 minutes. |
| [1.13](https://github.com/ipersids/devops-kubernetes/tree/main/todo_app) | [1.13](https://github.com/ipersids/devops-kubernetes/tree/1.13) | Update the Todo app's functionality. |
| [2.1](https://github.com/ipersids/devops-kubernetes/tree/main/log_output) | [2.1](https://github.com/ipersids/devops-kubernetes/tree/2.1) | Connect the Log output application and the Ping-pong application with HTTP. Remove the shared volume between those two applications. |
| [2.2](https://github.com/ipersids/devops-kubernetes/tree/main/todo_backend) | [2.2](https://github.com/ipersids/devops-kubernetes/tree/2.2) | Create a Todo backend service. It should have a GET /todos endpoint for fetching the list of todos and a POST /todos endpoint for creating a new todo. |
| [2.3](https://github.com/ipersids/devops-kubernetes/tree/main/log_output) | [2.3](https://github.com/ipersids/devops-kubernetes/tree/2.3) | Move the Log output and Ping-pong apps to the `exercises` namespace. |
| [2.4](https://github.com/ipersids/devops-kubernetes/tree/main/todo_app) | [2.4](https://github.com/ipersids/devops-kubernetes/tree/2.4) | Move the Todo client and Todo backend to the `project` namespace. |
| [2.5](https://github.com/ipersids/devops-kubernetes/tree/main/log_output) | [2.5](https://github.com/ipersids/devops-kubernetes/tree/2.5) | Create a ConfigMap for the Log output application to define one file information.txt and env variable MESSAGE. |
| [2.6](https://github.com/ipersids/devops-kubernetes/tree/main/todo_app) | [2.6](https://github.com/ipersids/devops-kubernetes/tree/2.6) | Pass all the configurations to Todo as env variables that are defined either in a config map or in deployments. |
| [2.7](https://github.com/ipersids/devops-kubernetes/tree/main/pingpong) | [2.7](https://github.com/ipersids/devops-kubernetes/tree/2.7) | Run a Postgres database as a StatefulSet (with one replica) and save the Ping-pong application's counter in the database. |
| [2.8](https://github.com/ipersids/devops-kubernetes/tree/main/todo_backend) | [2.8](https://github.com/ipersids/devops-kubernetes/tree/2.8) | Run a Postgres database as a stateful set (with one replica) and save the Todo tasks into the database. |
| [2.9](https://github.com/ipersids/devops-kubernetes/tree/main/todo_cron) | [2.9](https://github.com/ipersids/devops-kubernetes/tree/2.9) | Create a CronJob that generates a new todo every hour to remind to read a random wikipedia article. |
| [2.10](https://github.com/ipersids/devops-kubernetes/tree/main/monitoring) | [2.10](https://github.com/ipersids/devops-kubernetes/tree/2.10) | Add request logging so that you can monitor every todo that is sent to the backend. |
| [3.1](https://github.com/ipersids/devops-kubernetes/tree/main/pingpong) | [release 3.1](https://github.com/ipersids/devops-kubernetes/tree/3.11) | Deploy the Ping-pong application to GKE. |
| [3.2](https://github.com/ipersids/devops-kubernetes/tree/main/log_output) | [3.2](https://github.com/ipersids/devops-kubernetes/tree/3.2) | Deploy the "Log output" and "Ping-pong" applications to GKE and expose them with Ingress. "Ping-pong" responds at `/pingpong`. |
| [3.3](https://github.com/ipersids/devops-kubernetes/tree/main/log_output) | [3.3](https://github.com/ipersids/devops-kubernetes/tree/3.3) | Replace the Ingress with Gateway API in the "Log output" and "Ping-pong" applications. |
| [3.4](https://github.com/ipersids/devops-kubernetes/tree/main/pingpong) | [3.4](https://github.com/ipersids/devops-kubernetes/tree/3.4) | Move Ping-pong's counter increment handler to `/` and configure the HTTPRoute to rewrite `/pingpong` to `/`. |
| [3.5](https://github.com/ipersids/devops-kubernetes/tree/main/todo_app) | [3.5](https://github.com/ipersids/devops-kubernetes/tree/3.5) | Configure the project to use Kustomize, and deploy it to Google Kubernetes Engine. |
| [3.6](https://github.com/ipersids/devops-kubernetes/tree/main/.github/workflows) | [3.6](https://github.com/ipersids/devops-kubernetes/tree/3.6) | Setup automatic deployment for the project using GitHub Actions. |
| [3.7](https://github.com/ipersids/devops-kubernetes/tree/main/.github/workflows) | [3.7](https://github.com/ipersids/devops-kubernetes/tree/3.7) | Improve the deployment so that each branch creates a separate environment. The main branch should still be deployed in the namespace project. |
| [3.8](https://github.com/ipersids/devops-kubernetes/tree/main/.github/workflows) | [3.8](https://github.com/ipersids/devops-kubernetes/tree/3.8) | Create a new workflow so that deleting a branch deletes the environment. |
| [3.9](#DBaaS-vs-DIY) | [3.9](https://github.com/ipersids/devops-kubernetes/tree/3.9) | Do a pros/cons comparison of the DBaaS and DIY solutions. |

<a id="DBaaS-vs-DIY"></a>
## DBaaS vs DIY

Database-as-a-service (DBaaS) is a managed database model provided as a service. DBaaS typically runs on a cloud computing platform and allows users to pay a cloud provider for database services and computing resources. This helps users focus more on development and outsource database management to the provider’s team.

Do-it-Yourself (DIY) means a self-hosted model where the operational work to set up and maintain the database is handled by an in-house team. It gives more control over infrastructure, compute, storage, network security, backups and software upgrades.

| Topic | Database-as-a-service | Self-Hosted Database |
| :---: | :-------------------: | :------------------- |
| Cost | Depends on compute, storage and managed features, including the provider’s management fees. | Infrastructure may cost less, but the overall cost includes operational work and engineering time to build and maintain the database infrastructure. |
| Setup | Choose a vendor, create an instance and configure access — usually fast and simple. | Choose infrastructure, configure the database, storage, networking, secrets and deployment — requires planning and development. |
| Maintenance | The provider’s team specialises in databases and handles routine maintenance. The client’s team still plans application compatibility, upgrades and maintenance timing. | An in-house team plans and handles database updates, patches and upgrades. Additional coordination with other development teams may be needed. |
| Backups and recovery | Managed backup and restore tools work according to the client’s configuration. The client is responsible for testing recovery. | The in-house team is responsible for creating backup jobs, storing backups and testing recovery. |
| Availability | Managed replication and automatic failover can reduce downtime when high availability is configured. | To ensure high availability, the team needs to configure replication, failover and routing, and perform failure testing. Backups are also needed for data recovery. A single database pod has downtime during recovery. |
| Security | The provider manages infrastructure security. The client is responsible for managing users, permissions, connections and data access. | In-house team manage database security as well as its deployment configuration. |
| Control | Limited by the vendor, which provides a high-level abstraction over the database. | More control over low-level database configuration, hosting and security features. |
| Monitoring | Built-in metrics and logs are available. The client still needs to configure useful alerts. | Additional work is needed to install and configure monitoring tools and alerts. |
| Required skills | 	Less database operations knowledge is needed, but database design and SQL skills remain necessary. A separate database operations team may not be needed. | Requires database, storage, networking, backup and recovery knowledge. Kubernetes knowledge is needed if the database runs in Kubernetes. Team size may grow with business needs. |


**In the DBaaS model**, vendors handle the underlying infrastructure, database engine maintenance and infrastructure security, and provide automated backups and recovery tools. The vendor creates a high-level abstraction over the database and provides user-friendly tools to configure it, monitor metrics and logs, and set up alerts. Clients still need to configure access, backup settings and recovery requirements.

This helps clients focus more on developing their core product and can reduce the need for a separate database operations team. DBaaS solutions are usually fast and simple to set up, which makes them a good choice when development speed is more important than low-level control over the database infrastructure.

However, there are some disadvantages as well. Vendors include management costs in their prices and can change prices over time. Less control over the database infrastructure can limit some business needs and desired database features or configurations. Another major risk is vendor lock-in, especially when the application depends on vendor-specific features. This can make it difficult to move away later when the service or its prices no longer meet business needs.

**In the DIY model**, the in-house team handles database setup, maintenance, security, backups and recovery. This requires more planning, technical knowledge and coordination with other development teams. Monitoring tools, alerts, replication and failover also need to be configured and maintained by the team.

This gives the business more control over database infrastructure, versions, extensions and configuration. DIY can be a good choice when specific business needs require features or settings that managed services do not support. It also gives more freedom to choose hosting and move the database to another provider.

However, this control comes with additional responsibilities. Infrastructure costs may be lower, but the overall cost includes engineering time and ongoing operational work. Updates, failures and recovery can take time away from developing the core product. The team also needs to test backups and failure scenarios regularly, because running the database successfully does not guarantee that its data can be recovered after a problem.

**Which approach works better depends on business needs**. Small startups may choose the DBaaS model because of limited team size and the need for fast development. Other companies with high database costs, specific technical needs or strict "no third-party" requirements for data storage may consider a self-hosted database a better choice. Both options should be compared based on the required level of control, total costs and the team’s ability to manage the database.
