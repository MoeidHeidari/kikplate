# Kubernetes Deployment

This guide covers deploying Kikplate to a Kubernetes cluster using the manifests in `kubernetes/`. All base manifests are managed with Kustomize and live under `kubernetes/base/`. An overlay for development is provided at `kubernetes/overlays/dev/`.

## Prerequisites

| Tool | Version |
|------|---------|
| kubectl | 1.27 or later |
| kustomize | 5.x (or kubectl with built-in kustomize) |
| A running Kubernetes cluster | Any CNCF-conformant cluster |
| An nginx ingress controller | For external access |
| A container registry token | To pull from GHCR |

## Architecture in Kubernetes

Kikplate deploys as four workloads:

| Workload | Kind | Description |
|----------|------|-------------|
| `api` | Deployment | Go REST API server |
| `sync` | Deployment | Background sync worker (runs `app:sync` command) |
| `web` | Deployment | Next.js frontend |
| `postgres` | Deployment | PostgreSQL with a PersistentVolumeClaim |

The API and sync worker use the same container image (`ghcr.io/kikplate/kikplate-api`) but run different commands. Both mount the same ConfigMap and Secret.

## Resource Overview

```
kubernetes/base/
  namespace.yaml                  Creates the kikplate namespace
  postgres-secret.yaml            PostgreSQL credentials (POSTGRES_USER, POSTGRES_PASSWORD)
  postgres-pvc.yaml               10Gi PersistentVolumeClaim for PostgreSQL data
  postgres-service.yaml           ClusterIP service for PostgreSQL
  postgres-deployment.yaml        PostgreSQL Deployment
  app-config-configmap.yaml       config.yaml mounted at /app/config/config.yaml
  app-credentials-secret.yaml     JWT_SECRET and OAuth client secrets
  api-service.yaml                ClusterIP service for the API
  api-deployment.yaml             API server Deployment
  sync-deployment.yaml            Sync worker Deployment
  web-service.yaml                ClusterIP service for the web frontend
  web-deployment.yaml             Next.js frontend Deployment
  ingress.yaml                    Nginx Ingress routing traffic to the web frontend
```

## Step-by-Step Deployment

### 1. Create the Secrets

The example secret file is at `kubernetes/base/app-credentials-secret.yaml.example`. Copy it and fill in real values:

```
cp kubernetes/base/app-credentials-secret.yaml.example kubernetes/base/app-credentials-secret.yaml
```

Edit `kubernetes/base/app-credentials-secret.yaml`:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-credentials
  namespace: kikplate
type: Opaque
stringData:
  JWT_SECRET: "a-long-random-secret-minimum-32-characters"
  SSO_GITHUB_CLIENT_SECRET: "your-github-client-secret"
  SSO_GOOGLE_CLIENT_SECRET: "your-google-client-secret"
  SSO_GITLAB_CLIENT_SECRET: "your-gitlab-client-secret"
```

Generate a strong JWT secret:

```
openssl rand -base64 48
```

Also update the PostgreSQL secret at `kubernetes/base/postgres-secret.yaml` if you want non-default credentials.

### 2. Update the Application Config

The application config is stored in `kubernetes/base/app-config-configmap.yaml`. Edit it to set your domain, OAuth provider details, and customization values. It follows the same schema as `config/config.yaml`. See [Configuration](configuration.md) for the full reference.

Key values to change for a real deployment:

```yaml
server:
  frontend_url: https://kikplate.yourdomain.com

sso:
  providers:
    - name: github
      client_id: YOUR_GITHUB_CLIENT_ID
      redirect_url: https://kikplate.yourdomain.com/api/auth/github/callback
```

Leave `client_secret` empty in the ConfigMap. The secret is injected as an environment variable from `app-credentials-secret.yaml`.

### 3. Update the Ingress Hostname

Edit `kubernetes/base/ingress.yaml` and replace the example hostname:

```yaml
spec:
  rules:
    - host: kikplate.yourdomain.com
