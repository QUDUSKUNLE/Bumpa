-- name: CreateUserBadge :one
INSERT INTO user_badges (user_id, badge_code)
VALUES ($1, $2)
ON CONFLICT (user_id, badge_code) DO NOTHING
RETURNING user_id;

