package achievements

import (
	"context"

	"github.com/google/uuid"
	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/ports"
)

type AchievementService struct {
	repo                   ports.Repository
	achievementDefinitions []events.AchievementDefinition
}

func (service *AchievementService) ProcessPurchase(
	ctx context.Context,
	userID uuid.UUID,
	purchase events.Purchase,
) error {
	return nil
}
