package repositories

import (
	"context"
	"errors"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/core/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *Repository) GetUnlockedAchievementCount(ctx context.Context, userID pgtype.UUID) (int, error) {
	count, err := r.queries.CountUserAchievement(ctx, userID)

	if err != nil {
		utils.LogError("CountUserAchievement Repositories Error: %v", err)
		return 0, err
	}
	return int(count), nil
}

func (r *Repository) UnlockBadgeIfNew(ctx context.Context, userID pgtype.UUID, badgeCode string) (bool, error) {
	_, err := r.queries.CreateUserBadge(
		ctx,
		db.CreateUserBadgeParams{UserID: userID, BadgeCode: badgeCode})
	if errors.Is(err, pgx.ErrNoRows) {
		// User Badge Already Exist
		return false, nil
	}
	if err != nil {
		utils.LogError("CreateUserBadge Repositories Error: %v", err)
		return false, err
	}
	return true, nil
}
