-- name: MessageSummaryMysql :one
SELECT COUNT(*)                                                                 AS total_messages,

       COUNT(CASE WHEN status = 0 THEN 1 END)                                   AS success_messages,
       COUNT(CASE WHEN status = 2 THEN 1 END)                                   AS failed_messages,

       CAST(COUNT(CASE WHEN status = 0 THEN 1 END) * 100.0 / COUNT(*) AS FLOAT) AS success_rate,
       CAST(COUNT(CASE WHEN status = 2 THEN 1 END) * 100.0 / COUNT(*) AS FLOAT) AS error_rate,

       AVG(TIMESTAMPDIFF(SECOND, created_at, sent_at))                          AS avg_send_time
FROM messages
WHERE workspace_id = ?;

-- name: GetProviderStatsMysql :one
SELECT COUNT(*)                                                                 AS total,
       COUNT(CASE WHEN status = 0 THEN 1 END)                                   AS success,
       COUNT(CASE WHEN status = 2 THEN 1 END)                                   AS failed,
       CAST(COUNT(CASE WHEN status = 0 THEN 1 END) * 100.0 / COUNT(*) AS FLOAT) AS success_rate,
       CAST(COUNT(CASE WHEN status = 2 THEN 1 END) * 100.0 / COUNT(*) AS FLOAT) AS error_rate,
       AVG(TIMESTAMPDIFF(SECOND, created_at, sent_at))                          AS avg_send_time
FROM messages
WHERE workspace_id = ?
  AND provider_id = ?
ORDER BY total;

-- name: GetAvgSendTimeMysql :one
SELECT COUNT(*)                                                                 AS total,
       COUNT(CASE WHEN status = 0 THEN 1 END)                                   AS success,
       COUNT(CASE WHEN status = 2 THEN 1 END)                                   AS failed,
       CAST(COUNT(CASE WHEN status = 0 THEN 1 END) * 100.0 / COUNT(*) AS FLOAT) AS success_rate,
       CAST(COUNT(CASE WHEN status = 2 THEN 1 END) * 100.0 / COUNT(*) AS FLOAT) AS error_rate,
       AVG(TIMESTAMPDIFF(SECOND, created_at, sent_at))                          AS avg_send_time
FROM messages
WHERE workspace_id = ?
  AND provider_id = ?
ORDER BY total;

-- name: GetTimeDistributionMysql :one
SELECT CASE
           WHEN TIMESTAMPDIFF(DAY, sqlc.narg('start_date'), sqlc.narg('end_date')) > 365
               THEN DATE_FORMAT(created_at, '%Y-%m')
           WHEN TIMESTAMPDIFF(DAY, sqlc.narg('start_date'), sqlc.narg('end_date')) > 30
               THEN DATE_FORMAT(created_at, '%Y-%m-%d')
           ELSE DATE_FORMAT(created_at, '%Y-%m-%d %H:00:00')
           END  AS time_bucket,

       COUNT(*) AS message_count
FROM messages
WHERE created_at BETWEEN sqlc.narg('start_date') AND sqlc.narg('end_date')
  AND workspace_id = ?
  AND provider_id = ?
GROUP BY time_bucket
ORDER BY time_bucket;