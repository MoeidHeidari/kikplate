# Kikplate Helm Chart

Kikplate is a self hosted platform for discovering, sharing, and submitting starter templates. This chart deploys the full Kikplate stack on Kubernetes.

## What This Chart Deploys

1. `api` service for the Go backend
2. `web` service for the Next.js frontend
3. `sync` worker for background synchronization
4. `postgres` for the application database
5. `ingress` for external access when enabled

## Install

```bash
helm install kikplate oci://ghcr.io/kikplate/helm-charts/kikplate \
  --namespace kikplate \
  --create-namespace
```

## Install With Required Secrets

```bash
helm install kikplate oci://ghcr.io/kikplate/helm-charts/kikplate \
  --namespace kikplate \
  --create-namespace \
  --set secrets.jwtSecret="$(openssl rand -base64 32)" \
  --set secrets.sso.githubClientSecret="YOUR_GITHUB_CLIENT_SECRET"
```

## Upgrade

```bash
helm upgrade kikplate oci://ghcr.io/kikplate/helm-charts/kikplate \
  --namespace kikplate
```

## Uninstall

```bash
helm uninstall kikplate --namespace kikplate
```

## Configuration

The following table lists the main configurable parameters and their default values.

| Parameter | Description | Default / Example |
|-----------|-------------|-------------------|
| `nameOverride` | Optional short name override for generated resources | `""` |
| `fullnameOverride` | Optional full name override for generated resources | `""` |
| `namespace` | Optional namespace override. Leave empty to use the Helm release namespace | `""` |
| `api.replicaCount` | Number of backend replicas | `1` |
| `api.image.repository` | Backend image repository | `ghcr.io/kikplate/kikplate-api` |
| `api.image.tag` | Backend image tag | `main` |
| `api.image.pullPolicy` | Backend image pull policy | `IfNotPresent` |
| `api.service.port` | Backend service port | `3001` |
| `api.resources.requests.cpu` | Backend CPU request | `100m` |
| `api.resources.requests.memory` | Backend memory request | `128Mi` |
| `api.resources.limits.cpu` | Backend CPU limit | `500m` |
| `api.resources.limits.memory` | Backend memory limit | `512Mi` |
| `sync.replicaCount` | Number of sync worker replicas | `1` |
| `sync.image.repository` | Sync worker image repository | `ghcr.io/kikplate/kikplate-api` |
| `sync.image.tag` | Sync worker image tag | `main` |
| `sync.image.pullPolicy` | Sync worker image pull policy | `IfNotPresent` |
| `web.replicaCount` | Number of frontend replicas | `1` |
| `web.image.repository` | Frontend image repository | `ghcr.io/kikplate/kikplate-web` |
| `web.image.tag` | Frontend image tag | `main` |
| `web.image.pullPolicy` | Frontend image pull policy | `IfNotPresent` |
| `web.service.port` | Frontend service port | `3000` |
| `postgres.image.repository` | PostgreSQL image repository | `postgres` |
| `postgres.image.tag` | PostgreSQL image tag | `16` |
| `postgres.service.port` | PostgreSQL service port | `5432` |
| `postgres.persistence.size` | PostgreSQL persistent volume size | `10Gi` |
| `postgres.persistence.accessMode` | PostgreSQL persistent volume access mode | `ReadWriteOnce` |
| `ingress.enabled` | Enable ingress resource generation | `true` |
| `ingress.className` | Ingress class name | `nginx` |
| `ingress.host` | Public ingress hostname | `kikplate.example.com` |
| `secrets.jwtSecret` | JWT signing secret for the application | `change-me-jwt-secret` |
| `secrets.sso.githubClientSecret` | GitHub OAuth client secret | `""` |
| `secrets.sso.googleClientSecret` | Google OAuth client secret | `""` |
| `secrets.sso.gitlabClientSecret` | GitLab OAuth client secret | `""` |
| `postgresql.database` | PostgreSQL database name | `kikplate` |
| `postgresql.username` | PostgreSQL username | `postgres` |
| `postgresql.password` | PostgreSQL password | `password` |
| `config.server.port` | Backend listen port written to config | `3001` |
| `config.server.frontendUrl` | Public frontend URL used by backend callbacks and links | `http://localhost:3000` |
| `config.server.log.level` | Backend log level | `info` |
| `config.sync.interval` | Full synchronization interval | `20m` |
| `config.sync.pollInterval` | Polling interval for updates | `5m` |
| `config.sync.batchSize` | Batch size for synchronization jobs | `25` |
| `config.customization.logo` | Logo path used by the frontend | `/kikplate-logo-on-dark.svg` |
| `config.customization.bannerTitle` | Main homepage banner title | `The Home of your starter boilerplates` |
| `config.customization.badgeRequestUrl` | URL used for badge requests | `https://github.com/kikplate/kikplate/issues/new?template=badge-request.yml` |
| `config.customization.socialMedia` | Social links shown in the UI | `github`, `slack`, `linkedin`, `x` entries |
| `config.customization.preparedQueries` | Prepared search suggestions shown in the homepage search | See `values.yaml` |
| `config.sso.providers` | OAuth provider definitions rendered into config | `github`, `google`, `gitlab` |
| `config.badges` | Badge catalog rendered by the application | See badge table below |

## Badge Value Schema

Artifact Hub reads these descriptions from `values.schema.json` and can display them for badge entries.

| Badge Field | Description | Example |
|-------------|-------------|---------|
| `config.badges[].slug` | Stable unique badge identifier used internally and in URLs | `official` |
| `config.badges[].name` | Display label shown in the application UI | `Official` |
| `config.badges[].description` | Longer explanation of what the badge means | `Officially recognized and maintained by the project owners` |
| `config.badges[].icon` | Frontend icon identifier used for the badge | `award` |
| `config.badges[].tier` | Badge ownership tier | `official` or `community` |

## Values File

For the complete defaults and schema metadata, see the packaged `values.yaml` and `values.schema.json` files included in this chart.