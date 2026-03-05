# Write API

Command & mutation service for the ODA poetry platform. Handles authentication, poem CRUD, likes, bookmarks, emotion tagging, and admin operations.

## Stack

| Layer | Tech |
|-------|------|
| Router | Chi v5 |
| Framework | Huma v2 (typed I/O, auto OpenAPI) |
| ORM | GORM + pgx v5 |
| DB | PostgreSQL |
| Auth | JWT (access + refresh tokens), bcrypt |

## Architecture

```
cmd/server/main.go          ← Entrypoint, wiring, graceful shutdown
internal/
├── config/                  ← Env-based configuration
├── database/                ← GORM connection + migrations
├── middleware/               ← Auth, security, logging (Chi + Huma)
├── seed/                    ← Admin user + emotion catalog seeding
└── features/
    ├── auth/                ← Register, login, refresh, profile, users
    │   ├── delivery/http/   ← Huma handlers + route registration
    │   ├── repository/      ← GORM persistence (users, refresh tokens)
    │   └── usecase/         ← Business logic (auth flows, profile mgmt)
    ├── poems/               ← Poem CRUD, likes, bookmarks, emotions
    │   ├── delivery/http/
    │   ├── repository/
    │   └── usecase/
    └── admin/               ← Dashboard stats, user/poem/association mgmt
        ├── delivery/http/
        ├── repository/
        └── usecase/
```

Each feature follows **Clean Architecture**: `repository → usecase → delivery/http`.

Files are split by domain concern (e.g., `poems_handler.go`, `interactions_handler.go`) rather than monolithic single files.

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | Yes | — | PostgreSQL connection string |
| `PORT` | No | `8082` | HTTP listen port |
| `JWT_SECRET` | Yes | — | Secret for JWT signing |
| `INTERNAL_SECRET` | Yes | — | Shared secret for service-to-service auth |
| `ADMIN_EMAIL` | No | `admin@oda.local` | Seeded admin email |
| `ADMIN_PASSWORD` | No | — | Seeded admin password |

## API Docs

Auto-generated OpenAPI spec available at `/docs` when the service is running.

## Endpoints

### Auth (`/api/auth/*`)
- `POST /auth/register` — Create account
- `POST /auth/login` — Authenticate
- `POST /auth/refresh` — Rotate tokens
- `POST /auth/logout` — Revoke refresh token

### Profile (`/api/me`, `/api/auth/*`, `/api/users/*`)
- `GET /me` — Current user profile
- `PUT /auth/profile` — Update profile
- `POST /auth/change-password` — Change password
- `GET /users/:username` — Public profile
- `GET /users/search` — Search users

### Poems (`/api/poems/*`)
- `POST /poems` — Create poem
- `PUT /poems/:id` — Update poem
- `DELETE /poems/:id` — Delete poem
- `POST /poems/:id/like` — Toggle like
- `POST /poems/:id/bookmark` — Toggle bookmark
- `POST /poems/:id/emotions` — Tag emotion
- `DELETE /poems/:id/emotions` — Remove current user's emotion tag

### Admin (`/api/admin/*`)
- Dashboard stats, CRUD for users/poems, manage likes/bookmarks/emotions, emotion catalog management

## Running

```bash
# Standalone
DATABASE_URL="..." JWT_SECRET="..." INTERNAL_SECRET="..." go run ./cmd/server

# With Air (hot-reload) — configured via .air.toml
air
```
