package services

import (
	"encoding/json"

	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/ports"
	"github.com/QUDUSKUNLE/Bumpa/core/services/achievements"
	"github.com/QUDUSKUNLE/Bumpa/core/services/badges"
)

type ServicesHandler struct {
	ports              ports.RepositoryPorts
	AchievementService achievements.AchievementService
	BadgeService       badges.BadgeService
	// OutboxProcessorService outboxprocessor.OutboxProcessor
}

func NewServiceAdapter(repositoryPort ports.RepositoryPorts, bus events.EventBus) *ServicesHandler {
	return &ServicesHandler{
		ports: repositoryPort,
		AchievementService: *achievements.NewAchievementService(
			repositoryPort,
			achievements.AchievementDefinition(),
			bus,
		),
		BadgeService: *badges.NewBadgeService(
			repositoryPort,
			badges.BadgeDefinition(),
			bus,
		),
		// OutboxProcessorService: *outboxprocessor.NewOutboxProcessor(repositoryPort, bus),
	}
}

func MustJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
