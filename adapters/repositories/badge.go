package repositories

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) GetUnlockedAchievementCount(ctx context.Context, userID uuid.UUID) (int, error) {
	return 1, nil
}

func (r *Repository) UnlockBadgeIfNew(ctx context.Context, userID uuid.UUID, badgeCode string) (bool, error) {
	return true, nil
}
