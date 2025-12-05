# API Surface and Auth

Reference for working in `api/`, request validation, and auth controls.

## Stack and layout
- Echo server with routes grouped in `api/router/` (one file per resource). Handlers live in `api/handler/<resource>/` with one action per file (e.g., `create.go`, `list.go`).
- Middleware: auth enforcement in `api/middleware/`; apply `AuthRequiredMiddleware` on groups needing JWTs.
- Responses: use `utils/response` helpers (`response.Success` / `response.Error`) for consistent envelopes.
- Swagger: annotations above handlers feed `api/docs/`; regenerate with `make generate-docs` when request/response schemas change.

## Auth flow
- JWT-based auth; tokens validated by middleware and user ID extracted via `utils/auth.GetUserIDFromContext`.
- OAuth: accounts/users created via Google OAuth helpers in auth handlers; refresh tokens stored in DB; `CreateRefreshToken` and `UpdateRefreshTokenUsedAt` manage lifecycle.
- Team scoping: most routes accept `teamID` path params; handlers verify membership via `GetTeamMemberByUserID` before returning data.

## Handler patterns
- Parse and validate path params early (strconv to int64); return 400 on parse errors.
- Decode JSON bodies to request structs and validate with `validator.v10` tags.
- Start a transaction, perform membership and ownership checks, then business logic. Always defer `DeferRollback` and commit before returning.
- Logging: use `zap.L()` with structured fields; prefer contextual errors over generic logs.
- Public flags: incident event creation respects `public` booleans; default is true.

## Monitor and incident endpoints
- Monitor CRUD under `api/router/monitor.go`; config stored as JSON and validated via model helper methods.
- Incident endpoints under `api/router/incident.go`:
  - Manual creation when no open incident exists; defaults to `detected` status.
  - Status updates map statuses to event types; `resolved` sets `resolved_at`.
  - Event listing/creation are scoped by monitor and incident IDs with membership checks.

## Notifications and routing
- Notification CRUD under `api/router/notification.go`; configs are raw JSON stored in DB and interpreted by `core/notification/*` when dispatching.
- Monitor-to-notification associations managed via `CreateMonitorNotifications`/`DeleteMonitorNotifications`; router ensures monitor belongs to team before linking.

## Error handling and codes
- Use specific HTTP codes: 400 for invalid params/bodies, 401 for missing auth, 404 for missing scoped resources, 409 for conflict (e.g., open incident exists), 500 for unexpected errors.
- Avoid leaking internal errors to clients; log details with zap and return generic messages.
