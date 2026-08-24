package repositories


import (
	"context"
	"testing"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_CreateUser(t *testing.T) {
	ctx := context.Background()
	repo := setupTestRepository(t)

	email := uuid.New().String() + "@example.com"

	user, err := repo.CreateUser(ctx, db.CreateUserParams{
		Name:  "Create User Test",
		Email: email,
	})

	require.NoError(t, err)

	assert.True(t, user.ID.Valid)
	assert.Equal(t, "Create User Test", user.Name)
	assert.Equal(t, email, user.Email)
}

func TestRepository_GetUser(t *testing.T) {
	ctx := context.Background()
	repo := setupTestRepository(t)

	email := uuid.New().String() + "@example.com"

	createdUser, err := repo.CreateUser(ctx, db.CreateUserParams{
		Name:  "Get User Test",
		Email: email,
	})

	require.NoError(t, err)
	require.True(t, createdUser.ID.Valid)

	user, err := repo.GetUser(ctx, createdUser.ID)

	require.NoError(t, err)

	assert.Equal(t, createdUser.ID, user.ID)
	assert.Equal(t, "Get User Test", user.Name)
	assert.Equal(t, email, user.Email)
}
