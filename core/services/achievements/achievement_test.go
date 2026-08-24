package achievements

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/QUDUSKUNLE/Bumpa/core/ports"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockAchievementRepository struct {
	ports.RepositoryPorts

	insertPurchaseFn           func(context.Context, events.Purchase) (bool, error)
	getPurchaseStatsFn         func(context.Context, pgtype.UUID) (events.PurchaseStats, error)
	unlockAchievementFn        func(context.Context, pgtype.UUID, string) (bool, error)
	addOutboxEventFn           func(context.Context, domain.Event) error
	getPendingOutboxEventsFn   func(context.Context, int) ([]db.OutboxEvent, error)
	getUserFn                  func(context.Context, pgtype.UUID) (db.GetUserRow, error)
	markOutboxEventProcessedFn func(context.Context, pgtype.UUID) error
	getUserAchievementsFn      func(context.Context, pgtype.UUID) (db.GetUserAchievementsRow, error)
}

func (m *mockAchievementRepository) WithTx(
	ctx context.Context,
	fn func(ports.RepositoryPorts) error,
) error {
	return fn(m)
}

func (m *mockAchievementRepository) InsertPurchaseIfNew(
	ctx context.Context,
	purchase events.Purchase,
) (bool, error) {
	return m.insertPurchaseFn(ctx, purchase)
}

func (m *mockAchievementRepository) GetPurchaseStats(
	ctx context.Context,
	userID pgtype.UUID,
) (events.PurchaseStats, error) {
	return m.getPurchaseStatsFn(ctx, userID)
}

func (m *mockAchievementRepository) UnlockAchievementIfNew(
	ctx context.Context,
	userID pgtype.UUID,
	code string,
) (bool, error) {
	return m.unlockAchievementFn(ctx, userID, code)
}

func (m *mockAchievementRepository) AddOutboxEvent(
	ctx context.Context,
	event domain.Event,
) error {
	return m.addOutboxEventFn(ctx, event)
}

func (m *mockAchievementRepository) GetPendingOutboxEvents(
	ctx context.Context,
	batchSize int,
) ([]db.OutboxEvent, error) {
	return m.getPendingOutboxEventsFn(ctx, batchSize)
}

func (m *mockAchievementRepository) GetUser(
	ctx context.Context,
	userID pgtype.UUID,
) (db.GetUserRow, error) {
	return m.getUserFn(ctx, userID)
}

func (m *mockAchievementRepository) MarkOutboxEventProcessed(
	ctx context.Context,
	id pgtype.UUID,
) error {
	return m.markOutboxEventProcessedFn(ctx, id)
}

func (m *mockAchievementRepository) GetUserAchievements(
	ctx context.Context,
	userID pgtype.UUID,
) (db.GetUserAchievementsRow, error) {
	return m.getUserAchievementsFn(ctx, userID)
}

type mockAchievementBus struct {
	events.EventBus

	published  []domain.Event
	publishErr error
}

func (m *mockAchievementBus) Publish(
	ctx context.Context,
	event domain.Event,
) error {
	m.published = append(m.published, event)

	return m.publishErr
}

func TestAchievementDefinition(t *testing.T) {
	definitions := AchievementDefinition()

	require.Len(t, definitions, 5)

	tests := []struct {
		code      string
		name      string
		purchases int
		expected  bool
	}{
		{"first_purchase", "First Purchase", 0, false},
		{"first_purchase", "First Purchase", 1, true},
		{"three_purchases", "Three Purchases", 2, false},
		{"three_purchases", "Three Purchases", 3, true},
		{"five_purchases", "Five Purchases", 4, false},
		{"five_purchases", "Five Purchases", 5, true},
		{"ten_purchases", "Ten Purchases", 9, false},
		{"ten_purchases", "Ten Purchases", 10, true},
		{"twenty_purchases", "Twenty Purchases", 19, false},
		{"twenty_purchases", "Twenty Purchases", 20, true},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			var definition events.AchievementDefinition

			for _, def := range definitions {
				if def.Code == tt.code {
					definition = def
					break
				}
			}

			require.NotNil(t, definition.Condition)

			result := definition.Condition(events.PurchaseStats{
				TotalPurchases: tt.purchases,
			})

			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.name, definition.Name)
		})
	}
}

