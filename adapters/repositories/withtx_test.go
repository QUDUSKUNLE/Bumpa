package repositories

import (
	"context"
	"testing"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/core/ports"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRepository_WithTx_Commits(t *testing.T) {
	ctx := context.Background()
	repo := setupTestRepository(t)

	email := uuid.New().String() + "@example.com"

	err := repo.WithTx(ctx, func(txRepo ports.RepositoryPorts) error {
		_, err := txRepo.CreateUser(ctx, db.CreateUserParams{
			Name:  "Transaction Commit User",
			Email: email,
		})

		return err
	})

	require.NoError(t, err)
}
