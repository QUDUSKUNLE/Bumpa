package outboxprocessor

import (
	"context"
	"time"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/QUDUSKUNLE/Bumpa/core/ports"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type mockOutboxRepository struct {
	getPendingFn func(
		context.Context,
		*int,
	) ([]db.OutboxEvent, error)

	markProcessedFn func(
		context.Context,
		pgtype.UUID,
	) error

	getPendingCalls    int
	markProcessedCalls int
	markedIDs          []pgtype.UUID
}

func (m *mockOutboxRepository) WithTx(
	ctx context.Context,
	fn func(ports.RepositoryPorts) error,
) error {
	return fn(m)
}
func (m *mockOutboxRepository) CreateUser(
	ctx context.Context,
	user db.CreateUserParams,
) (db.User, error) {
	return db.User{}, nil
}

func (m *mockOutboxRepository) GetUser(
	ctx context.Context,
	id pgtype.UUID,
) (db.GetUserRow, error) {
	return db.GetUserRow{}, nil
}

func (m *mockOutboxRepository) InsertPurchaseIfNew(
	ctx context.Context,
	purchase events.Purchase,
) (bool, error) {
	return true, nil
}

func (m *mockOutboxRepository) GetPurchaseStats(
	ctx context.Context,
	userID pgtype.UUID,
) (events.PurchaseStats, error) {
	return events.PurchaseStats{}, nil
}

func (m *mockOutboxRepository) UnlockAchievementIfNew(
	ctx context.Context,
	userID pgtype.UUID,
	achievementCode string,
) (bool, error) {
	return false, nil
}

func (m *mockOutboxRepository) GetUnlockedAchievementCount(
	ctx context.Context,
	userID pgtype.UUID,
) (int, error) {
	return 0, nil
}

func (m *mockOutboxRepository) UnlockBadgeIfNew(
    ctx context.Context,
    userID pgtype.UUID,
    badgeCode string,
) (bool, error) {
    return false, nil
}

func (m *mockOutboxRepository) AddOutboxEvent(
    ctx context.Context,
    event domain.Event,
) error {
    return nil
}

func (m *mockOutboxRepository) CreatePaymentIfNew(
    ctx context.Context,
    payment domain.Payment,
) (bool, error) {
    return false, nil
}

func (m *mockOutboxRepository) MarkPaymentSuccessful(
	ctx context.Context,
	id pgtype.UUID,
	providerRef string,
) error {
	return nil
}

func (m *mockOutboxRepository) GetPaymentByUserAndBadge(
	ctx context.Context,
	userID pgtype.UUID,
	badgeCode string,
) (*domain.Payment, error) {
	panic("not implemented")
}

func (m *mockOutboxRepository) MarkPaymentFailed(
	ctx context.Context,
	id pgtype.UUID,
	reason string,
) error {
	return nil
}

func (m *mockOutboxRepository) GetPendingOutboxEvents(
	ctx context.Context,
	batchSize int,
) ([]db.OutboxEvent, error) {
	return []db.OutboxEvent{}, nil
}

func (m *mockOutboxRepository) GetUserAchievements(
	ctx context.Context,
	userID pgtype.UUID,
) (db.GetUserAchievementsRow, error) {
	panic("not implemented")
}

func (m *mockOutboxRepository) MarkOutboxEventProcessed(
	ctx context.Context,
	id pgtype.UUID,
) error {
	m.markProcessedCalls++
	m.markedIDs = append(m.markedIDs, id)

	if m.markProcessedFn != nil {
		return m.markProcessedFn(ctx, id)
	}

	return nil
}

var _ ports.RepositoryPorts = (*mockOutboxRepository)(nil)

// ------------------------------------------------------------
// Mock Event Publisher
// ------------------------------------------------------------

type mockEventPublisher struct {
	publishFn func(
		context.Context,
		domain.Event,
	) error

	publishedEvents []domain.Event
}

func (m *mockEventPublisher) Publish(
	ctx context.Context,
	event domain.Event,
) error {
	m.publishedEvents = append(
		m.publishedEvents,
		event,
	)

	if m.publishFn != nil {
		return m.publishFn(ctx, event)
	}

	return nil
}

var _ events.EventPublisher = (*mockEventPublisher)(nil)

// ------------------------------------------------------------
// Helpers
// ------------------------------------------------------------

func newOutboxEvent(
	eventType string,
) db.OutboxEvent {
	id := uuid.New()
	aggregateID := uuid.New()

	return db.OutboxEvent{
		ID: pgtype.UUID{
			Bytes: id,
			Valid: true,
		},
		EventType: eventType,
		AggregateID: pgtype.UUID{
			Bytes: aggregateID,
			Valid: true,
		},
		Payload: []byte(`{"test":"payload"}`),
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	}
}
