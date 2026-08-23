package repositories

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/core/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *Repository) UnlockAchievementIfNew(ctx context.Context, userID uuid.UUID, achievementCode string) (bool, error) {
	_, err := r.queries.InsertUserAchievement(ctx, db.InsertUserAchievementParams{
		UserID:          pgtype.UUID{Bytes: userID, Valid: true},
		AchievementCode: achievementCode,
	})
	if err != nil {
		utils.LogError("InserUserAchievement error", err)
		return false, err
	}
	return true, nil
}
