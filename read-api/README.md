# Read API

Read-only query service for the ODA poetry platform. Serves the public feed, poem/user search, stats, bookmarks, and emotion catalog. Designed to scale independently from the write path and optionally connect to a read replica.

## Stack

| Layer | Tech |
|-------|------|
| Router | Chi v5 |
| Framework | Huma v2 (typed I/O, auto OpenAPI) |
| ORM | GORM + pgx v5 |
| DB | PostgreSQL (supports read replica) |

## Architecture

```
cmd/server/main.go          ← Entrypoint, wiring, graceful shutdown
internal/
├── config/                  ← Env-based configuration
├── database/                ← GORM connection
├── middleware/               ← Auth, security, logging (Chi + Huma)
└── features/feed/
    ├── delivery/http/       ← Huma handlers + route registration
    │   ├── handler_base.go  ← ReadHandler struct + constructor
    │   ├── poems_handler.go ← Feed, search, single poem, stats, emotions
    │   ├── users_handler.go ← Public profile, search users, user stats
    │   ├── bookmarks_handler.go
    │   ├── emotions_handler.go
    │   └── routes.go
    ├── repository/          ← GORM read queries
    │   ├── repository.go    ← ReadRepository struct + constructor
    │   ├── poems_repository.go
    │   ├── users_repository.go
    │   ├── bookmarks_repository.go
    │   └── emotions_repository.go
    └── usecase/             ← Thin pass-through with pagination
        ├── usecase.go
        ├── poems_usecase.go
        └── users_usecase.go
```

Single feature module (`feed`) containing all read-side domain concerns, split into per-domain files.

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL_READ` | No | fallback to `DATABASE_URL` | Read replica connection string |
| `DATABASE_URL` | Yes | — | Primary PostgreSQL connection |
| `PORT` | No | `8083` | HTTP listen port |
| `INTERNAL_SECRET` | Yes | — | Shared secret for service-to-service auth |

## API Docs

Auto-generated OpenAPI spec available at `/docs` when the service is running.

## Endpoints

### Poems
- `GET /api/poems/feed` — Public feed (paginated)
- `GET /api/poems/search?q=` — Search poems
- `GET /api/poems/:id` — Single poem
- `GET /api/poems/:id/stats` — Likes, views, emotion count
- `GET /api/poems/:id/emotions/distribution` — Emotion breakdown

### Users
- `GET /api/users/search?q=` — Search users
- `GET /api/users/:username` — Public profile
- `GET /api/users/:userID/poems` — User's poems
- `GET /api/users/:userID/stats` — Aggregate stats

### Bookmarks (authenticated)
- `GET /api/bookmarks` — Current user's bookmarks

### Emotions
- `GET /api/emotions` — Full emotion catalog

## Running

```bash
# Standalone
DATABASE_URL="..." INTERNAL_SECRET="..." go run ./cmd/server

# With Air (hot-reload) — configured via .air.toml
air
```
