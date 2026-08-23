-- name: InsertUserAchievement :one
INSERT INTO user_achievements (user_id, achievement_code)
VALUES ($1, $2)
ON CONFLICT (user_id, achievement_code) DO NOTHING
RETURNING user_id;

-- name: CountUserAchievement :one
SELECT COUNT(*)
FROM user_achievements
WHERE user_id = $1;
