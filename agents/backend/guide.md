# Backend Agent Guide

- Scope: Go services only (API, worker, schedular, repositories, migrations, utils). Hand off any `website/` work to the frontend track unless backend support is explicitly required for UI.
- Startup: `cmd/main.go` runs API, worker, and schedular; use `api`/`worker`/`schedular` args for targeted runs.
- HTTP: handlers/middleware live under `api/`; routes under `api/router/`; keep DTOs beside handlers and regenerate Swagger when routes change.
- Data: use `repository/` for queries, `models/` for structs, and `migrations/` for schema (sonyflake IDs via `utils/id`). One active incident per monitor enforced in `migrations/2_*`.
- Background: scheduler enqueues `monitor:ping:{region}` tasks per `APP_REGIONS`; worker consumes the current `APP_REGION` queue and evaluates incidents—see `agents/backend/monitoring.md` before touching ping/incident logic.
- Config/logging: add env vars in `utils/config` with defaults; use `utils/logger` with contexts for structured logs.
- Testing: prefer table-driven tests alongside code; cover new public funcs plus validation/data edges; run `go test ./...` before shipping.
