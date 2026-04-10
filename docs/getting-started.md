# Getting Started

This guide walks you through running Kikplate on your local machine using Docker Compose.

## Prerequisites

The following tools must be installed before you begin:

| Tool | Minimum Version | Purpose |
|------|----------------|---------|
| Docker | 24.x | Container runtime |
| Docker Compose | 2.x | Multi-container orchestration |
| Go | 1.22 | Building the API and CLI |
| Node.js | 20 LTS | Building the web frontend |
| Git | 2.x | Cloning the repository |

## Clone the Repository

```
git clone https://github.com/kikplate/kikplate.git
cd kikplate
```

## Configure the Application

Copy the example configuration file:

```
cp config/examples/config.yaml.example config/config.yaml
```

Open `config/config.yaml` in your editor and set at minimum:

```yaml
database:
  host: localhost
  port: 5432
  database: kikplate
  username: postgres
  password: password

server:
  port: 3001
  frontend_url: http://localhost:3000

sso:
  providers:
    - name: github
      client_id: YOUR_GITHUB_CLIENT_ID
      client_secret: YOUR_GITHUB_CLIENT_SECRET
      redirect_url: http://localhost:3000/api/auth/github/callback
      scopes:
        - read:user
        - user:email
```

For a complete reference of every configuration key see [Configuration](configuration.md).

## Start the Stack

```
docker compose up --build
```

This starts four services:

| Service | Port | Description |
|---------|------|-------------|
| postgres | 5432 | PostgreSQL database |
| api | 3001 | Go REST API |
| sync | (internal) | Background sync worker |
| web | 3000 | Next.js frontend |

The web UI is available at `http://localhost:3000`.
The API is available at `http://localhost:3001`.

## Seed Initial Data

The API automatically runs GORM auto-migration on startup. To seed badge catalog data run:

```
docker compose exec api ./api db:seed
```

## Register an Account

Open `http://localhost:3000` and register a new account. You will receive an email verification link. In local development, check the container logs for the verification URL:

```
docker compose logs api | grep verify-email
```

## Install the CLI

The Kikplate CLI lets you interact with any Kikplate server from your terminal. To install it from source:

```
cd cli
go build -o kikplate .
sudo mv kikplate /usr/local/bin/kikplate
```

Pre-built binaries for Linux and macOS (amd64 and arm64) are attached to every GitHub release.

After installing, initialize the CLI configuration:

```
kikplate config init
```

This creates `~/.kikplate/config.yaml`. Edit it to point to your server:

```yaml
server:
  address: http://localhost:3001
```

Then log in:

```
kikplate login --email you@example.com --password yourpassword
```

For full CLI documentation see [CLI Reference](cli.md).

## Development Setup

### API

```
cd api
go mod download
go run . app:serve
```

The API reads `config/config.yaml` by default. Override the path with the `CONFIG_PATH` environment variable.

### Web

```
cd web
npm install
npm run dev
```

The development server starts at `http://localhost:3000` with hot reload.

### Sync Worker

```
cd api
go run . app:sync
```

The sync worker can run alongside the API server. They share the same database and configuration.

## Running Tests

### API

```
cd api
go test ./...
```

### Web

```
cd web
npm run build
```

### CLI

```
cd cli
go test ./...
```

## Next Steps

Once everything is running, read:

- [How It Works](how-it-works.md) to understand the plate lifecycle.
- [CLI Reference](cli.md) to submit and manage plates from the terminal.
- [Configuration](configuration.md) for production tuning.
- [Kubernetes](kubernetes.md) or [Helm](helm.md) to deploy to a cluster.
