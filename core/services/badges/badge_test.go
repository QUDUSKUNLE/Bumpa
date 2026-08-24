package badges

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/QUDUSKUNLE/Bumpa/core/ports"
	"github.com/QUDUSKUNLE/Bumpa/core/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockBadgeRepository struct {
	ports.RepositoryPorts

	purchaseStats events.PurchaseStats
	statsErr      error

	unlockBadgeFn func(
		context.Context,
		pgtype.UUID,
		string,
	) (bool, error)

	addOutboxEventFn func(
		context.Context,
		domain.Event,
	) error

	unlockedBadges []string
	outboxEvents   []domain.Event
}

func (m *mockBadgeRepository) WithTx(
	ctx context.Context,
	fn func(ports.RepositoryPorts) error,
) error {
	return fn(m)
}

func (m *mockBadgeRepository) GetPurchaseStats(
	ctx context.Context,
	userID pgtype.UUID,
) (events.PurchaseStats, error) {
	return m.purchaseStats, m.statsErr
}

func (m *mockBadgeRepository) UnlockBadgeIfNew(
	ctx context.Context,
	userID pgtype.UUID,
	badgeCode string,
) (bool, error) {
	m.unlockedBadges = append(m.unlockedBadges, badgeCode)

	if m.unlockBadgeFn != nil {
		return m.unlockBadgeFn(ctx, userID, badgeCode)
	}

	return true, nil
}

func (m *mockBadgeRepository) AddOutboxEvent(
	ctx context.Context,
	event domain.Event,
) error {
	m.outboxEvents = append(m.outboxEvents, event)

	if m.addOutboxEventFn != nil {
		return m.addOutboxEventFn(ctx, event)
	}

	return nil
}

func TestBadgeDefinition(t *testing.T) {
	definitions := BadgeDefinition()

	require.Len(t, definitions, 5)

	expected := []struct {
		code     string
		name     string
		required int
	}{
		{"bronze", "Bronze", 1},
		{"silver", "Silver", 3},
		{"gold", "Gold", 5},
		{"platinum", "Platinum", 10},
		{"diamond", "Diamond", 20},
	}

	for i, want := range expected {
		t.Run(want.code, func(t *testing.T) {
			assert.Equal(t, want.code, definitions[i].Code)
			assert.Equal(t, want.name, definitions[i].Name)
			assert.Equal(t, want.required, definitions[i].RequiredRewards)
		})
	}
}

func TestBadgeService_HandleAchievementUnlocked_UnlocksEligibleBadges(
	t *testing.T,
) {
	ctx := context.Background()

	userID := uuid.New()

	repo := &mockBadgeRepository{
		purchaseStats: events.PurchaseStats{
			TotalPurchases: 5,
		},
	}

	service := NewBadgeService(
		repo,
		BadgeDefinition(),
		nil,
	)

	payload := events.AchievementUnlockedPayload{
		AchievementName: "Five Purchases",
	}

	payloadJSON, err := json.Marshal(payload)
	require.NoError(t, err)

	event := events.Event{
		ID:             uuid.New(),
		UserID:         userID,
		AggregateID:    userID,
		Type:           "AchievementUnlocked",
		Payload:        payloadJSON,
		PaymentAccount: "test-account",
	}

	err = service.HandleAchievementUnlocked(ctx, event)

	require.NoError(t, err)

	// 5 purchases qualifies for:
	// Bronze = 1
	// Silver = 3
	// Gold = 5
	//
	// It should NOT unlock Platinum or Diamond.
	assert.Equal(
		t,
		[]string{
			"bronze",
			"silver",
			"gold",
		},
		repo.unlockedBadges,
	)

	require.Len(t, repo.outboxEvents, 3)

	assert.Equal(t, "BadgeUnlocked", repo.outboxEvents[0].Type)
	assert.Equal(t, "BadgeUnlocked", repo.outboxEvents[1].Type)
	assert.Equal(t, "BadgeUnlocked", repo.outboxEvents[2].Type)
}

func TestBadgeService_HandleAchievementUnlocked_DoesNotUnlockIneligibleBadges(
	t *testing.T,
) {
	ctx := context.Background()

	userID := uuid.New()

	repo := &mockBadgeRepository{
		purchaseStats: events.PurchaseStats{
			TotalPurchases: 0,
		},
	}

	service := NewBadgeService(
		repo,
		BadgeDefinition(),
		nil,
	)

	event := events.Event{
		ID:          uuid.New(),
		UserID:      userID,
		AggregateID: userID,
		Type:        "AchievementUnlocked",
		Payload: utils.MustJSON(events.AchievementUnlockedPayload{
			AchievementName: "Something",
		}),
	}

	err := service.HandleAchievementUnlocked(ctx, event)

	require.NoError(t, err)

	assert.Empty(t, repo.unlockedBadges)
	assert.Empty(t, repo.outboxEvents)
}

