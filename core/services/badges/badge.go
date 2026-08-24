package badges

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

type BadgeService struct {
	repo             ports.RepositoryPorts
	badgeDefinitions []events.BadgeDefinition
	bus              events.EventBus
}

func BadgeDefinition() []events.BadgeDefinition {
	return []events.BadgeDefinition{
		{
			Code:            "bronze",
			Name:            "Bronze",
			RequiredRewards: 1,
		},
		{
			Code:            "silver",
			Name:            "Silver",
			RequiredRewards: 3,
		},
		{
			Code:            "gold",
			Name:            "Gold",
			RequiredRewards: 5,
		},
		{
			Code:            "platinum",
			Name:            "Platinum",
			RequiredRewards: 10,
		},
		{
			Code:            "diamond",
			Name:            "Diamond",
			RequiredRewards: 20,
		},
	}
}

func NewBadgeService(repo ports.RepositoryPorts, defs []events.BadgeDefinition, bus events.EventBus) *BadgeService {
	return &BadgeService{
		repo:             repo,
		badgeDefinitions: defs,
		bus:              bus,
	}
}

func (s *BadgeService) HandleAchievementUnlocked(
	ctx context.Context,
	event events.Event,
) error {
	utils.LogInfo("Triggered HandleAchievementUnlocked: %v", nil)
	utils.LogInfo(
		"AchievementUnlocked Event: ID=%v UserID=%v AggregateID=%v Type=%v",
		event.ID,
		event.UserID,
		event.AggregateID,
		event.Type,
	)
	var payload events.AchievementUnlockedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}

	err := s.repo.WithTx(ctx, func(tx ports.RepositoryPorts) error {

		// 1. Get the user's total number of purchases.
		stats, err := tx.GetPurchaseStats(
			ctx,
			pgtype.UUID{Bytes: event.AggregateID, Valid: true},
		)
		utils.LogInfo(
			"Achievement count for user %s = %d",
			event.UserID,
			stats.TotalPurchases,
		)

		purchaseCount := stats.TotalPurchases

		if err != nil {
			utils.LogError("GetPurchaseStats Service Error: %v", err)
			return err
		}

		for _, badge := range s.badgeDefinitions {
			utils.LogInfo(
				"Checking badge=%s required=%d count=%d",
				badge.Code,
				badge.RequiredRewards,
				purchaseCount,
			)
			if purchaseCount < badge.RequiredRewards {
				continue
			}

			// 3. Unlock only if this badge hasn't already been unlocked.
			unlocked, err := tx.UnlockBadgeIfNew(
				ctx,
				pgtype.UUID{Bytes: event.UserID, Valid: true},
				badge.Code,
			)
			if err != nil {
				utils.LogError(
					"UnlockBadgeIfNew Service Error: %v",
					err,
				)
				return err
			}

			// Badge was already unlocked.
			if !unlocked {
				utils.LogInfo(
					"Badge already unlocked: %s",
					badge.Code,
				)
				continue
			}

			utils.LogInfo(
				"NEW BADGE UNLOCKED: %s (%s)",
				badge.Code,
				badge.Name,
			)

			out := events.BadgeUnlockedPayload{
				BadgeName: badge.Name,
				BadgeCode: badge.Code,
				User:      pgtype.UUID{Bytes: event.UserID, Valid: true},
			}

			evt := domain.Event{
				ID:             uuid.New(),
				UserID:         event.UserID,
				Type:           "BadgeUnlocked",
				OccurredAt:     time.Now().UTC(),
				AggregateID:    event.AggregateID,
				Payload:        utils.MustJSON(out),
				PaymentAccount: event.PaymentAccount,
			}

			if err := tx.AddOutboxEvent(ctx, evt); err != nil {
				utils.LogError("AddOutboxEvent Service Error: %v", err)
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	utils.LogInfo(
		"Finished HandleAchievementUnlocked: user=%s",
		event.UserID,
	)
	return nil
}
