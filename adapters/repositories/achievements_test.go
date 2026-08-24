package repositories

import (
	"context"
	"testing"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_UnlockAchievementIfNew(t *testing.T) {
	ctx := context.Background()
	repo := setupTestRepository(t)

	user, err := repo.CreateUser(ctx, db.CreateUserParams{
		Name:  "Achievement Test User",
		Email: uuid.New().String() + "@example.com",
	})
	require.NoError(t, err)

	userUUID := user.ID

	achievementCode := "first_purchase"

	unlocked, err := repo.UnlockAchievementIfNew(
		ctx,
		userUUID,
		achievementCode,
	)
	require.NoError(t, err)
	assert.True(t, unlocked)

	unlocked, err = repo.UnlockAchievementIfNew(
		ctx,
		userUUID,
		achievementCode,
	)
	require.NoError(t, err)
	assert.False(t, unlocked)
}

func TestRepository_GetUserAchievements(t *testing.T) {
	ctx := context.Background()
	repo := setupTestRepository(t)

	user, err := repo.CreateUser(ctx, db.CreateUserParams{
		Name:  "Get Achievement Test User",
		Email: uuid.New().String() + "@example.com",
	})
	require.NoError(t, err)

	userUUID := user.ID

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

	result, err := repo.GetUserAchievements(
		ctx,
		userUUID,
	)
	require.NoError(t, err)

	assert.Contains(
		t,
		result.UnlockedAchievements,
		"First Purchase",
	)

	assert.Contains(
		t,
		result.UnlockedAchievements,
		"Three Purchases",
	)
}
