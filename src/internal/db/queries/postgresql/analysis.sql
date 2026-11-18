-- name: MessageSummaryPostgresql :one
SELECT COUNT(*) AS total_messages,
       COUNT(*) FILTER (WHERE status = 0) AS success_messages,
       COUNT(*) FILTER (WHERE status = 2) AS failed_messages,
       COALESCE(COUNT(*) FILTER (WHERE status = 0) * 100.0 / NULLIF(COUNT(*), 0), 0)::float8 AS success_rate,
       COALESCE(COUNT(*) FILTER (WHERE status = 2) * 100.0 / NULLIF(COUNT(*), 0), 0)::float8 AS error_rate,
       COALESCE(AVG(EXTRACT(EPOCH FROM sent_at - created_at)), 0)::float8 AS avg_send_time
FROM messages
WHERE workspace_id = $1;

-- name: GetProviderStatsPostgresql :one
SELECT COUNT(*) AS total,
       COUNT(*) FILTER (WHERE status = 0) AS success,
       COUNT(*) FILTER (WHERE status = 2) AS failed,
       COALESCE(COUNT(*) FILTER (WHERE status = 0) * 100.0 / NULLIF(COUNT(*), 0), 0)::float8 AS success_rate,
       COALESCE(COUNT(*) FILTER (WHERE status = 2) * 100.0 / NULLIF(COUNT(*), 0), 0)::float8 AS error_rate,
       COALESCE(AVG(EXTRACT(EPOCH FROM sent_at - created_at)), 0)::float8 AS avg_send_time
FROM messages
WHERE workspace_id = $1
  AND provider_id = $2;

-- name: GetTimeDistributionPostgresql :one
WITH bounds AS (
    SELECT sqlc.narg('start_date')::timestamp AS start_date,
           sqlc.narg('end_date')::timestamp AS end_date
)
SELECT date_trunc(
           CASE
               WHEN bounds.end_date - bounds.start_date > INTERVAL '365 days' THEN 'month'
               WHEN bounds.end_date - bounds.start_date > INTERVAL '30 days' THEN 'day'
               ELSE 'hour'
           END,
           messages.created_at
       )::text AS time_bucket,
       COUNT(*) AS message_count
FROM messages CROSS JOIN bounds
WHERE messages.created_at BETWEEN bounds.start_date AND bounds.end_date
  AND messages.workspace_id = sqlc.arg('workspace_id')
  AND messages.provider_id = sqlc.arg('provider_id')
GROUP BY time_bucket
ORDER BY time_bucket;
