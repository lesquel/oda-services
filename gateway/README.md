# Gateway

API gateway / reverse proxy for the ODA poetry platform. Single public entrypoint on port 8080 that routes requests to the write-api and read-api based on HTTP method and path.

## Stack

| Layer | Tech |
|-------|------|
| Router | Chi v5 |
| CORS | go-chi/cors |
| Auth | JWT validation (middleware) |
| Proxy | Custom reverse proxy with circuit breaker |

## Architecture

```
cmd/server/main.go          ← Entrypoint, route table, graceful shutdown
internal/
├── config/                  ← Env-based configuration
├── middleware/
│   ├── auth.go              ← JWT extraction + validation
│   ├── logger.go            ← Structured request logging
│   └── security.go          ← Security headers
└── proxy/
    ├── proxy.go             ← Reverse proxy (request forwarding)
    └── circuitbreaker.go    ← Circuit breaker pattern for resilience
```

## Routing Strategy

| Method / Path | Target | Auth |
|---------------|--------|------|
| `POST /api/auth/*` | write-api | None |
| `GET /api/*` (reads) | read-api | Optional / Required JWT |
| `POST/PUT/DELETE /api/*` (writes) | write-api | Required JWT |
| `GET/POST/PUT/DELETE /api/admin/*` | write-api | Admin JWT |
| `GET /docs` | Static HTML | None |
| `GET /healthz` | Local 200 | None |

The gateway injects the `X-Internal-Secret` header before forwarding, so backend services can trust that requests came through the gateway.

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PORT` | No | `8080` | HTTP listen port |
| `JWT_SECRET` | Yes | — | Secret for JWT validation |
| `INTERNAL_SECRET` | Yes | — | Shared secret injected into forwarded requests |
| `WRITE_API_URL` | Yes | `http://localhost:8082` | Write API base URL |
| `READ_API_URL` | Yes | `http://localhost:8083` | Read API base URL |

## Circuit Breaker

The proxy includes a circuit breaker per backend service:
- **Threshold**: Configurable consecutive failure count
- **Open state**: Returns 503 immediately without forwarding
- **Recovery**: Half-open after timeout, re-closes on success

## Running

```bash
# Standalone
JWT_SECRET="..." INTERNAL_SECRET="..." WRITE_API_URL="..." READ_API_URL="..." go run ./cmd/server

# With Air (hot-reload) — configured via .air.toml
air
```
