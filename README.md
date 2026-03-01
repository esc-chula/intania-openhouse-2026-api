# intania-openhouse-2026-api

Backend starter for Intania Openhouse 2026 API. Built in Go with a clean layering (handlers -> usecases -> repositories) and SQL-first persistence.

## Quick start

Prerequisites
- [Go](https://go.dev/) 1.24+
- Docker + Docker Compose (for local Postgres)
- [Air](https://github.com/air-verse/air) (optional, for hot reload)

Setup
```bash
# 1) Download Go libraries
go mod download

# 2) Configure env and copy `service-account.json` file into the root directory
cp .env.example .env.dev

# 3) Start dependencies (Postgres)
make up-deps

# 4) Setup migrations (migrate reset and migrate up)
make setup

# 5) Run API with hot-reload
make up

# Or run without hot-reload
go run . --env-file .env.dev serve
```

Migrations (goose)
```bash
make migrate-up                          # migrate up to the latest version
make migrate-down                        # migrate down one version
make migrate-reset                       # migrate down all versions
make migrate-create ARGS=add_rls_support # create new migration file
```

Docker build/run
```bash
make docker-build
make docker-run
```

## Configuration

Configuration is defined in `pkg/config/config_template.yaml` and can be overridden by environment variables (in `.env.dev`). Example:

```bash
APP_ADDRESS=0.0.0.0:1234
APP_ALLOWED_ORIGINS=https://intania-openhouse-2026.vercel.app
APP_IS_PRODUCTION=true
DATABASE_DSN=postgres://realuser:realpassword@realhost:5432/intania-openhouse-2026?sslmode=disable
```

Air loads `.env.dev` automatically via `.air.toml`.

## Project structure

```
cmd/                    # Cobra CLI commands (serve, migrate)
internal/
  handlers/             # Huma handlers (HTTP layer)
  middlewares/          # Middlewares
  migrations/           # SQL migrations
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
Makefile                # Dev commands
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

## Data access pattern (Executor + Transactioner)

`pkg/baserepo` keeps repositories transaction-agnostic by reading `bun.IDB` from context. Usecases define transaction boundaries via `Transactioner`, and repositories execute queries via `Executor` that automatically uses the transaction if present.

## API docs

Huma provides generated docs and OpenAPI spec when `APP_IS_PRODUCTION=false` (see `internal/server/server.go`). Default paths:
- Docs: `/docs`
- OpenAPI: `/openapi.json`
- Schemas: `/schemas`
