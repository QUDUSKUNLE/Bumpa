package repositories

import (
	"context"

  "github.com/google/uuid"
	// "github.com/QUDUSKUNLE/Bumpa/adapters/events"
)



func (r *Repository) UnlockAchievementIfNew(ctx context.Context, userID uuid.UUID, achievementCode string) (bool, error) {
    // insert into user_achievements with ON CONFLICT DO NOTHING
    return true, nil
}
