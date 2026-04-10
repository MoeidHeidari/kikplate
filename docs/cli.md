# CLI Reference

The Kikplate CLI (`kikplate`) is a standalone Go binary that lets you browse, submit, manage, and scaffold from plates on any Kikplate server.

## Installation

### From a Release (Recommended)

Download the binary for your platform from the [GitHub Releases](https://github.com/kikplate/kikplate/releases) page.

```
tar -xzf kikplate-<version>-linux-amd64.tar.gz
sudo mv kikplate /usr/local/bin/kikplate
```

Supported platforms: Linux amd64, Linux arm64, macOS amd64, macOS arm64.

### From Source

```
git clone https://github.com/kikplate/kikplate.git
cd kikplate/cli
go build -o kikplate .
sudo mv kikplate /usr/local/bin/kikplate
```

## Configuration

The CLI stores its configuration at `~/.kikplate/config.yaml`.

### Initialize

```
kikplate config init
```

This creates a default config file:

```yaml
server:
  address: http://localhost:3001
auth:
  token: ""
```

### Use a Custom Config Path

Every command accepts a global `--config` flag:

```
kikplate --config /path/to/config.yaml search --name golang
```

### View Active Config

```
kikplate config view
```

### Config File Reference

| Key | Description | Default |
|-----|-------------|---------|
| `server.address` | Base URL of the Kikplate API server | `http://localhost:3001` |
| `auth.token` | JWT stored after login. Set automatically by `kikplate login`. | `""` |

## Authentication

### Login

```
kikplate login --email you@example.com --password yourpassword
```

On success, the JWT is saved to `~/.kikplate/config.yaml`. Subsequent commands use it automatically.

### Who Am I

```
kikplate whoami
```

Prints the username, display name, email, and account ID of the authenticated account.

### Logout

```
kikplate logout
```

Removes the stored token from the config file.

## Searching and Browsing

### Search Plates

```
kikplate search [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--name` | Search by name or keyword | |
| `--category` | Filter by category | |
| `--tag` | Filter by tag | |
| `--limit` | Results per page | 20 |
| `--page` | Page number | 1 |

Examples:

```
kikplate search --name "golang rest"
kikplate search --category backend --tag gin
kikplate search --name next.js --limit 10 --page 2
```

Output:

```
SLUG                   NAME                 CATEGORY   RATING   VERIFIED
myorg/go-chi-rest      Go Chi REST API      backend    4.8      yes
myorg/gin-starter      Gin Starter          backend    4.5      yes

Total: 2
```

### Describe a Plate

```
kikplate describe myorg/go-chi-rest
```

Output includes name, category, status, visibility, repository URL, branch, rating, star count, verification status, tags, owner, organization, and timestamps.

## Managing Your Plates

### View Your Plates

```
kikplate my plates
```

### View Your Bookmarks

```
kikplate my bookmarks
```

### View Your Organizations

```
kikplate my orgs
```

### View Sync Status

```
kikplate my sync
```

Shows each plate's sync status, last sync time, and next scheduled sync.

## Managing Local Plates

The CLI maintains a local list of plates you use frequently at `~/.kikplate/plates.json`. This is separate from the server.

### Add a Plate Locally

```
kikplate plates add myorg/go-chi-rest
```

Fetches the plate details from the server and stores them locally.

### List Local Plates

```
kikplate plates list
```

### Remove a Local Plate

```
kikplate plates remove myorg/go-chi-rest
```

## Submitting a Plate

Submission requires authentication.

```
kikplate submit https://github.com/myorg/my-template
```

Optional flags:

| Flag | Description |
|------|-------------|
| `--branch` | Git branch to use (default: `main`) |
| `--org` | Organization ID to submit under |

Example for an organization-scoped plate:

```
kikplate submit https://github.com/myorg/my-template --org <org-uuid> --branch main
```

After submission, the CLI prints the assigned slug and the `verification_token` you need to add to your `kikplate.yaml`.

## Verifying a Plate

After adding the verification token to `kikplate.yaml` and pushing the commit:

```
kikplate verify myorg/my-template
```

On success the plate becomes approved, public, and verified. The CLI confirms the new status.

## Scaffolding a Project

Scaffolding creates a new project from a plate. The source repository is cloned, the `kikplate.yaml` manifest is stripped, and `.kikplate.origin` is written to record the provenance.

### Scaffold to a Local Directory

```
kikplate scaffold myorg/go-chi-rest my-new-project
```

### Scaffold to a Remote Repository

```
kikplate scaffold myorg/go-chi-rest https://github.com/you/new-repo.git
```

This clones the template, creates an initial commit, and pushes to the remote. If the remote already contains a scaffolded project, the command refuses to push unless you pass `--force`.

Optional flags:

| Flag | Description |
|------|-------------|
| `--local` | Force scaffold to local directory even if a URL-shaped argument is given |
| `--force` | Overwrite a remote that was already scaffolded from a kikplate plate |

## Global Flags

| Flag | Description |
|------|-------------|
| `--config` | Path to config file |
| `--help`, `-h` | Show help for any command |

## Command Reference Summary

| Command | Description | Auth Required |
|---------|-------------|--------------|
| `config init` | Create default config file | No |
| `config view` | Print active config | No |
| `login` | Authenticate and save token | No |
| `logout` | Remove stored token | No |
| `whoami` | Show authenticated user | Yes |
| `search` | Search plates on the server | No |
| `describe [slug]` | Show plate details | No |
| `submit [repo-url]` | Submit repository as plate | Yes |
| `verify [slug]` | Verify submitted plate | Yes |
| `scaffold [slug] [target]` | Scaffold project from plate | No |
| `plates add [slug]` | Add plate to local list | No |
| `plates list` | List locally added plates | No |
| `plates remove [slug]` | Remove plate from local list | No |
| `my plates` | List your plates on server | Yes |
| `my bookmarks` | List bookmarked plates | Yes |
| `my orgs` | List your organizations | Yes |
| `my sync` | View sync status of your plates | Yes |
