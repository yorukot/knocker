# Repository Guidelines

Use this guide to work consistently in the Knocker Go codebase and ship changes safely. Deep dives: `agents/overview.md` (runtime flow), `agents/monitoring.md` (pings/incidents), `agents/data.md` (schema/repository), `agents/api.md` (HTTP/auth), and `agents/notifications.md` (queues/dispatch).

## Project Structure & Module Organization
- Entry point: `cmd/main.go` starts API, worker, and schedular (or targeted via `api`/`worker`/`schedular` args).
- HTTP: handlers/middleware/routing in `api/`; Swagger in `api/docs/`.
- Data: `db/` connections, `repository/` queries, `models/` structs, `migrations/` SQL (one active incident per monitor enforced in `migrations/2_*`). IDs come from `utils/id` (sonyflake).
- Background: `worker/` (Asynq consumers, ping handling, notifications) and `schedular/` (queues monitor tasks, updates next_check).
- Utilities: `utils/` for config (`utils/config`), logging (`utils/logger`), IDs, helpers. Artifacts in `tmp/`. Local services in `compose.yaml`.

## How the system operates
- Scheduler polls every 2s for monitors whose `next_check` has elapsed, enqueues `monitor:ping:{region}` tasks for each region (`APP_REGIONS`), and updates `last_checked`/`next_check` with jitter. Worker consumes only the current `APP_REGION` queue.
- Worker executes monitors via `core/monitor.Run` (HTTP or ping), buffers ping results to `pings` through `PingRecorder`, then evaluates incidents (`worker/handler/monitor_ping.go`). Notifications dispatch asynchronously when incidents are created or resolved.
- Refer to `agents/monitoring.md` before modifying ping execution, thresholds, incident logic, or notification triggers.

## API Design Conventions
- RESTful resources (e.g., `/users`, `/teams/{id}`) with Auth-required routes in `api/router/`.
- Handlers live under `api/handler/<resource>/`, one handler per file (e.g., `list.go`, `create.go`); define routes per resource file in `api/router/` (e.g., `monitor.go`, `incident.go`).
- Keep DTOs beside their handlers; generated docs live with the API (`make generate-docs`).

## Coding Style & Naming Conventions
- Go defaults: tabs; always commit gofmt output. Follow the Uber Go style guide.
- Keep package-scoped names short and descriptive; place request/response DTOs near handlers.
- Prefix repository methods by resource (e.g., `UserCreate`, `TeamList`); add env vars in `utils/config` with documented defaults.
- Use `utils/logger` for structured logs and pass contexts for request-scoped work.

## Testing Guidelines
- Prefer table-driven tests in the same package; name `TestFunction_Scenario`.
- Cover new public functions plus validation and data-access edges; prefer fakes over real services.
- Run `go test ./...` before pushing; keep fixtures near the tested package.

## Security & Configuration Tips
- `.env` auto-loads; never commit secrets - use example placeholders. Set `JWT_SECRET_KEY`, OAuth, and DB/Redis values per `utils/config`.
- Avoid default passwords locally; when migrations change, document rollback behavior in the PR.
