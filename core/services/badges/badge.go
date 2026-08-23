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
)

type BadgeService struct {
	repo             ports.RepositoryPorts
	badgeDefinitions []events.BadgeDefinition
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

func NewBadgeService(repo ports.RepositoryPorts, defs []events.BadgeDefinition) *BadgeService {
	return &BadgeService{
		repo:             repo,
		badgeDefinitions: defs,
	}
}

func (s *BadgeService) HandleAchievementUnlocked(
	ctx context.Context,
	event events.Event,
) error {
	var payload events.AchievementUnlockedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}
	return s.repo.WithTx(ctx, func(tx ports.RepositoryPorts) error {
		count, err := tx.GetUnlockedAchievementCount(ctx, event.AggregateID)
		if err != nil {
			utils.LogError("GetUnlockedAchievementCount error", err)
			return err
		}
		user, err := tx.GetUser(ctx, event.UserID)
		if err != nil {
			utils.LogError("GetUser error", err)
			return err
		}
		for _, badge := range s.badgeDefinitions {
			if count < badge.RequiredRewards {
				continue
			}

			unlocked, err := tx.UnlockBadgeIfNew(ctx, event.UserID, badge.Code)
			if err != nil {
				utils.LogError("UnlockBadgeIfNew error", err)
				return err
			}
			if !unlocked {
				continue
			}

			out := events.BadgeUnlockedPayload{
				BadgeName: badge.Name,
				BadgeCode: badge.Code,
				User:      user.ID.String(),
			}

			ev := domain.Event{
				ID:          uuid.New(),
				UserID:      event.UserID,
				Type:        "BadgeUnlocked",
				OccurredAt:  time.Now().UTC(),
				AggregateID: event.AggregateID,
				Payload:     utils.MustJSON(out),
			}

			if err := tx.AddOutboxEvent(ctx, ev); err != nil {
				utils.LogError("AddOutboxEvent error", err)
				return err
			}
		}
		return nil
	})
}
