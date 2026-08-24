package repositories

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *Repository) AddOutboxEvent(ctx context.Context, event domain.Event) error {
	_, err := r.queries.CreateBusEvent(ctx, db.CreateBusEventParams{
		ID:          pgtype.UUID{Bytes: event.ID, Valid: true},
		EventType:   event.Type,
		AggregateID: pgtype.UUID{Bytes: event.AggregateID, Valid: true},
		Payload:     event.Payload,
		CreatedAt:   pgtype.Timestamptz{Time: event.OccurredAt, Valid: true},
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetPendingOutboxEvents(ctx context.Context, batchSize int) ([]db.OutboxEvent, error) {
	return r.queries.GetPendingOutBusEvent(ctx, int32(batchSize))
}

func (r *Repository) MarkOutboxEventProcessed(ctx context.Context, id pgtype.UUID) error {
	return r.queries.MarkOutboxEventProcessed(ctx, id)
}
