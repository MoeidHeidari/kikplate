<div align="center">
  <img src="docs/images/logo.svg" alt="Kikplate Logo" width="200"/>
  <br><br>
  <a href="https://join.slack.com/t/kikplate/shared_invite/zt-3uddz8y51-4ko_TBozERrr7jhl2AzYSg">
    <img src="https://img.shields.io/badge/Join%20kikplate%20on%20Slack-4A154B?style=for-the-badge&logo=slack&logoColor=white" alt="Join on Slack"/>
  </a>


  [![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/kikplate)](https://artifacthub.io/packages/search?repo=kikplate)
  [![GitHub stars](https://img.shields.io/github/stars/kikplate/kikplate?style=social)](https://github.com/kikplate/kikplate)
  
  
# Kikplate
  
  Kikplate is an open source marketplace for project templates and boilerplates. Teams publish templates backed by real GitHub repositories, developers browse and filter by category, language, or badge tier, and the CLI scaffolds new projects in seconds.

  <img src="docs/images/kikplate.png" alt="Kikplate Logo" />

<br><br>
</div>



## Overview

| Topic | Description |
|-------|-------------|
| **Templates as code** | Every template ("plate") is a public GitHub repository containing a `kikplate.yaml` manifest. There is nothing to upload — the source of truth is always the repository. |
| **Verification workflow** | Submitting a plate issues a one-time token. Placing that token in `kikplate.yaml` and calling `verify` proves ownership and moves the plate to the approved state. |
| **Continuous sync** | The background sync worker periodically re-reads each approved plate's `kikplate.yaml`, updating metadata and catching drift or revocation automatically. |
| **Organizations** | Plates can be scoped to an organization, enabling teams and companies to maintain shared template catalogs under a single namespace. |
| **Badges** | Admins award badges (community, verified, official, sponsored) to signal quality and curation level. |
| **CLI** | The `kikplate` CLI handles authentication, searching, describing, and scaffolding templates without touching a browser. |

---

## Quick Start

The fastest way to run Kikplate locally is with Docker Compose.

```sh
git clone https://github.com/your-org/kikplate.git
cd kikplate

cp config/examples/config.yaml.example config/config.yaml
# Edit config/config.yaml — set database.dsn and jwt_secret at minimum.

docker compose up -d
docker compose exec api ./api db:seed
```

The API is available at `http://localhost:3001` and the web UI at `http://localhost:3000`.

Install the CLI to interact from a terminal:

```sh
go install github.com/your-org/kikplate/cli@latest

kikplate config init
kikplate login
kikplate search --category backend
```

---

## Documentation

| Document | Description |
|----------|-------------|
| [Getting Started](docs/getting-started.md) | Prerequisites, local setup, dev workflow, running tests |
| [How It Works](docs/how-it-works.md) | Core concepts: plates, kikplate.yaml, verification lifecycle, sync, badges |
| [Architecture](docs/architecture.md) | Component map, layers, HTTP surface, DI modules, data flow diagrams |
| [Database](docs/database.md) | Schema reference, table definitions, entity relationships |
| [Configuration](docs/configuration.md) | All config file keys and environment variable overrides |
| [CLI Reference](docs/cli.md) | Full reference for every `kikplate` command with flags and examples |
| [Kubernetes](docs/kubernetes.md) | Deploying to Kubernetes with `kubectl apply -k`, secrets, ingress, scaling |
| [Helm](docs/helm.md) | Helm chart install, upgrade, full values reference, production example |
| [Contributing](docs/contributing.md) | Repo structure, branching, commit conventions, adding endpoints, PR process, release flow |
| [OpenAPI](docs/openapi.yaml) | Machine-readable API specification (OpenAPI 3.0.3) |

---

## Repository Layout

```
api/          Go API server (Chi, Uber Fx, GORM)
cli/          Go CLI client (Cobra)
web/          Next.js frontend
config/       Default configuration file and examples
docs/         All documentation including OpenAPI spec
helm/         Helm chart for production deployments
kubernetes/   Kustomize manifests for direct kubectl deployments
.github/      CI/CD workflows
```

---

## License

[LICENSE](LICENSE)

