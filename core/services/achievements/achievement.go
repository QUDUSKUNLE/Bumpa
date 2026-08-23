package achievements

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/QUDUSKUNLE/Bumpa/core/ports"
	"github.com/QUDUSKUNLE/Bumpa/core/utils"
	"github.com/google/uuid"
)

type AchievementService struct {
	bus                    events.EventBus
	repo                   ports.RepositoryPorts
	achievementDefinitions []events.AchievementDefinition
}

func AchievementDefinition() []events.AchievementDefinition {
	return []events.AchievementDefinition{
		{
			Code:     "first_purchase",
			Name:     "First Purchase",
			Group:    "shopping",
			Position: 1,
			Condition: func(stats events.PurchaseStats) bool {
				return stats.TotalPurchases >= 1
			},
		},
		{
			Code:     "three_purchases",
			Name:     "Three Purchases",
			Group:    "shopping",
			Position: 2,
			Condition: func(stats events.PurchaseStats) bool {
				return stats.TotalPurchases >= 3
			},
		},
		{
			Code:     "five_purchases",
			Name:     "Five Purchases",
			Group:    "shopping",
			Position: 3,
			Condition: func(stats events.PurchaseStats) bool {
				return stats.TotalPurchases >= 5
			},
		},
		{
			Code:     "ten_purchases",
			Name:     "Ten Purchases",
			Group:    "shopping",
			Position: 4,
			Condition: func(stats events.PurchaseStats) bool {
				return stats.TotalPurchases >= 10
			},
		},
		{
			Code:     "twenty_purchases",
			Name:     "Twenty Purchases",
			Group:    "shopping",
			Position: 5,
			Condition: func(stats events.PurchaseStats) bool {
				return stats.TotalPurchases >= 20
			},
		},
	}
}

func NewAchievementService(repo ports.RepositoryPorts, defs []events.AchievementDefinition) *AchievementService {
	return &AchievementService{
		repo:                   repo,
		achievementDefinitions: defs,
		bus:                    *events.NewEventBus(),
	}
}

func (s *AchievementService) ProcessPurchase(
	ctx context.Context,
	userID uuid.UUID,
	purchase events.Purchase,
) error {
	return s.repo.WithTx(ctx, func(tx ports.RepositoryPorts) error {
		inserted, err := tx.InsertPurchaseIfNew(ctx, purchase)
		if err != nil {
			utils.LogError("Log InserPurchaseIfNew", err)
			return err
		}
		if !inserted {
			return nil
		}
		stats, err := tx.GetPurchaseStats(ctx, userID)
		if err != nil {
			utils.LogError("Get purchases stats", err)
			return err
		}
		user, err := tx.GetUser(ctx, userID)
		if err != nil {
			utils.LogError("Get user error", err)
			return err
		}

		for _, definition := range s.achievementDefinitions {
			if !definition.Condition(stats) {
				continue
			}

			unlocked, err := tx.UnlockAchievementIfNew(ctx, userID, definition.Code)
			if err != nil {
				utils.LogError("UnlockAchievementIfNew error", err)
				return err
			}
			if !unlocked {
				continue
			}

			payload := events.AchievementUnlockedPayload{
				AchievementName: definition.Name,
				User:            user.ID.String(),
			}

			body, err := json.Marshal(payload)
			if err != nil {
				utils.LogError("JSON Marshal error", err)
				return err
			}

			evt := domain.Event{
				ID:          uuid.New(),
				Type:        "AchievementUnlocked",
				OccurredAt:  time.Now().UTC(),
				AggregateID: userID,
				Payload:     body,
			}

			if err := s.bus.Publish(ctx, evt); err != nil {
				fmt.Println("error publishing event")
				return err
			}

			if err := tx.AddOutboxEvent(ctx, evt); err != nil {
				utils.LogError("Error adding eventOutbox", err)
				return err
			}
		}
		return nil
	})
}
