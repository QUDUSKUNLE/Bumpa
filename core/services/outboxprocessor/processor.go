package outboxprocessor

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/QUDUSKUNLE/Bumpa/core/ports"
	"github.com/QUDUSKUNLE/Bumpa/core/utils"
	"github.com/google/uuid"
)

type OutboxProcessor struct {
	repo ports.RepositoryPorts
	bus  events.EventBus
}

func NewOutboxProcessor(repo ports.RepositoryPorts, bus events.EventBus) *OutboxProcessor {
	return &OutboxProcessor{
		repo: repo,
		bus:  bus,
	}
}

func (p *OutboxProcessor) Process(ctx context.Context) error {
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
