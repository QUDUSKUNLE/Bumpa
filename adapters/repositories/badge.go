package repositories

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/core/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *Repository) GetUnlockedAchievementCount(ctx context.Context, userID uuid.UUID) (int, error) {
	count, err := r.queries.CountUserAchievement(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		utils.LogError("CountUserAchievement error", err)
		return 0, err
	}
	return int(count), nil
}

func (r *Repository) UnlockBadgeIfNew(ctx context.Context, userID uuid.UUID, badgeCode string) (bool, error) {
	_, err := r.queries.CreateUserBadge(ctx, db.CreateUserBadgeParams{UserID: pgtype.UUID{Bytes: userID, Valid: true}, BadgeCode: badgeCode})
	if err != nil {
		utils.LogError("CreateUserBadge error", err)
		return false, err
	}
	return true, nil
}
