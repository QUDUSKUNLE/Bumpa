package achievements

import (
	"context"
	"encoding/json"
	"time"

	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/QUDUSKUNLE/Bumpa/core/ports"
	"github.com/google/uuid"
)

type AchievementService struct {
	repo                   ports.Repository
	achievementDefinitions []events.AchievementDefinition
}

func NewAchievementService(repo ports.Repository, defs []events.AchievementDefinition) *AchievementService {
	return &AchievementService{
		repo:                   repo,
		achievementDefinitions: defs,
	}
}

func (s *AchievementService) ProcessPurchase(
	ctx context.Context,
	userID uuid.UUID,
	purchase events.Purchase,
) error {
	return s.repo.WithTx(ctx, func(tx ports.Repository) error {
		inserted, err := tx.InsertPurchaseIfNew(ctx, purchase)
		if err != nil {
			return err
		}
		if !inserted {
			return nil
		}
		stats, err := tx.GetPurchaseStats(ctx, userID)
		if err != nil {
			return err
		}
		user, err := tx.GetUser(ctx, userID)
		if err != nil {
			return err
		}

		for _, definition := range s.achievementDefinitions {
			if !definition.Condition(stats) {
				continue
			}

			unlocked, err := tx.UnlockAchievementIfNew(ctx, userID, definition.Code)
			if err != nil {
				return err
			}
			if !unlocked {
				continue
			}

			payload := events.AchievementUnlockedPayload{
				AchievementName: definition.Name,
				User:            user,
			}

			body, err := json.Marshal(payload)
			if err != nil {
				return err
			}

			if err := tx.AddOutboxEvent(ctx, domain.Event{
				ID:          uuid.New(),
				Type:        "AchievementUnlocked",
				OccurredAt:  time.Now().UTC(),
				AggregateID: userID,
				Payload:     body,
			}); err != nil {
				return err
			}
		}
		return nil
	})
}