```

Ensure your DNS points to the ingress controller's external IP.

### 4. Apply the Manifests

```
kubectl apply -k kubernetes/base
```

Or to use the development overlay:

```
kubectl apply -k kubernetes/overlays/dev
```

Check that all pods are running:

```
kubectl get pods -n kikplate
```

Expected output once healthy:

```
NAME                        READY   STATUS    RESTARTS   AGE
api-xxxx-xxxx               1/1     Running   0          2m
sync-xxxx-xxxx              1/1     Running   0          2m
web-xxxx-xxxx               1/1     Running   0          2m
postgres-xxxx-xxxx          1/1     Running   0          2m
```

### 5. Seed the Database

Run a one-off job to seed the badge catalog:

```
kubectl exec -n kikplate deploy/api -- ./api db:seed
```

## Configuration Updates

When you change the ConfigMap, restart affected Deployments to pick up the new config:

```
kubectl rollout restart deployment/api -n kikplate
kubectl rollout restart deployment/sync -n kikplate
```

When a Secret changes, restart similarly.

## Scaling

### Horizontal Scaling

The API and web deployments are stateless and can be scaled horizontally:

```
kubectl scale deployment/api --replicas=3 -n kikplate
kubectl scale deployment/web --replicas=3 -n kikplate
```

The sync worker must run as a single replica unless your workload requires multiple workers. Multiple sync workers can run safely as each plate row is locked atomically in the database before processing, but you should confirm your batch size and poll interval support it.

### Resource Requests and Limits

Default resource budgets defined in the base manifests:

| Workload | CPU Request | CPU Limit | Memory Request | Memory Limit |
|----------|-------------|-----------|---------------|-------------|
| api | 100m | 500m | 128Mi | 512Mi |
| sync | 100m | 500m | 128Mi | 512Mi |
| web | 100m | 500m | 256Mi | 512Mi |
| postgres | 100m | 500m | 256Mi | 512Mi |

Adjust these using a Kustomize patch in your overlay. An example patch is provided at `kubernetes/overlays/dev/api-deployment-patch.yaml`.

## Kustomize Overlays

The base manifests are environment-agnostic. Create overlays for different environments in `kubernetes/overlays/<env>/`. An overlay typically:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
  - ../../base
patchesStrategicMerge:
  - api-deployment-patch.yaml
  - sync-deployment-patch.yaml
  - web-deployment-patch.yaml
```

A patch file can override any field. For example, to change the API replica count and image tag for production:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
  namespace: kikplate
spec:
  replicas: 3
  template:
    spec:
      containers:
        - name: api
          image: ghcr.io/kikplate/kikplate-api:v1.2.0
```

## Health Checks

The API exposes a health endpoint at `GET /hello`. Both the readiness and liveness probes for the API deployment target this path.

Check pod health manually:

```
kubectl exec -n kikplate deploy/api -- wget -qO- http://localhost:3001/hello
```

## Viewing Logs

```
kubectl logs -n kikplate deploy/api
kubectl logs -n kikplate deploy/sync
kubectl logs -n kikplate deploy/web
kubectl logs -n kikplate deploy/postgres
```

Follow logs:

```
kubectl logs -n kikplate deploy/api -f
```

## Upgrading

To deploy a new image tag, update the image in the relevant deployment and apply:

```
kubectl set image deployment/api api=ghcr.io/kikplate/kikplate-api:v1.3.0 -n kikplate
kubectl set image deployment/sync sync=ghcr.io/kikplate/kikplate-api:v1.3.0 -n kikplate
kubectl set image deployment/web web=ghcr.io/kikplate/kikplate-web:v1.3.0 -n kikplate
```

Or update your Kustomize overlay and re-apply.

## Uninstalling

To remove all Kikplate resources from the cluster:

```
kubectl delete -k kubernetes/base
```

This deletes all resources including the PVC. To keep the database volume, delete individual resources instead.
