package badges

import (
	"context"
	"encoding/json"
	"time"

	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/QUDUSKUNLE/Bumpa/core/ports"
	"github.com/QUDUSKUNLE/Bumpa/core/services"
	"github.com/google/uuid"
)

type BadgeService struct {
	repo             ports.Repository
	badgeDefinitions []events.BadgeDefinition
}

func NewBadgeService(repo ports.Repository, defs []events.BadgeDefinition) *BadgeService {
	return &BadgeService{
		repo:             repo,
		badgeDefinitions: defs,
	}
}

func (s *BadgeService) HandleAchievementUnlocked(
	ctx context.Context,
	event domain.Event,
) error {
	var payload events.AchievementUnlockedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}
	return s.repo.WithTx(ctx, func(tx ports.Repository) error {
		count, err := tx.GetUnlockedAchievementCount(ctx, event.AggregateID)
		if err != nil {
			return err
		}
		user, err := tx.GetUser(ctx, event.AggregateID)
		if err != nil {
			return err
		}
		for _, badge := range s.badgeDefinitions {
			if count < badge.RequiredRewards {
				continue
			}

			unlocked, err := tx.UnlockBadgeIfNew(ctx, event.AggregateID, badge.Code)
			if err != nil {
				return err
			}
			if !unlocked {
				continue
			}

			out := events.BadgeUnlockedPayload{
				BadgeName: badge.Name,
				BadgeCode: badge.Code,
				User:      user,
			}

			ev := domain.Event{
				ID:          uuid.New(),
				Type:        "BadgeUnlocked",
				OccurredAt:  time.Now().UTC(),
				AggregateID: event.AggregateID,
				Payload:     services.MustJSON(out),
			}

			if err := tx.AddOutboxEvent(ctx, ev); err != nil {
				return err
			}
		}
		return nil
	})
}
