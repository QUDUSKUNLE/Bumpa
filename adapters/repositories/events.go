package repositories

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *Repository) AddOutboxEvent(ctx context.Context, event domain.Event) error {
	_, err := r.queries.CreateBusEvent(ctx, db.CreateBusEventParams{
		ID:          pgtype.UUID{Bytes: uuid.New(), Valid: true},
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
