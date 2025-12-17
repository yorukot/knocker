-- Remove policies (no-op if already absent).
SELECT remove_continuous_aggregate_policy('monitor_hourly_summary');
SELECT remove_retention_policy('pings');

DROP MATERIALIZED VIEW IF EXISTS monitor_hourly_summary;
