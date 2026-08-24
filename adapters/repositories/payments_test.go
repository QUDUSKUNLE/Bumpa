package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createPaymentTestUser(t *testing.T, ctx context.Context, repo *Repository) pgtype.UUID {
	t.Helper()

	user, err := repo.CreateUser(ctx, db.CreateUserParams{
		Name:  "Payment Test User",
		Email: uuid.New().String() + "@example.com",
	})
	require.NoError(t, err)

	return user.ID
}

func TestRepository_CreatePaymentIfNew(t *testing.T) {
	ctx := context.Background()
	repo := setupTestRepository(t)

	userID := createPaymentTestUser(t, ctx, repo)

	badgeCode := "bronze"

	payment := &domain.Payment{
		UserID:     userID,
		BadgeCode:  badgeCode,
		AmountKobo: 50000,
		Status:     "pending",
		CreatedAt:  time.Now().UTC(),
	}

	created, err := repo.CreatePaymentIfNew(ctx, *payment)

	require.NoError(t, err)
	require.NotNil(t, created)
	assert.True(t, created)
}

func TestRepository_MarkPaymentFailed(t *testing.T) {
	ctx := context.Background()
	repo := setupTestRepository(t)

	userID := createPaymentTestUser(t, ctx, repo)

	badgeCode := "bronze"

	payment := &domain.Payment{
		UserID:     userID,
		BadgeCode:  badgeCode,
		AmountKobo: 50000,
		Status:     "pending",
		CreatedAt:  time.Now().UTC(),
	}

	created, err := repo.CreatePaymentIfNew(ctx, *payment)

	require.NoError(t, err)
	require.True(t, created)

	savedPayment, err := repo.GetPaymentByUserAndBadge(
		ctx,
		userID,
		badgeCode,
	)

	require.NoError(t, err)
	require.NotNil(t, savedPayment)

	providerRef := "failed-provider-ref"

	err = repo.MarkPaymentFailed(
		ctx,
		savedPayment.ID,
		providerRef,
	)

	require.NoError(t, err)

	updatedPayment, err := repo.GetPaymentByUserAndBadge(
		ctx,
		userID,
		badgeCode,
	)

	require.NoError(t, err)
	assert.Equal(t, providerRef, updatedPayment.ProviderRef)
}

func TestRepository_MarkPaymentSuccessful(t *testing.T) {
	ctx := context.Background()
	repo := setupTestRepository(t)

	userID := createPaymentTestUser(t, ctx, repo)

	badgeCode := "silver"

	payment := &domain.Payment{
		UserID:     userID,
		BadgeCode:  badgeCode,
		AmountKobo: 100000,
		Status:     "pending",
		CreatedAt:  time.Now().UTC(),
	}

	created, err := repo.CreatePaymentIfNew(ctx, *payment)

	require.NoError(t, err)
	require.True(t, created)

	savedPayment, err := repo.GetPaymentByUserAndBadge(
		ctx,
		userID,
		badgeCode,
	)

	require.NoError(t, err)

	providerRef := "successful-provider-ref"

	err = repo.MarkPaymentSuccessful(
		ctx,
		savedPayment.ID,
		providerRef,
	)

	require.NoError(t, err)

	updatedPayment, err := repo.GetPaymentByUserAndBadge(
		ctx,
		userID,
		badgeCode,
	)

	require.NoError(t, err)
	assert.Equal(t, providerRef, updatedPayment.ProviderRef)
}

func TestRepository_GetPaymentByUserAndBadge(t *testing.T) {
	ctx := context.Background()
	repo := setupTestRepository(t)

	userID := createPaymentTestUser(t, ctx, repo)

	badgeCode := "gold"

	payment := &domain.Payment{
		UserID:     userID,
		BadgeCode:  badgeCode,
		AmountKobo: 250000,
		Status:     "pending",
		CreatedAt:  time.Now().UTC(),
	}

	created, err := repo.CreatePaymentIfNew(ctx, *payment)

	require.NoError(t, err)
	require.True(t, created)

	result, err := repo.GetPaymentByUserAndBadge(
		ctx,
		userID,
		badgeCode,
	)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, badgeCode, result.BadgeCode)
	assert.Equal(t, int64(250000), result.AmountKobo)
	assert.Equal(t, "pending", result.Status)
}
