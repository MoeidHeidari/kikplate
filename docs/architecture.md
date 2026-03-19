# Kickplate — Architecture

> **Living document.** Grows incrementally alongside the codebase. Each section is added as the system evolves.
> Database design is covered separately in `database.md`.

---

## Table of Contents

1. [System Overview](#system-overview)
2. [Repository Structure](#repository-structure)
3. [Backend Layer Architecture](#backend-layer-architecture)
4. [Dependency Injection — fx Wiring](#dependency-injection--fx-wiring)
5. [Auth Flow](#auth-flow)
6. [Plate Lifecycle](#plate-lifecycle)
7. [Synchronizer Design](#synchronizer-design)
8. [Request Lifecycle](#request-lifecycle)
9. [Deployment Topology](#deployment-topology)
10. [Roadmap Phases](#roadmap-phases)

---

## System Overview

Kickplate is an open-source template registry — a community-driven marketplace where developers discover, share, and scaffold production-ready project templates via a web UI or CLI. Think *npm for project templates*.

```mermaid
graph TD
    subgraph Clients["Clients — Phase 1"]
        WEB["🌐 Web UI\n(Next.js + Tailwind + shadcn/ui)"]
    end

    subgraph CLI_FUTURE["Clients — Phase 2"]
        CLI["🖥️ kickplate CLI\n(Go + Cobra)\n[planned]"]
    end

    subgraph Backend["API Server  ·  Go + Chi + fx"]
        API["REST API\n(handlers → services → repositories)"]
        SYNC["Synchronizer\n(background goroutine)"]
        RL["Rate Limiter\n(Redis-based, per IP + per user)"]
    end

    subgraph Storage
        PG[("PostgreSQL\n(plates, users, reviews, badges)")]
        REDIS[("Redis\n(sessions, rate limiting)")]
    end

    subgraph External
        GH["GitHub Public API\n(fetch kickplate.yaml + README)"]
        IDP["Identity Provider\n(SSO — config-driven)"]
    end

    WEB -->|"HTTP + JWT"| API
    CLI -.->|"HTTP + JWT\n[phase 2]"| API
    API --> PG
    API --> REDIS
    API --> RL
    SYNC -->|"GET kickplate.yaml"| GH
    SYNC --> PG
    API <-->|"SSO token exchange"| IDP

    style CLI_FUTURE opacity:0.45
```

---

## Repository Structure

```mermaid
graph LR
    subgraph root["kickplate/"]
        direction TB
        API_DIR["api/\nGo backend"]
        WEB_DIR["web/\nNext.js frontend"]
        CLI_DIR["cli/\nGo CLI  [phase 2]"]
        DC["docker-compose.yml"]
        CFG["config.yaml"]
    end

    subgraph api_detail["api/"]
        direction TB
        BOOT["bootstrap/\nfx wiring, root command"]
        CMD["command/\nCobra sub-commands"]
        LIB["lib/\nLogger · Env · Database · RequestHandler"]
        MDL["model/\nGORM structs (11 entities)"]
        REPO["repository/\nInterfaces + postgres/ implementations"]
        SVC["service/\n[planned]"]
        HDL["handler/handlers/\n[planned]"]
        RTE["handler/routes/\n[planned]"]
    end

    API_DIR --> api_detail

    style CLI_DIR opacity:0.45
    style SVC opacity:0.6
    style HDL opacity:0.6
    style RTE opacity:0.6
```

---

## Backend Layer Architecture

Dependencies flow **downward only**. No layer imports from a layer above it. Uber fx enforces this at wire-up time.

```mermaid
flowchart TB
    subgraph Transport["① Transport"]
        CHI["Chi Router"]
        MW["Middleware\nRequestID · RealIP · Recoverer · Auth [→]"]
    end

    subgraph Handler["② Handler  [→]"]
        AH["AuthHandler"]
        PH["PlateHandler"]
        ADM["AdminHandler"]
    end

    subgraph Service["③ Service  [→]"]
        AS["AuthService"]
        PS["PlateService"]
        SS["SyncService"]
    end

    subgraph Repository["④ Repository"]
        RI["Interfaces\nUserRepo · AccountRepo · PlateRepo\nPlateTagRepo · PlateReviewRepo\nBadgeRepo · PlateBadgeRepo\nPlateMemberRepo · SyncLogRepo\nEmailVerificationRepo"]
        RP["Postgres Implementations\n(GORM + WithContext)"]
    end

    subgraph Infra["⑤ Infrastructure"]
        DB["lib.Database\ngorm.DB wrapping PostgreSQL"]
        RDS["Redis\n(via go-redis)  [→]"]
        ENV["lib.Env\nconfig.yaml + env var overrides"]
        LOG["lib.Logger\n(zap — dev colour / prod JSON)"]
    end

    CHI --> MW --> Handler
    Handler --> Service
    Service --> Repository
    RI --> RP
    RP --> DB
    Service --> RDS
    DB -.->|reads DSN| ENV
    Handler & Service & Repository -.->|injected| LOG

    style Transport fill:#1e3a5f,color:#e8f0fe
    style Handler fill:#1a4731,color:#d4edda
    style Service fill:#3d2b00,color:#fff3cd
    style Repository fill:#3b1f2b,color:#f8d7da
    style Infra fill:#212121,color:#e0e0e0
```

> `[→]` = planned, not yet implemented.

---

## Dependency Injection — fx Wiring

All modules are composed in `bootstrap/modules.go` and injected by Uber fx at startup.

```mermaid
graph TD
    subgraph bootstrap["bootstrap.CommonModules"]
        LM["lib.Module\nLogger · Env · Database · RequestHandler"]
        RM["routes.Module\nAll route registrations"]
        HM["handlers.Module\nAll HTTP handlers"]
        SM["service.Module  [→]"]
        RPM["repository.Module  [→]"]
    end

    LM --> HM
    LM --> RM
    LM --> RPM
    RPM --> SM
    SM --> HM
    HM --> RM

    subgraph cmd["command/"]
        SERVE["ServeCommand\nfx.Invoke(cmd.Run())"]
    end

    bootstrap --> SERVE
```

---

## Auth Flow

Three auth modes all resolve to a single `account.id` used by every downstream table. The mode is selected by `config.yaml` — no code change needed to switch modes.

```mermaid
sequenceDiagram
    actor User
    participant API as API Server
    participant DB as PostgreSQL
    participant IDP as Identity Provider

    note over User,DB: Mode 1 — Session (username / password)
    User->>API: POST /auth/register {username, email, password}
    API->>DB: INSERT user (is_active=false)
    API->>DB: INSERT email_verification (token, expires_at)
    API-->>User: 201 — check your email

    User->>API: GET /auth/verify?token=…
    API->>DB: find token WHERE is_used=false AND expires_at > now()
    API->>DB: UPDATE user SET is_active=true
    API->>DB: UPDATE email_verification SET is_used=true
    API->>DB: find or create account WHERE provider='local'
    API-->>User: JWT { account_id }

    note over User,IDP: Mode 2 — SSO (OAuth / OIDC)
    User->>API: GET /auth/sso/callback?code=…
    API->>IDP: exchange code → ID token
    IDP-->>API: { sub, email, name }
    API->>DB: find or create account WHERE provider='github', provider_user_id=sub
    API-->>User: JWT { account_id }

    note over User,DB: Mode 3 — Header (proxy / internal)
    User->>API: Any request + X-Auth-User: value
    API->>DB: find or create account WHERE provider='header', provider_user_id=value
    API->>API: continue with account.id
```

---

## Plate Lifecycle

```mermaid
stateDiagram-v2
    direction LR

    [*] --> Pending : User submits plate

    Pending --> Approved  : Admin approves
    Pending --> Rejected  : Admin rejects\nor ownership mismatch
    Rejected --> Pending  : Owner fixes + resubmits
    Approved --> Archived : Owner or admin archives

    state Approved {
        direction TB
        [*] --> Public
        Public --> Private   : Owner changes visibility
        Public --> Unlisted  : Owner changes visibility
        Private --> Public
        Unlisted --> Public
    }

    note right of Pending
        Repository plate:
        backend fetches kickplate.yaml immediately,
        checks owner == authenticated username.
        ✓ match → is_verified=true, sync_status=synced
        ✗ mismatch → rejected
        File plate:
        always is_verified=true (implicit ownership)
    end note
```

---

## Synchronizer Design

Runs as a background goroutine on a 1-minute ticker. Only acts on repository plates — file plates have no remote source to sync.

```mermaid
flowchart TD
    TICK["⏱️ Ticker\nevery 1 minute"]

    QUERY["SELECT plates\nWHERE type = 'repository'\nAND next_sync_at ≤ now()\nAND sync_status ≠ 'syncing'\nORDER BY next_sync_at ASC\nLIMIT 50"]

    WORKERS["Parallel goroutines\none per plate"]
    LOCK["SET sync_status = 'syncing'"]

    FETCH["GET github.com/…/kickplate.yaml\n(single service-level token from config.yaml)"]

    OK{"HTTP 200?"}

    OWNER{"owner field\nstill matches?"}
    CHANGED{"metadata\nchanged?"}

    UPDATE_META["UPDATE plate\nSET metadata, updated_at"]
    MARK_SYNCED["sync_status = 'synced'\nconsecutive_failures = 0\nnext_sync_at = now() + sync_interval\nlast_synced_at = now()"]
    LOG_OK["INSERT sync_log\n(status=success, changes_detected)"]

    UNVERIFIED["is_verified = false\nsync_status = 'unverified'\nnotify owner"]
    LOG_UV["INSERT sync_log\n(status=failed, 'owner mismatch')"]

    GONE{"Repo gone or\nyaml missing?"}
    FAIL["sync_status = 'failed'\nconsecutive_failures += 1\nnext_sync_at = backoff(n)\nsync_error = message"]
    LOG_FAIL["INSERT sync_log\n(status=failed, error)"]
    LOG_GONE["INSERT sync_log\n(status=failed, 'repo not found')"]
    GONE_STATUS["sync_status = 'unverified'\nnotify owner"]

    TICK --> QUERY --> WORKERS --> LOCK --> FETCH
    FETCH --> OK
    OK -->|Yes| OWNER
    OWNER -->|Yes| CHANGED
    OWNER -->|No| UNVERIFIED --> LOG_UV
    CHANGED -->|Yes| UPDATE_META --> MARK_SYNCED --> LOG_OK
    CHANGED -->|No| MARK_SYNCED
    OK -->|No| GONE
    GONE -->|No| FAIL --> LOG_FAIL
    GONE -->|Yes| GONE_STATUS --> LOG_GONE
```

**Backoff schedule**

| Consecutive failures | Next retry delay |
|---|---|
| 1 | 30 minutes |
| 2 | 2 hours |
| 3 | 12 hours |
| 4+ | 24 hours + owner notification |

---

## Request Lifecycle

A typical authenticated API request from ingress to storage and back.

```mermaid
sequenceDiagram
    participant C as Client (Browser)
    participant CHI as Chi Router
    participant MW as Middleware Stack
    participant H as Handler  [→]
    participant SVC as Service  [→]
    participant REPO as Repository
    participant DB as PostgreSQL
    participant RDS as Redis

    C->>CHI: HTTP request + Bearer JWT
    CHI->>MW: RequestID · RealIP · Recoverer
    MW->>RDS: rate limit check (per IP + per user)  [→]
    MW->>MW: verify JWT, extract account.id  [→]
    MW->>H: r.Context() carries account.id
    H->>H: decode + validate request body
    H->>SVC: call service method
    SVC->>REPO: interface method
    REPO->>DB: GORM query (WithContext)
    DB-->>REPO: rows / error
    REPO-->>SVC: model structs / error
    SVC-->>H: result / domain error
    H->>C: JSON  200 / 4xx / 5xx
```

---

## Deployment Topology

Self-hosting is a first-class requirement. Target experience: `docker-compose up` and done.

```mermaid
graph TD
    subgraph Host["Docker host"]
        subgraph DC["docker-compose stack"]
            WEB_C["web\nNext.js · port 3000"]
            API_C["api\nGo binary · port 8080"]
            PG_C["postgres\nport 5432"]
            RDS_C["redis\nport 6379"]
        end
        CFG_VOL[("config.yaml\nvolume mount")]
        PG_VOL[("postgres data\nnamed volume")]
    end

    BROWSER["Browser"] -->|":80 / :443"| WEB_C
    BROWSER -->|"direct API calls"| API_C
    WEB_C -->|"SSR API calls"| API_C
    API_C --> PG_C
    API_C --> RDS_C
    CFG_VOL -.-> API_C
    PG_C --- PG_VOL

    subgraph CI["GitHub Actions CI"]
        BUILD["build + test + lint\non every PR"]
    end
```

```bash
git clone https://github.com/kickplate/kickplate
cd kickplate
cp config.example.yaml config.yaml   # fill in secrets
docker-compose up
# → running at http://localhost:8080
```

---

## Roadmap Phases

```mermaid
gantt
    title Kickplate — Roadmap
    dateFormat  YYYY-MM
    axisFormat  %b %Y

    section Phase 1 · Foundation
    Backend core (auth, plates, sync, admin)   :active, p1b, 2025-01, 3M
    Frontend (home, listing, detail, submit)   :active, p1f, 2025-01, 3M
    Docker Compose + GitHub Actions CI         :active, p1d, 2025-01, 3M
    Seed 20-30 plates                          :        p1s, 2025-03, 1M

    section Phase 2 · CLI
    kickplate scaf + search + submit           :p2, 2025-04, 2M

    section Phase 3 · Community
    Reviews, ratings, badges, profiles         :p3, 2025-06, 3M
    Trending, notifications                    :p3b, 2025-07, 2M

    section Phase 4 · Organizations
    Org namespaces, private plates, SSO        :p4, 2025-09, 4M

    section Phase 5 · Advanced Scaffolding
    values.yaml, multi-plate projects, versioning :p5, 2025-11, 4M

    section Phase 6 · AI
    Plate matching + composition + gap-fill    :p6, 2026-03, 4M
```

---

*Last updated: initial draft — system overview, repo structure, layers, fx wiring, auth, plate lifecycle, synchronizer, request lifecycle, deployment, roadmap.*