func TestBadgeService_HandleAchievementUnlocked_DoesNotCreateOutboxForExistingBadge(
	t *testing.T,
) {
	ctx := context.Background()

	userID := uuid.New()

	repo := &mockBadgeRepository{
		purchaseStats: events.PurchaseStats{
			TotalPurchases: 5,
		},
		unlockBadgeFn: func(
			ctx context.Context,
			userID pgtype.UUID,
			badgeCode string,
		) (bool, error) {
			// Badge already exists.
			return false, nil
		},
	}

	service := NewBadgeService(
		repo,
		BadgeDefinition(),
		nil,
	)

	event := events.Event{
		ID:          uuid.New(),
		UserID:      userID,
		AggregateID: userID,
		Type:        "AchievementUnlocked",
		Payload: utils.MustJSON(events.AchievementUnlockedPayload{
			AchievementName: "Five Purchases",
		}),
	}

	err := service.HandleAchievementUnlocked(ctx, event)

	require.NoError(t, err)

	assert.Equal(
		t,
		[]string{
			"bronze",
			"silver",
			"gold",
		},
		repo.unlockedBadges,
	)

	assert.Empty(t, repo.outboxEvents)
}

func TestBadgeService_HandleAchievementUnlocked_InvalidPayload(
	t *testing.T,
) {
	ctx := context.Background()

	repo := &mockBadgeRepository{}

	service := NewBadgeService(
		repo,
		BadgeDefinition(),
		nil,
	)

	event := events.Event{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		AggregateID: uuid.New(),
		Type:        "AchievementUnlocked",
		Payload:     []byte(`invalid-json`),
	}

	err := service.HandleAchievementUnlocked(ctx, event)

	require.Error(t, err)

	assert.Empty(t, repo.unlockedBadges)
	assert.Empty(t, repo.outboxEvents)
}

func TestBadgeService_HandleAchievementUnlocked_GetPurchaseStatsError(
	t *testing.T,
) {
	ctx := context.Background()

	expectedErr := errors.New("database error")

	repo := &mockBadgeRepository{
		statsErr: expectedErr,
	}

	service := NewBadgeService(
		repo,
		BadgeDefinition(),
		nil,
	)

	event := events.Event{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		AggregateID: uuid.New(),
		Type:        "AchievementUnlocked",
		Payload: utils.MustJSON(events.AchievementUnlockedPayload{
			AchievementName: "First Purchase",
		}),
	}

	err := service.HandleAchievementUnlocked(ctx, event)

	require.ErrorIs(t, err, expectedErr)

	assert.Empty(t, repo.unlockedBadges)
	assert.Empty(t, repo.outboxEvents)
}

func TestBadgeService_HandleAchievementUnlocked_UnlockBadgeError(
	t *testing.T,
) {
	ctx := context.Background()

	expectedErr := errors.New("failed to unlock badge")

	repo := &mockBadgeRepository{
		purchaseStats: events.PurchaseStats{
			TotalPurchases: 1,
		},
		unlockBadgeFn: func(
			ctx context.Context,
			userID pgtype.UUID,
			badgeCode string,
		) (bool, error) {
			return false, expectedErr
		},
	}

	service := NewBadgeService(
		repo,
		BadgeDefinition(),
		nil,
	)

	event := events.Event{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		AggregateID: uuid.New(),
		Type:        "AchievementUnlocked",
		Payload: utils.MustJSON(events.AchievementUnlockedPayload{
			AchievementName: "First Purchase",
		}),
	}

	err := service.HandleAchievementUnlocked(ctx, event)

	require.ErrorIs(t, err, expectedErr)

	assert.Empty(t, repo.outboxEvents)
}

func TestBadgeService_HandleAchievementUnlocked_AddOutboxEventError(
	t *testing.T,
) {
	ctx := context.Background()

	expectedErr := errors.New("outbox error")

	repo := &mockBadgeRepository{
		purchaseStats: events.PurchaseStats{
			TotalPurchases: 1,
		},
		addOutboxEventFn: func(
			ctx context.Context,
			event domain.Event,
		) error {
			return expectedErr
		},
	}

	service := NewBadgeService(
		repo,
		BadgeDefinition(),
		nil,
	)

	event := events.Event{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		AggregateID: uuid.New(),
		Type:        "AchievementUnlocked",
		Payload: utils.MustJSON(events.AchievementUnlockedPayload{
			AchievementName: "First Purchase",
		}),
	}

	err := service.HandleAchievementUnlocked(ctx, event)

	require.ErrorIs(t, err, expectedErr)

	require.Len(t, repo.outboxEvents, 1)

	assert.Equal(
		t,
		"BadgeUnlocked",
		repo.outboxEvents[0].Type,
	)
}