func TestAchievementService_ProcessPurchase_NewPurchaseUnlocksAchievement(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New()
	userUUID := pgtype.UUID{
		Bytes: userID,
		Valid: true,
	}

	inserted := true
	unlocked := true

	var unlockedCode string
	var createdEvent domain.Event

	repo := &mockAchievementRepository{
		insertPurchaseFn: func(
			ctx context.Context,
			purchase events.Purchase,
		) (bool, error) {
			return inserted, nil
		},

		getPurchaseStatsFn: func(
			ctx context.Context,
			id pgtype.UUID,
		) (events.PurchaseStats, error) {
			return events.PurchaseStats{
				TotalPurchases: 1,
			}, nil
		},

		unlockAchievementFn: func(
			ctx context.Context,
			id pgtype.UUID,
			code string,
		) (bool, error) {
			unlockedCode = code
			return unlocked, nil
		},

		addOutboxEventFn: func(
			ctx context.Context,
			event domain.Event,
		) error {
			createdEvent = event
			return nil
		},
	}

	bus := &mockAchievementBus{}

	service := NewAchievementService(
		repo,
		AchievementDefinition(),
		bus,
	)

	purchase := events.Purchase{
		ID:   uuid.New(),
		User: userID,
	}

	err := service.ProcessPurchase(
		ctx,
		userUUID,
		purchase,
	)

	require.NoError(t, err)

	assert.Equal(t, "first_purchase", unlockedCode)
	assert.Equal(t, "AchievementUnlocked", createdEvent.Type)
	assert.Equal(t, userID, createdEvent.UserID)
	assert.Equal(t, userID, createdEvent.AggregateID)

	var payload events.AchievementUnlockedPayload

	err = json.Unmarshal(createdEvent.Payload, &payload)

	require.NoError(t, err)
	assert.Equal(t, "First Purchase", payload.AchievementName)
}

func TestAchievementService_ProcessPurchase_ExistingPurchaseDoesNothing(t *testing.T) {
	ctx := context.Background()

	inserted := false

	unlockCalled := false
	outboxCalled := false

	repo := &mockAchievementRepository{
		insertPurchaseFn: func(
			ctx context.Context,
			purchase events.Purchase,
		) (bool, error) {
			return inserted, nil
		},

		getPurchaseStatsFn: func(
			ctx context.Context,
			id pgtype.UUID,
		) (events.PurchaseStats, error) {
			t.Fatal("GetPurchaseStats should not be called")
			return events.PurchaseStats{}, nil
		},

		unlockAchievementFn: func(
			ctx context.Context,
			id pgtype.UUID,
			code string,
		) (bool, error) {
			unlockCalled = true
			return unlockCalled, nil
		},

		addOutboxEventFn: func(
			ctx context.Context,
			event domain.Event,
		) error {
			outboxCalled = true
			return nil
		},
	}

	service := NewAchievementService(
		repo,
		AchievementDefinition(),
		&mockAchievementBus{},
	)

	err := service.ProcessPurchase(
		ctx,
		pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		events.Purchase{
			ID: uuid.New(),
		},
	)

	require.NoError(t, err)
	assert.False(t, unlockCalled)
	assert.False(t, outboxCalled)
}

func TestAchievementService_ProcessPurchase_RepositoryError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("database error")

	repo := &mockAchievementRepository{
		insertPurchaseFn: func(
			ctx context.Context,
			purchase events.Purchase,
		) (bool, error) {
			return false, expectedErr
		},
	}

	service := NewAchievementService(
		repo,
		AchievementDefinition(),
		&mockAchievementBus{},
	)

	err := service.ProcessPurchase(
		ctx,
		pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		events.Purchase{
			ID: uuid.New(),
		},
	)

	assert.ErrorIs(t, err, expectedErr)
}

func TestAchievementService_GetUserAchievements_InvalidUserID(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/not-a-uuid/achievements",
		nil,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/users/:user/achievements")
	c.SetParamNames("user")
	c.SetParamValues("not-a-uuid")

	service := NewAchievementService(
		&mockAchievementRepository{},
		AchievementDefinition(),
		&mockAchievementBus{},
	)

	err := service.GetUserAchievements(c)

	require.Error(t, err)

	httpErr, ok := err.(*echo.HTTPError)

	require.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "invalid user ID", httpErr.Message)
}

func TestAchievementService_GetUserAchievements_MissingUser(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/users//achievements",
		nil,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/users/:user/achievements")
	c.SetParamNames("user")
	c.SetParamValues("")

	service := NewAchievementService(
		&mockAchievementRepository{},
		AchievementDefinition(),
		&mockAchievementBus{},
	)

	err := service.GetUserAchievements(c)

	require.Error(t, err)

	httpErr, ok := err.(*echo.HTTPError)

	require.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "user identity is required", httpErr.Message)
}
