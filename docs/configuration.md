# Configuration

Kikplate is configured through a YAML file and a set of environment variables that override or extend it. The lookup order is: environment variables from `.env` file, then the YAML file, then hardcoded defaults.

## Config File Location

The API server resolves the config file in the following order:

1. The path set in the `CONFIG_PATH` environment variable.
2. `./config/config.yaml` relative to the working directory.
3. `../config/config.yaml`
4. `../../config/config.yaml`

An example file is provided at `config/examples/config.yaml.example`.

## Complete Reference

### database

PostgreSQL connection settings.

```yaml
database:
  host: localhost          # Hostname or IP of the PostgreSQL server
  port: 5432               # PostgreSQL port
  database: kikplate       # Database name
  username: postgres       # Database user
  password: password       # Database password
```

All database fields can also be set via environment variables:

| Environment Variable | Config Key |
|---------------------|-----------|
| `DB_HOST` | `database.host` |
| `DB_PORT` | `database.port` |
| `DB_NAME` | `database.database` |
| `DB_USER` | `database.username` |
| `DB_PASS` | `database.password` |

### server

HTTP server settings.

```yaml
server:
  port: 3001                          # Port the API listens on
  frontend_url: http://localhost:3000 # Used for OAuth callbacks and CORS
  log:
    level: info                       # trace, debug, info, warn, error
```

| Environment Variable | Config Key |
|---------------------|-----------|
| `SERVER_PORT` | `server.port` |
| `SERVER_LOG_LEVEL` | `server.log.level` |

### sync

Controls the background synchronization worker (`app:sync`).

```yaml
sync:
  interval: 20m      # How often a plate is eligible for re-sync after it was last synced
  poll_interval: 5m  # How often the worker wakes up and looks for due plates
  batch_size: 25     # Maximum number of plates processed per poll cycle
```

These values are intentionally conservative. In a deployment with thousands of plates you may want to reduce `poll_interval` and increase `batch_size`.

### sso

OAuth provider configuration. You can configure any combination of GitHub, Google, and GitLab.

```yaml
sso:
  providers:
    - name: github
      client_id: YOUR_CLIENT_ID
      client_secret: YOUR_CLIENT_SECRET
      redirect_url: http://localhost:3000/api/auth/github/callback
      scopes:
        - read:user
        - user:email

    - name: google
      client_id: YOUR_CLIENT_ID
      client_secret: YOUR_CLIENT_SECRET
      redirect_url: http://localhost:3000/api/auth/google/callback
      scopes:
        - openid
        - email
        - profile

    - name: gitlab
      client_id: YOUR_CLIENT_ID
      client_secret: YOUR_CLIENT_SECRET
      redirect_url: http://localhost:3000/api/auth/gitlab/callback
      scopes:
        - openid
        - email
        - profile
```

Client secrets should never be committed to source control. In production, set them via Kubernetes secrets or environment variables and leave them empty in the config file.

### auth.email_verification

Controls local (non-SSO) email/password signup behavior.

```yaml
auth:
  email_verification:
    enabled: false
    token_ttl: 24h
    verify_url_base: http://localhost:3000/verify-email
```

- `enabled=false`: users registered through `POST /auth/register` are active immediately.
- `enabled=true`: users are created inactive and must verify via emailed link before login.
- `verify_url_base` is optional. If omitted, the API uses `server.frontend_url + /verify-email`.

| Environment Variable | Config Key |
|---------------------|-----------|
| `AUTH_EMAIL_VERIFICATION_ENABLED` | `auth.email_verification.enabled` |
| `AUTH_EMAIL_VERIFICATION_TOKEN_TTL` | `auth.email_verification.token_ttl` |
| `AUTH_EMAIL_VERIFICATION_VERIFY_URL_BASE` | `auth.email_verification.verify_url_base` |

### smtp

SMTP sender configuration used when `auth.email_verification.enabled=true`.

```yaml
smtp:
  host: email-smtp.us-east-1.amazonaws.com
  port: 587
  username: your_ses_smtp_username
  password: your_ses_smtp_password
  from_email: no-reply@example.com
  from_name: Kikplate
  use_starttls: true
```

For Amazon SES, use the SMTP endpoint of your SES region (for example `email-smtp.us-east-1.amazonaws.com`) with SES SMTP credentials.

