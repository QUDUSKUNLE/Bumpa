package repositories

import (
	"context"
	"os"
	"testing"

	"github.com/QUDUSKUNLE/Bumpa/adapters/config"
	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/stretchr/testify/require"
)

func setupTestRepository(t *testing.T) *Repository {
	t.Helper()

	databaseURL := os.Getenv("TEST_DATABASE_URL")

	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}

	require.NotEmpty(
		t,
		databaseURL,
		"TEST_DATABASE_URL or DATABASE_URL must be set",
	)

	ctx := context.Background()

	dbConfig := config.DBConfig()

	store, conn, err := db.DatabaseConnection(
		ctx,
		databaseURL,
		dbConfig,
	)

	require.NoError(t, err)

	t.Cleanup(func() {
		conn.Close()
	})

	return NewRepository(store, conn)
}
