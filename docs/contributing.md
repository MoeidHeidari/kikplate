# Contributing

Thank you for your interest in contributing to Kikplate. This document explains everything you need to know to propose changes, run the project locally, and get your contribution merged.

## Code of Conduct

All contributors are expected to be respectful and constructive. Harassment of any kind is not tolerated.

## Ways to Contribute

- Report bugs by opening a GitHub issue.
- Propose features by opening a discussion or issue before writing code.
- Improve documentation.
- Submit pull requests that fix bugs, add tests, or implement agreed-upon features.

## Repository Structure

```
kikplate/
  api/              Go REST API server and sync worker
    bootstrap/      Cobra root command and Uber Fx assembly
    command/        CLI commands registered with Fx
    handler/        HTTP routes, handlers, and middleware
    lib/            Infrastructure: config, DB, logger, JWT
    model/          GORM models
    repository/     Data access layer (PostgreSQL implementations)
    service/        Business logic
  cli/              Standalone Go CLI
    cmd/            Cobra commands
  web/              Next.js frontend
  helm/             Helm chart
  kubernetes/       Kustomize base and overlays
  config/           Application configuration
  docs/             Documentation
  .github/          GitHub Actions workflows
```

## Setting Up the Development Environment

Follow [Getting Started](getting-started.md) to run the stack locally with Docker Compose. For active development of individual components, run each service directly.

### API (Go)

```
cd api
go mod download
go run . app:serve
```

Run tests:

```
go test ./...
```

Build the binary:

```
go build -o api .
```

Configuration is auto-detected from `config/config.yaml`. Override with `CONFIG_PATH`.

### Web (Next.js)

```
cd web
npm install
npm run dev
```

Run production build check:

```
npm run build
```

### CLI (Go)

```
cd cli
go mod download
go run . --help
```

Run tests:

```
go test ./...
```

## Making Changes

### Branching Strategy

| Branch | Purpose |
|--------|---------|
| `main` | Production-ready code. All releases come from here. |
| `pre-release` | Release candidates built by the RC workflow. |
| Feature branches | One branch per issue or feature. Branch from `main`. |

### Naming Branches

Use a short, lowercase, hyphenated name that reflects the work:

```
feat/badge-request-flow
fix/sync-backoff-overflow
docs/cli-reference
```

### Commit Messages

Follow the Conventional Commits format. The release workflow reads commit history to generate changelogs.

| Type | Use When |
|------|---------|
| `feat:` | Adding a new feature |
| `fix:` | Fixing a bug |
| `docs:` | Documentation only |
| `refactor:` | Code change with no behavior change |
| `test:` | Adding or updating tests |
| `chore:` | Build, tooling, CI changes |

Examples:

```
feat: add gitlab oauth provider support
fix: prevent sync worker from processing archived plates
docs: add helm deployment guide
```

## Adding a New API Endpoint

1. Define the model in `api/model/` if a new table is needed.
2. Add a repository method in `api/repository/postgres/`.
3. Add a service method in `api/service/`.
4. Create a handler in `api/handler/handlers/`.
5. Register the route in `api/handler/routes/module.go`.
6. Register any new providers in the relevant Fx module.
7. Add the endpoint to `docs/openapi.yaml`.
8. Update `docs/architecture.md` if the HTTP surface section changes.

## Adding a New CLI Command

1. Create a new file in `cli/cmd/` (e.g., `cli/cmd/mycommand.go`).
2. Define a `cobra.Command` variable.
3. Register it in `init()` by calling `rootCmd.AddCommand(myCmd)`.
4. Add command documentation to `docs/cli.md`.

## Pull Request Process

1. Fork the repository and create your branch from `main`.
2. Ensure the full test suite passes locally before opening a PR.
3. Write or update tests for any code you change.
4. Update documentation for any user-facing changes.
5. Open the PR against the `main` branch.
6. The PR title should follow the same Conventional Commits format as commits.
7. A maintainer will review your PR. Address all review comments before it can be merged.
8. Squash merges are used to keep the `main` branch history clean.

## CI Checks

Every pull request runs the CI workflow which:

| Job | What it checks |
|-----|---------------|
| API Quality | `go test ./...` and `go build ./...` in `api/` |
| Web Quality | `npm run build` in `web/` |
| CLI Quality | `go test ./...` and `go build ./...` in `cli/` |

All three jobs must pass before a PR can be merged. There are no exceptions.

Quality checks are implemented as reusable workflows in `.github/workflows/reusable-*-quality.yml`. If you need to change how tests run, edit those files. The change applies to all workflows that call them.

## Release Process

Releases are created manually by maintainers from the GitHub Actions UI. The Release workflow (`release.yml`) handles:

1. Validating that the workflow runs from `main`.
2. Computing the next semantic version from the latest git tag.
3. Running all quality jobs.
4. Building and pushing multi-platform Docker images to GHCR.
5. Building CLI binaries for Linux and macOS (amd64 and arm64).
6. Creating a GitHub Release with attached binaries.

The Helm chart is released separately using the Helm Chart Release workflow (`helm-release.yml`). It validates the requested version against an existing application release, runs Helm lint and template checks, packages the chart, and pushes it to both GHCR and the gh-pages Helm repository.

As a contributor you do not need to trigger releases. Maintainers handle that.

## Project Conventions

### Go

Modules use the `github.com/kickplate/api` import path for the API and `github.com/kikplate/cli` for the CLI. Follow standard Go project layout conventions. Avoid adding global state.

### Error Handling

Return errors up the call stack. Do not log and return simultaneously. Handlers are responsible for translating service errors into HTTP status codes.

### Database Migrations

GORM `AutoMigrate` runs on startup. Add new fields directly to the model structs. Backward-incompatible changes (column renames, type changes) require a migration strategy discussion with maintainers before implementation.

### Configuration

New configuration keys belong in `api/lib/env.go` (in the `Env` struct) and must be documented in `docs/configuration.md`. Add the corresponding key to `config/examples/config.yaml.example` with a comment explaining its purpose.

## Reporting Bugs

Open a GitHub issue and include:

1. A short description of the bug.
2. Steps to reproduce it.
3. Expected behavior.
4. Actual behavior.
5. Relevant logs or error messages.
6. Kikplate version or commit hash.
7. Operating system and Kubernetes version if relevant.

## Requesting Features

Open a GitHub discussion or issue. Describe the use case and proposed behavior before writing any code. This avoids effort on features that do not align with the project direction.

## License

Kikplate is licensed under the MIT License. By contributing, you agree that your contributions will be licensed under the same terms.
