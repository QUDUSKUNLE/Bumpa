-- name: CreatePurchase :one
INSERT INTO purchases (
  id,
  user_id,
  external_id,
  amount_kobo
) VALUES (
  $1, $2, $3, $4
)
ON CONFLICT (id) DO NOTHING
RETURNING *;

-- name: CountPurchasesByUser :one
SELECT COUNT(*)::int AS total_purchases
FROM purchases
WHERE user_id = $1;
