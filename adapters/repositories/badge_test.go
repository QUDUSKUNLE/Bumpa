package repositories

import (
	"context"
	"testing"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_GetUnlockedAchievementCount(t *testing.T) {
	ctx := context.Background()
	repo := setupTestRepository(t)

	// Create a real user so the user_id foreign key is valid.
	user, err := repo.CreateUser(ctx, db.CreateUserParams{
		Name:  "Achievement Count Test User",
		Email: uuid.New().String() + "@example.com",
	})
	require.NoError(t, err)

	userUUID := user.ID

	// A new user should have zero achievements.
	count, err := repo.GetUnlockedAchievementCount(ctx, userUUID)

	require.NoError(t, err)
	require.NotNil(t, count)
	assert.Equal(t, 0, count)

	// Unlock two achievements.
	_, err = repo.UnlockAchievementIfNew(
		ctx,
		userUUID,
		"first_purchase",
	)
	require.NoError(t, err)

	_, err = repo.UnlockAchievementIfNew(
		ctx,
		userUUID,
		"three_purchases",
	)
	require.NoError(t, err)

	// The count should now be two.
	count, err = repo.GetUnlockedAchievementCount(ctx, userUUID)

	require.NoError(t, err)
	require.NotNil(t, count)
	assert.Equal(t, 2, count)
}

func TestRepository_UnlockBadgeIfNew(t *testing.T) {
	ctx := context.Background()
	repo := setupTestRepository(t)

	// Create a real user so the user_id foreign key is valid.
	user, err := repo.CreateUser(ctx, db.CreateUserParams{
		Name:  "Badge Test User",
		Email: uuid.New().String() + "@example.com",
	})
	require.NoError(t, err)

	userUUID := user.ID

	badgeCode := "bronze"

	// First unlock should succeed.
	unlocked, err := repo.UnlockBadgeIfNew(
		ctx,
		userUUID,
		badgeCode,
	)

	require.NoError(t, err)
	require.NotNil(t, unlocked)
	assert.True(t, unlocked)

	// Unlocking the same badge again should return false.
	unlocked, err = repo.UnlockBadgeIfNew(
		ctx,
		userUUID,
		badgeCode,
	)

	require.NoError(t, err)
	require.NotNil(t, unlocked)
	assert.False(t, unlocked)
}
