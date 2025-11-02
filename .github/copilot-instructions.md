<!--
Guidance for AI coding agents working on the nexivent-Backend repository.
Keep this concise and actionable. Update when project layout or conventions change.
-->
# Copilot / AI assistant instructions for nexivent-Backend

This project is a small Go backend (net/http + database/sql + pgx) for the Nexivent ticketing service. Below are the essential facts an AI agent needs to be immediately productive.

- Service entrypoint: `cmd/api/main.go` — constructs an `application` struct, loads config via `loadConfig`, and starts the HTTP server on `ADDR` (default `:4000`). Routes are assembled in `application.routes()`.
- Configuration: `cmd/api/config.go` — environment-driven. Uses `github.com/joho/godotenv` to optionally load a `.env` file. Key env vars: `ADDR`, `DB_DSN`, `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_MAX_IDLE_TIME`.
- Database: `internal/data/postgres.go` exposes `OpenDB(DBConfig) (*sql.DB, error)` which pings the DB. Uses the `pgx` driver (imported for `database/sql`). DB configuration object is `data.DBConfig`.
- SQL patterns / helpers: `internal/data/InsertGenerico.go` defines reusable helpers:
  - `InsertReturningID(ctx, db, table, cols, vals, idcol, dest)` — builds an `INSERT ... RETURNING id` using `$1` placeholders.
  - `UpdateByID(ctx, db, table, cols, vals, idcol, id)` — builds `UPDATE ... SET ... WHERE idcol=$N`.
  - `ExecerQueryer` is the interface used by these helpers; almost any `*sql.DB` or `*sql.Tx` that implements `ExecContext` and `QueryRowContext` will work.
- Error sentinel: `internal/data/errors.go` defines `ErrNotFound` — code should translate this to 404 when building HTTP handlers.

Architecture / conventions (how code is organized and why):
- Minimal layered structure follows common patterns:
  - `cmd/api` — process entrypoint, routing/middlewares, config. Keep only orchestration here.
  - `internal/app` and `internal/*` — application logic and infra. `internal` packages are consumed only within the repo.
  - `internal/data` — database adapters/repos. Repos accept a `context.Context` and `*sql.DB` and return domain types.
  - `internal/domain` — simple domain structs (e.g. `Categoria`, `Ticket`, `Usuario`).
- Handlers/middleware: `application` methods (see `cmd/api/middleware.go`, `cmd/api/health.go`) — use `application` receiver so handlers can access logger and config.

Typical code patterns to follow (concrete examples):
- Use `context.Context` in all DB calls (repo methods accept `cont context.Context`). See `internal/data/categoria_repo.go` for Save/GetById/Delete examples.
- Use the DB helpers rather than hand-building placeholder sequences. Example: `InsertReturningID(ctx, db, "categorias", cols, vals, "id_categoria", &c.IDCategoria)`.
- Treat `ErrNotFound` specially: map it to HTTP 404 in handlers.
- SQL placeholders: Postgres-style `$1`, `$2`, ... (helper functions already generate these).

Run/build/test guidance (what works locally):
- Build/run service: `go run ./cmd/api` or use `go build ./cmd/api` then run the binary. Service listens on `ADDR` env var (default `:4000`).
- DB: default `DB_DSN` in `config.go` points to `postgres://postgres:postgres@localhost:5432/nexivent?sslmode=disable`. Use a local Postgres or Docker compose (this repo has `docker-compose.yml`).
- DB migrations: SQL files live in `migrations/` (e.g. `001_create_categorias.sql`). Keep migrations in sync with repo implementations.
- Tests: there are currently no automated tests in the repo. Before adding tests, follow repository patterns: use `context.Context`, create test DB or use mocking for `ExecerQueryer`.

Integration & dependencies to watch for:
- External: `github.com/jackc/pgx/v5/stdlib` (pgx driver), `github.com/joho/godotenv`.
- Keep SQL error wrapping in mind: repository returns native errors (e.g., sql.ErrNoRows converted to `ErrNotFound`). Handlers must inspect these.

What an AI agent should NOT change without human confirmation:
- Large refactors of package layout or public API. The repo is small and changes can affect orchestration in `cmd/api`.
- DB connection defaults and migration files—confirm with the team before renaming columns or changing migration order.

When adding new code, prefer:
- Additive, small commits that modify a single area (config, repo, handler).
- Use existing helper functions in `internal/data` for consistent SQL generation.

Files to inspect when you need context (quick jump list):
- `cmd/api/main.go`, `cmd/api/config.go`, `cmd/api/middleware.go`, `cmd/api/health.go`
- `internal/data/postgres.go`, `internal/data/InsertGenerico.go`, `internal/data/errors.go`, `internal/data/categoria_repo.go`
- `internal/domain/*.go` for data shapes

If anything above is unclear or you'd like the agent to include extra examples (more handlers, repo patterns, or a suggested test harness), reply with what you want expanded and I'll iterate.
