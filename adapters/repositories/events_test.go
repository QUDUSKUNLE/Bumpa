package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_AddOutboxEvent(t *testing.T) {
	ctx := context.Background()
	repo := setupTestRepository(t)

	eventID := uuid.New()
	aggregateID := uuid.New()
	occurredAt := time.Now().UTC()

	event := &domain.Event{
		ID:          eventID,
		Type:        "achievement.unlocked",
		AggregateID: aggregateID,
		Payload:     []byte(`{"achievement":"first_purchase"}`),
		OccurredAt:  occurredAt,
	}

	err := repo.AddOutboxEvent(ctx, *event)

	require.NoError(t, err)
}

func TestRepository_GetPendingOutboxEvents(t *testing.T) {
	ctx := context.Background()
	repo := setupTestRepository(t)

	eventID := uuid.New()

	event := &domain.Event{
		ID:          eventID,
		Type:        "achievement.unlocked",
		AggregateID: uuid.New(),
		Payload:     []byte(`{"achievement":"first_purchase"}`),
		OccurredAt:  time.Now().UTC(),
	}

	err := repo.AddOutboxEvent(ctx, *event)
	require.NoError(t, err)

	events, err := repo.GetPendingOutboxEvents(ctx, 1000)

	require.NoError(t, err)

	var found bool

	for _, pendingEvent := range events {
		if pendingEvent.ID.Bytes == eventID && pendingEvent.ID.Valid {
			found = true
			break
		}
	}

	assert.True(t, found, "created event should be returned as a pending event")
}

func TestRepository_MarkOutboxEventProcessed(t *testing.T) {
	ctx := context.Background()
	repo := setupTestRepository(t)

	eventID := uuid.New()

	event := &domain.Event{
		ID:          eventID,
		Type:        "achievement.unlocked",
		AggregateID: uuid.New(),
		Payload:     []byte(`{"achievement":"first_purchase"}`),
		OccurredAt:  time.Now().UTC(),
	}

	err := repo.AddOutboxEvent(ctx, *event)
	require.NoError(t, err)

	eventUUID := pgtype.UUID{
		Bytes: eventID,
		Valid: true,
	}

	err = repo.MarkOutboxEventProcessed(ctx, eventUUID)

	require.NoError(t, err)

	events, err := repo.GetPendingOutboxEvents(ctx, 10)

	require.NoError(t, err)

	for _, pendingEvent := range events {
		assert.NotEqual(
			t,
			eventUUID,
			pendingEvent.ID,
			"processed event should not remain pending",
		)
	}
}