| Environment Variable | Config Key |
|---------------------|-----------|
| `SMTP_HOST` | `smtp.host` |
| `SMTP_PORT` | `smtp.port` |
| `SMTP_USERNAME` | `smtp.username` |
| `SMTP_PASSWORD` | `smtp.password` |
| `SMTP_FROM_EMAIL` | `smtp.from_email` |
| `SMTP_FROM_NAME` | `smtp.from_name` |
| `SMTP_USE_STARTTLS` | `smtp.use_starttls` |

### Authentication Header

To enable trusted reverse-proxy header authentication:

```
AUTH_HEADER=X-Remote-User
```

When this variable is set, the API reads that header on every request. If present, it resolves or creates an account for the value. This mode is additive: JWT and OAuth continue to work alongside it.

### JWT

```
JWT_SECRET=a-long-random-string-at-least-32-characters
```

This secret signs all JWT access tokens. Rotating it invalidates all existing sessions. In production use a randomly generated value of at least 64 characters.

### customization

UI customization applied to the web frontend. These values are served from `GET /api/config` and consumed by the Next.js app.

```yaml
customization:
  logo: /kikplate-logo-on-dark.svg
  banner_title: "The Home of your starter boilerplates"
  badge_request_url: "https://github.com/MoeidHeidari/kikplate/issues/new?template=badge-request.yml"
  social_media:
    - type: github
      link: "https://github.com/yourorg"
    - type: slack
      link: "https://yourworkspace.slack.com"
    - type: linkedin
      link: "https://linkedin.com/company/yourcompany"
    - type: x
      link: "https://x.com/yourhandle"
  prepared_queries:
    - "Golang starter"
    - "Next.js"
    - "Gin framework"
```

The `logo` field accepts any public URL, a CDN path, or a path served by the Next.js frontend.

### badges

The badge catalog. These are seeded into the database the first time `db:seed` runs. The badge definition is for display and documentation only. The actual award of a badge to a plate is stored in the `plate_badge` table.

```yaml
badges:
  - slug: official
    name: Official
    description: Officially recognized and maintained by the project owners
    icon: award
    tier: official

  - slug: production-ready
    name: Production Ready
    description: Ready for production environments with high reliability
    icon: rocket
    tier: official

  - slug: documented
    name: Documented
    description: Well-documented with clear and comprehensive documentation
    icon: book-open
    tier: community
```

Supported `tier` values are `official` and `community`. Badge icons use the Lucide icon name format.

## Kubernetes and Helm

In Kubernetes, the config file is stored in a ConfigMap and mounted at `/app/config/config.yaml`. Secrets such as `JWT_SECRET` and OAuth client secrets are stored in a Kubernetes Secret and injected as environment variables.

See [Kubernetes](kubernetes.md) and [Helm](helm.md) for deployment-specific configuration guidance.

## Environment Variable Summary

| Variable | Description | Default |
|----------|-------------|---------|
| `CONFIG_PATH` | Absolute path to the YAML config file | Auto-detected |
| `DB_HOST` | PostgreSQL hostname | From config |
| `DB_PORT` | PostgreSQL port | From config |
| `DB_NAME` | PostgreSQL database name | From config |
| `DB_USER` | PostgreSQL username | From config |
| `DB_PASS` | PostgreSQL password | From config |
| `SERVER_PORT` | API listen port | `3001` |
| `SERVER_LOG_LEVEL` | Log verbosity | `info` |
| `JWT_SECRET` | JWT signing secret | Required |
| `AUTH_HEADER` | Trusted header name for reverse-proxy auth | Disabled |
| `AUTH_EMAIL_VERIFICATION_ENABLED` | Enable email verification for local signup | `false` |
| `AUTH_EMAIL_VERIFICATION_TOKEN_TTL` | Verification token expiration duration | `24h` |
| `AUTH_EMAIL_VERIFICATION_VERIFY_URL_BASE` | Base verification URL in outgoing email | `FRONTEND_URL/verify-email` |
| `SMTP_HOST` | SMTP host for verification emails | None |
| `SMTP_PORT` | SMTP port for verification emails | `587` |
| `SMTP_USERNAME` | SMTP username | None |
| `SMTP_PASSWORD` | SMTP password | None |
| `SMTP_FROM_EMAIL` | Sender email address | None |
| `SMTP_FROM_NAME` | Sender display name | `Kikplate` |
| `SMTP_USE_STARTTLS` | Use STARTTLS instead of implicit TLS | `true` |
| `ENV` | Environment name (`development`, `production`) | `development` |
