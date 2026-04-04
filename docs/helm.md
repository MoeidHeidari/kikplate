# Helm Deployment

Kikplate provides an official Helm chart that deploys the full stack — API, sync worker, web frontend, and PostgreSQL — with a single command.

## Prerequisites

| Tool | Version |
|------|---------|
| Helm | 3.x |
| kubectl | 1.27 or later |
| A Kubernetes cluster | Any CNCF-conformant cluster |
| An nginx ingress controller | For external access |

## Installing the Chart

### From the OCI Registry (Recommended)

The chart is published to GHCR on every release. No repository registration needed:

```
helm install kikplate oci://ghcr.io/moeidheidari/helm-charts/kikplate \
  --namespace kikplate \
  --create-namespace
```

With required secrets set at install time:

```
helm install kikplate oci://ghcr.io/moeidheidari/helm-charts/kikplate \
  --namespace kikplate \
  --create-namespace \
  --set secrets.jwtSecret="$(openssl rand -base64 32)" \
  --set secrets.sso.githubClientSecret="YOUR_GITHUB_CLIENT_SECRET"
```

### From the Helm Repository

```
helm repo add kikplate https://moeidheidari.github.io/kikplate
helm repo update
helm install kikplate kikplate/kikplate \
  --namespace kikplate \
  --create-namespace
```

## Upgrading

```
helm upgrade kikplate oci://ghcr.io/moeidheidari/helm-charts/kikplate \
  --namespace kikplate
```

Or from the repository:

```
helm repo update
helm upgrade kikplate kikplate/kikplate --namespace kikplate
```

## Uninstalling

```
helm uninstall kikplate --namespace kikplate
```

This removes all Helm-managed resources. The PostgreSQL PersistentVolumeClaim is not deleted automatically. To remove data permanently:

```
kubectl delete pvc -n kikplate -l app.kubernetes.io/instance=kikplate
```

## Values Reference

The chart exposes the following configurable values. Pass them with `--set key=value` or in a `values.yaml` file using `-f values.yaml`.

### nameOverride and fullnameOverride

```yaml
nameOverride: ""
fullnameOverride: ""
```

Override the base name used for generated Kubernetes resource names.

### namespace

```yaml
namespace: ""
```

Optional namespace override. By default, Helm uses the namespace passed to `helm install --namespace`.

### API

```yaml
api:
  replicaCount: 1
  image:
    repository: ghcr.io/moeidheidari/kikplate-api
    tag: main
    pullPolicy: IfNotPresent
  service:
    port: 3001
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: 500m
      memory: 512Mi
```

### Sync Worker

```yaml
sync:
  replicaCount: 1
  image:
    repository: ghcr.io/moeidheidari/kikplate-api
    tag: main
    pullPolicy: IfNotPresent
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: 500m
      memory: 512Mi
```

The sync worker uses the same image as the API server. The chart sets `args: ["app:sync"]` to run the sync command.

### Web Frontend

```yaml
web:
  replicaCount: 1
  image:
    repository: ghcr.io/moeidheidari/kikplate-web
    tag: main
    pullPolicy: IfNotPresent
  service:
    port: 3000
  resources:
    requests:
      cpu: 100m
      memory: 256Mi
    limits:
      cpu: 500m
      memory: 512Mi
```

### PostgreSQL

```yaml
postgres:
  replicaCount: 1
  image:
    repository: postgres
    tag: "16"
    pullPolicy: IfNotPresent
  service:
    port: 5432
  resources:
    requests:
      cpu: 100m
      memory: 256Mi
    limits:
      cpu: 500m
      memory: 512Mi
  persistence:
    size: 10Gi
    accessMode: ReadWriteOnce
```

### Ingress

```yaml
ingress:
  enabled: true
  className: nginx
  host: kikplate.example.com
  annotations: {}
```

Set `host` to your actual domain. Add cert-manager annotations here for TLS:

```yaml
ingress:
  enabled: true
  className: nginx
  host: kikplate.yourdomain.com
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
```

### Secrets

```yaml
secrets:
  jwtSecret: "change-me-jwt-secret"
  sso:
    githubClientSecret: ""
    googleClientSecret: ""
    gitlabClientSecret: ""
```

Never commit real secrets. Pass them at install time with `--set` or use an external secret management solution such as External Secrets Operator or Vault.

### PostgreSQL Credentials

```yaml
postgresql:
  database: kikplate
  username: postgres
  password: password
```

### Application Config

The `config` block maps directly to the `config.yaml` structure and is rendered into a ConfigMap mounted at `/app/config/config.yaml`:

```yaml
config:
  server:
    port: 3001
    frontendUrl: "https://kikplate.yourdomain.com"
    log:
      level: info
  sync:
    interval: 20m
    pollInterval: 5m
    batchSize: 25
  sso:
    providers:
      - name: github
        clientId: "YOUR_GITHUB_CLIENT_ID"
        redirectUrl: "https://kikplate.yourdomain.com/api/auth/github/callback"
        scopes:
          - read:user
          - user:email
  customization:
    logo: /kikplate-logo-on-dark.svg
    bannerTitle: "The Home of your starter boilerplates"
    badgeRequestUrl: "https://github.com/MoeidHeidari/kikplate/issues/new?template=badge-request.yml"
```

## Production Values File Example

Create a `production-values.yaml` file and pass it to Helm:

```yaml
api:
  replicaCount: 2

web:
  replicaCount: 2

ingress:
  enabled: true
  className: nginx
  host: kikplate.yourdomain.com
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod

postgresql:
  password: "a-strong-generated-password"

config:
  server:
    frontendUrl: "https://kikplate.yourdomain.com"
  sso:
    providers:
      - name: github
        clientId: "Ov23lid..."
        redirectUrl: "https://kikplate.yourdomain.com/api/auth/github/callback"
        scopes:
          - read:user
          - user:email
```

Install:

```
helm install kikplate oci://ghcr.io/moeidheidari/helm-charts/kikplate \
  --namespace kikplate \
  --create-namespace \
  -f production-values.yaml \
  --set secrets.jwtSecret="$(openssl rand -base64 48)" \
  --set secrets.sso.githubClientSecret="YOUR_SECRET"
```

## Scaling

To scale the API or web frontend after installation:

```
kubectl scale deployment -n kikplate -l app.kubernetes.io/component=api --replicas=3
kubectl scale deployment -n kikplate -l app.kubernetes.io/component=web --replicas=3
```

Or pass updated replicas to `helm upgrade`:

```
helm upgrade kikplate oci://ghcr.io/moeidheidari/helm-charts/kikplate \
  --namespace kikplate \
  --set api.replicaCount=3 \
  --set web.replicaCount=3
```

## Viewing Chart Values

```
helm show values oci://ghcr.io/moeidheidari/helm-charts/kikplate
```

## Checking Deployment Status

```
kubectl get pods -n kikplate
kubectl rollout status deployment/kikplate-api -n kikplate
kubectl rollout status deployment/kikplate-web -n kikplate
```

## Debugging

Show rendered templates without installing:

```
helm template kikplate oci://ghcr.io/moeidheidari/helm-charts/kikplate \
  --namespace kikplate \
  -f production-values.yaml
```

Lint the chart locally:

```
helm lint helm/
```
