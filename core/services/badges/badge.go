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
		count, err := tx.GetUnlockedAchievementCount(ctx, pgtype.UUID{Bytes: event.AggregateID, Valid: true})
		if err != nil {
			utils.LogError("GetUnlockedAchievementCount Service Error: %v", err)
			return err
		}
		user, err := tx.GetUser(ctx, pgtype.UUID{Bytes: event.UserID, Valid: true})
		if err != nil {
			utils.LogError("GetUser Service Error: %v", err)
			return err
		}
		for _, badge := range s.badgeDefinitions {
			if count < badge.RequiredRewards {
				continue
			}

			unlocked, err := tx.UnlockBadgeIfNew(ctx, pgtype.UUID{Bytes: event.UserID, Valid: true}, badge.Code)
			if err != nil {
				utils.LogError("UnlockBadgeIfNew Service Error: %v", err)
				return err
			}
			if !unlocked {
				continue
			}

			out := events.BadgeUnlockedPayload{
				BadgeName: badge.Name,
				BadgeCode: badge.Code,
				User:      user.ID,
			}

			evt := domain.Event{
				ID:          uuid.New(),
				UserID:      event.UserID,
				Type:        "BadgeUnlocked",
				OccurredAt:  time.Now().UTC(),
				AggregateID: event.AggregateID,
				Payload:     utils.MustJSON(out),
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
	return nil
}
