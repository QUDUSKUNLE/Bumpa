-- name: CreateUser :one
INSERT INTO users (
    name,
    email,
    phone,
    payment_account
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, name, email, phone, payment_account, created_at;

-- name: GetUser :one
SELECT id, name, email, phone, payment_account, created_at
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, name, email, phone, payment_account, created_at
FROM users
WHERE email = $1
LIMIT 1;

-- name: ListUsers :many
SELECT id, name, email, phone, payment_account, created_at
FROM users
ORDER BY created_at DESC;

-- name: UpdateUserPaymentAccount :one
UPDATE users
SET payment_account = $2
WHERE id = $1
RETURNING id, name, email, phone, payment_account, created_at;
