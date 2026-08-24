-- name: InsertUserAchievement :one
INSERT INTO user_achievements (user_id, achievement_code)
VALUES ($1, $2)
ON CONFLICT (user_id, achievement_code) DO NOTHING
RETURNING user_id;

-- name: CountUserAchievement :one
SELECT COUNT(*)
FROM user_achievements
WHERE user_id = $1;

-- name: GetUserAchievements :one
SELECT
    COALESCE(
        ARRAY_AGG(a.name ORDER BY a.achievement_group, a.position)
        FILTER (WHERE ua.user_id IS NOT NULL),
        '{}'
    ) AS unlocked_achievements,

    COALESCE(
        ARRAY_AGG(a.name ORDER BY a.achievement_group, a.position)
        FILTER (
            WHERE ua.user_id IS NULL
              AND NOT EXISTS (
                  SELECT 1
                  FROM achievements previous
                  WHERE previous.achievement_group = a.achievement_group
                    AND previous.position < a.position
                    AND NOT EXISTS (
                        SELECT 1
                        FROM user_achievements previous_ua
                        WHERE previous_ua.user_id = $1
                          AND previous_ua.achievement_code = previous.code
                    )
              )
        ),
        '{}'
    ) AS next_available_achievements

FROM achievements a
LEFT JOIN user_achievements ua
    ON ua.achievement_code = a.code
    AND ua.user_id = $1;
