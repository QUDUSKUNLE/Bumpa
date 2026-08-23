-- name: CreatePayment :one
INSERT INTO payments (
  id,
  user_id,
  badge_code,
  amount_kobo,
  status,
  created_at
) VALUES (
  $1, $2, $3, $4, $5, $6
)
ON CONFLICT (user_id, badge_code) DO NOTHING
RETURNING id;

-- name: MarkPaymentSuccessful :exec
UPDATE payments
SET status = 'successful',
  provider_reference = $2
WHERE id = $1;

-- name: MarkPaymentFailed :exec
UPDATE payments
SET status = 'failed',
  provider_reference = $2
WHERE id = $1;
