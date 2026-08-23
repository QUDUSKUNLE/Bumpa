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
