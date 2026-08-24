-- name: CreateBusEvent :one
INSERT INTO outbox_events (
  id,
  event_type,
  aggregate_id,
  payload,
  created_at
) VALUES (
  $1, $2, $3, $4, $5
)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: GetPendingOutBusEvent :many
SELECT *
FROM outbox_events
WHERE published_at IS NULL
ORDER BY created_at
LIMIT $1;

-- name: MarkOutboxEventProcessed :exec
UPDATE outbox_events
SET published_at = NOW()
WHERE id = $1
  AND published_at IS NULL;
