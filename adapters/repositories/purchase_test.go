package repositories

import (
	"context"
	"testing"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_InsertPurchaseIfNew(t *testing.T) {
	ctx := context.Background()
	repo := setupTestRepository(t)

	// Create a real user because purchases.user_id
	// has a foreign key to users.id.
	user, err := repo.CreateUser(ctx, db.CreateUserParams{
		Name:  "Purchase Test User",
		Email: uuid.New().String() + "@example.com",
	})
	require.NoError(t, err)

	purchaseID := uuid.New()

	purchase := events.Purchase{
		ID:         purchaseID,
		User:       user.ID.Bytes,
		ExternalID: "purchase-" + uuid.New().String(),
		AmountKobo: int64(50000),
	}

	// First insert should succeed.
	inserted, err := repo.InsertPurchaseIfNew(ctx, purchase)

	require.NoError(t, err)
	require.NotNil(t, inserted)
	assert.True(t, inserted)

	// Same purchase ID should not be inserted again.
	inserted, err = repo.InsertPurchaseIfNew(ctx, purchase)

	require.NoError(t, err)
	require.NotNil(t, inserted)
	assert.False(t, inserted)
}

func TestRepository_GetPurchaseStats(t *testing.T) {
	ctx := context.Background()
	repo := setupTestRepository(t)

	// Create a real user.
	user, err := repo.CreateUser(ctx, db.CreateUserParams{
		Name:  "Purchase Stats Test User",
		Email: uuid.New().String() + "@example.com",
	})
	require.NoError(t, err)

	userID := user.ID

	// A new user should have zero purchases.
	stats, err := repo.GetPurchaseStats(ctx, userID)

	require.NoError(t, err)
	assert.Equal(t, 0, stats.TotalPurchases)

	// Add first purchase.
	purchase1 := events.Purchase{
		ID:         uuid.New(),
		User:       userID.Bytes,
		ExternalID: "purchase-" + uuid.New().String(),
		AmountKobo: int64(50000),
	}

	inserted, err := repo.InsertPurchaseIfNew(ctx, purchase1)

	require.NoError(t, err)
	require.True(t, inserted)

	// Add second purchase.
	purchase2 := events.Purchase{
		ID:         uuid.New(),
		User:       userID.Bytes,
		ExternalID: "purchase-" + uuid.New().String(),
		AmountKobo: int64(75000),
	}

	inserted, err = repo.InsertPurchaseIfNew(ctx, purchase2)

	require.NoError(t, err)
	require.True(t, inserted)

	// User should now have two purchases.
	stats, err = repo.GetPurchaseStats(ctx, userID)

	require.NoError(t, err)
	assert.Equal(t, 2, stats.TotalPurchases)
}

func int64Ptr(value int64) *int64 {
	return &value
}
