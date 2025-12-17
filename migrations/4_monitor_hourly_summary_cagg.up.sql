-- Ensure the pings table is a hypertable on the time column so continuous aggregates work.
SELECT create_hypertable('pings', 'time', if_not_exists => TRUE);

-- Hourly rollup for availability and latency percentiles.
CREATE MATERIALIZED VIEW monitor_hourly_summary
WITH (timescaledb.continuous) AS
SELECT
    monitor_id,
    region,
    time_bucket('1 hour', time) AS bucket,
    count(*) AS total_count,
    count(*) FILTER (
        WHERE status = 'successful' AND latency <= 5000
    ) AS good_count,
    percentile_cont(0.50) WITHIN GROUP (ORDER BY latency / 1000.0) AS p50,
    percentile_cont(0.75) WITHIN GROUP (ORDER BY latency / 1000.0) AS p75,
    percentile_cont(0.90) WITHIN GROUP (ORDER BY latency / 1000.0) AS p90,
    percentile_cont(0.95) WITHIN GROUP (ORDER BY latency / 1000.0) AS p95,
    percentile_cont(0.99) WITHIN GROUP (ORDER BY latency / 1000.0) AS p99
FROM pings
GROUP BY monitor_id, region, bucket
WITH NO DATA;

-- Keep the materialized view fresh: refresh last 24 hours every 15 minutes,
-- skipping the most recent hour to avoid hot buckets.
SELECT add_continuous_aggregate_policy(
    'monitor_hourly_summary',
    start_offset => INTERVAL '24 hours',
    end_offset   => INTERVAL '1 hour',
    schedule_interval => INTERVAL '15 minutes'
);

-- Retention: drop raw ping data older than 90 days.
SELECT add_retention_policy('pings', INTERVAL '90 days');
