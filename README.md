# intania-openhouse-2026-api

Backend starter for Intania Openhouse 2026 API. Built in Go with a clean layering (handlers -> usecases -> repositories) and SQL-first persistence.

## Quick start

Prerequisites
- Go 1.24+
- Docker + Docker Compose (for local Postgres)
- Air (for hot reload)

Setup
```bash
# 1) Configure env
cp .env.example .env.dev

# 2) Start dependencies (Postgres)
make up-deps

# 3) Setup migrations
make setup

# 4) Run API with hot-reload
make up
```

Run without Air
```bash
go run . --env-file .env.dev serve
```

Migrations (goose)
```bash
# up | down | reset | create
go run . --env-file .env.dev migrate up
```

Docker build/run
```bash
make docker-build
make docker-run
```

## Configuration

Configuration is defined in `pkg/config/config_template.yaml` and can be overridden by environment variables (viper + env replacer). Example:

```bash
APP_ADDRESS=0.0.0.0:8000
APP_ALLOWED_ORIGINS=http://localhost:3000
APP_IS_PRODUCTION=false
DATABASE_DSN=postgres://user:password@localhost:5432/intania-openhouse-2026?sslmode=disable
```

Air loads `.env.dev` automatically via `.air.toml`.

## Project structure

```
cmd/                    # Cobra CLI commands (serve, migrate)
internal/
  handlers/             # Huma handlers (HTTP layer)
  middlewares/          # Middleware definitions
  migrations/           # Goose SQL migrations (embedded)
  models/               # Domain models
  repositories/         # Data access layer (Bun)
  server/               # Server wiring (chi + huma)
  usecases/             # Business logic layer
pkg/
  baserepo/             # Generic repo helpers + transactions
  config/               # Config loading + validation (viper)
  database/             # Postgres connection (Bun)
Dockerfile              # Distroless container build
docker-compose.yaml     # Local Postgres
Makefile                # Dev shortcuts
.air.toml               # Hot-reload config
```

## Libraries (and why)

- Go: backend language (fast, simple deployment)
- chi + huma-chi: HTTP routing + auto OpenAPI docs (frontend can generate client from spec)
- bun: SQL-first ORM for Postgres with query builder and struct mapping
- goose: migration tool, SQL-based with embedded migration files
- cobra: CLI entrypoint (`serve`, `migrate`)
- viper + godotenv: config loading from YAML + env overrides
- validator: config validation
- cors middleware: CORS handling for browser clients

## API docs

Huma provides generated docs and OpenAPI spec when `APP_IS_PRODUCTION=false` (see `internal/server/server.go`). Default paths:
- Docs: `/docs`
- OpenAPI: `/openapi.json`
- Schemas: `/schemas`
