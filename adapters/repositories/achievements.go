package repositories

import (
	"context"
	"errors"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/core/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *Repository) UnlockAchievementIfNew(
	ctx context.Context,
	userID pgtype.UUID,
	achievementCode string,
) (bool, error) {
	_, err := r.queries.InsertUserAchievement(ctx, db.InsertUserAchievementParams{
		UserID:          userID,
		AchievementCode: achievementCode,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		// Achievement Already Exist
		return false, nil
	}
	if err != nil {
		utils.LogError(
			"InsertUserAchievement Repository Error: %v",
			err,
		)
		return false, err
	}
	return true, nil
}

func (r *Repository) GetUserAchievements(
	ctx context.Context,
	userID pgtype.UUID,
) (db.GetUserAchievementsRow, error) {
	return r.queries.GetUserAchievements(ctx, userID)
}
