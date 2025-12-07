# Notifications, Queues, and Messaging

Guide to how notifications are queued, formatted, and sent.

## Queue types and flow
- Task types: `monitor:ping:{region}` for monitor execution and `notification:dispatch` for outbound alerts. Queue names come from task type strings; workers consume only tasks matching their `APP_REGION` for monitor pings.
- Enqueue points: scheduler enqueues monitor ping tasks; incident handling enqueues notification dispatch tasks only when an incident is opened or resolved.
- Asynq config: worker concurrency and queue weights are set in `worker/worker.go` (critical/default/low). Monitor ping handlers are registered per region; notification dispatch handler listens on the default queue.

## Notification dispatch pipeline
1. `HandleNotificationDispatch` (`worker/handler/notification_dispatch.go`) unmarshals payload and loads monitor + notification via repository inside a transaction.
2. Message is built with `core/notification.FormatMessage`, combining monitor name, status, region, latency, timestamp, and optional detail.
3. `core/notification.Send` routes by notification type:
   - Discord: sends webhook payload from `core/notification/discord.go` using `DiscordNotificationConfig` (`webhook_url`).
   - Telegram: uses bot token + chat ID (`TelegramNotificationConfig`) via `core/notification/telegram.go`.
   - Email: placeholder; returns "not implemented" error if used.
4. Errors are logged with zap and stop the task (will be retried by Asynq policy); successful sends log notification metadata.

## Payloads and detail
- NotificationPayload includes `TeamID`, `MonitorID`, `NotificationID`, `Region`, `Ping` snapshot (status/latency/time), and `Detail` string.
- Detail string usually comes from ping execution or incident message. It is trimmed before formatting and appears in the description when present.
- Title format: `<monitor name> is <STATUS>` (status uppercased). Description includes monitor, region, status, latency (if >0), checked time, and optional detail.

## Configuration and safety
- Ensure notification configs are validated on creation (handlers should validate required fields for the chosen type before storing raw JSON).
- Secrets (webhook URLs, bot tokens) must stay in `.env` or DB—never commit them. Example placeholders only.
- When adding new notification channels, implement a send function in `core/notification`, extend `NotificationType` enum in `models/monitor.go`, and wire handling in `Send` plus API validation.
