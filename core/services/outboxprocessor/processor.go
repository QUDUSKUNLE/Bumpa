package outboxprocessor

import (
	"context"
	"fmt"
	"time"

	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/QUDUSKUNLE/Bumpa/core/ports"
	"github.com/QUDUSKUNLE/Bumpa/core/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

// Process publishes pending outbox events.
func (p *OutboxProcessor) OutboxProcessor(ctx context.Context) error {
	eventArray, err := p.repo.GetPendingOutboxEvents(ctx, 100)
	if err != nil {
		return err
	}

	utils.LogInfo(
		"OutboxProcessor: found %d pending events",
		len(eventArray),
	)

	for _, event := range eventArray {
		evt := events.Event{
			ID:          uuid.UUID(event.ID.Bytes),
			UserID:      uuid.UUID(event.AggregateID.Bytes),
			Type:        event.EventType,
			AggregateID: uuid.UUID(event.AggregateID.Bytes),
			OccurredAt:  event.CreatedAt.Time,
			Payload:     event.Payload,
		}

		if err := p.processEvent(ctx, evt); err != nil {
			utils.LogError(
				"Failed processing outbox event %v: %v",
				event.ID,
				err,
			)
			// Don't stop the other events because one failed.
			continue
		}

	}

	return nil
}

func (p *OutboxProcessor) processEvent(
	ctx context.Context,
	event events.Event,
) error {

	evt := domain.Event{
		ID:             event.ID,
		UserID:         event.AggregateID,
		Type:           event.Type,
		AggregateID:    event.AggregateID,
		OccurredAt:     event.OccurredAt,
		Payload:        event.Payload,
		PaymentAccount: event.PaymentAccount,
	}

	utils.LogInfo(
		"Publishing outbox event ID=%s Type=%s AggregateID=%s",
		evt.ID,
		evt.Type,
		evt.AggregateID,
	)

	// Publish FIRST.
	if err := p.bus.Publish(ctx, evt); err != nil {
		utils.LogError(
			"Failed publishing outbox event %s: %v",
			event.ID,
			err,
		)
		return fmt.Errorf("publish event: %w", err)
	}

	utils.LogInfo(
		"Successfully published outbox event ID=%s Type=%s",
		event.ID,
		event.Type,
	)

	id := pgtype.UUID{Bytes: event.ID, Valid: true}

	// Only mark as processed after successful publication.
	if err := p.repo.MarkOutboxEventProcessed(ctx, id); err != nil {
		utils.LogError(
			"Failed marking outbox event %s processed: %v",
			event.ID,
			err,
		)
		return fmt.Errorf("mark event processed: %w", err)
	}

	utils.LogInfo(
		"Successfully processed outbox event ID=%s Type=%s",
		evt.ID,
		evt.Type,
	)

	return nil
}

func (p *OutboxProcessor) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	utils.LogInfo("Outbox processor started")

	// Process immediately on startup.
	p.runOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			utils.LogInfo("Outbox processor stopped")
			return

		case <-ticker.C:
			p.runOnce(ctx)
		}
	}
}

func (p *OutboxProcessor) runOnce(ctx context.Context) {
	if err := p.OutboxProcessor(ctx); err != nil {
		utils.LogError(
			"Outbox processor error: %v",
			err,
		)
	}
}
