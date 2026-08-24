package achievements

import (
	"context"
	"encoding/json"
	"time"

	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/QUDUSKUNLE/Bumpa/core/ports"
	"github.com/QUDUSKUNLE/Bumpa/core/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

func NewAchievementService(repo ports.RepositoryPorts, defs []events.AchievementDefinition, bus events.EventBus) *AchievementService {
	return &AchievementService{
		repo:                   repo,
		achievementDefinitions: defs,
		bus:                    bus,
	}
}

func (s *AchievementService) ProcessPurchase(
	ctx context.Context,
	userID pgtype.UUID,
	purchase events.Purchase,
) error {
	utils.LogInfo("Triggered ProcessPurchase: %v", nil)
	err := s.repo.WithTx(ctx, func(tx ports.RepositoryPorts) error {
		inserted, err := tx.InsertPurchaseIfNew(ctx, purchase)

		if err != nil {
			utils.LogError("InserPurchaseIfNew Service Error: %v", err)
			return err
		}
		if !inserted {
			return nil
		}
		stats, err := tx.GetPurchaseStats(ctx, userID)
		if err != nil {
			utils.LogError("Get Purchases Stats Service Error: %v", err)
			return err
		}

		for _, definition := range s.achievementDefinitions {
			if !definition.Condition(stats) {
				continue
			}

			unlocked, err := tx.UnlockAchievementIfNew(ctx, userID, definition.Code)
			if err != nil {
				utils.LogError("UnlockAchievementIfNew Service Error: %v", err)
				return err
			}
			if !unlocked {
				continue
			}

			payload := events.AchievementUnlockedPayload{
				AchievementName: definition.Name,
			}

			body, err := json.Marshal(payload)
			if err != nil {
				utils.LogError("JSON Marshal Service Error: %v", err)
				return err
			}

			evt := domain.Event{
				ID:          uuid.New(),
				UserID:      uuid.UUID(userID.Bytes),
				Type:        "AchievementUnlocked",
				OccurredAt:  time.Now().UTC(),
				AggregateID: uuid.UUID(userID.Bytes),
				Payload:     body,
			}

			if err := tx.AddOutboxEvent(ctx, evt); err != nil {
				utils.LogError("AddOutboxEvent Service Error: %v", err)
				return err
			}
		}
		utils.LogInfo("Finished ProcessPurchase: %v", nil)
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *AchievementService) Process(ctx context.Context) error {
	events, err := p.repo.GetPendingOutboxEvents(ctx, 100)
	if err != nil {
		return err
	}

	for _, event := range events {
		evt := domain.Event{
			ID:          uuid.UUID(event.ID.Bytes),
			UserID:      uuid.UUID(event.AggregateID.Bytes),
			Type:        event.EventType,
			AggregateID: uuid.UUID(event.AggregateID.Bytes),
			OccurredAt:  event.CreatedAt.Time,
			Payload:     event.Payload,
		}

		if err := p.bus.Publish(ctx, evt); err != nil {
			utils.LogError(
				"Failed publishing outbox event %s: %v",
				event.ID,
				err,
			)
			continue
		}

		// if err := p.repo.MarkOutboxEventProcessed(
		//     ctx,
		//     event.ID,
		// ); err != nil {
		//     utils.LogError(
		//         "Failed marking outbox event %s as processed: %v",
		//         event.ID,
		//         err,
		//     )
		// }
	}

	return nil
}
