package repositories

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/core/ports"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) WithTx(ctx context.Context, fn func(tx ports.RepositoryPorts) error) error {
	_, err := db.WithTx(ctx, r.pool, func(tx pgx.Tx) (*Repository, error) {
		txQueries := db.New(tx)
		txRepo := &Repository{
			queries: txQueries,
			pool:    r.pool,
		}
		return txRepo, fn(txRepo)
	})
	return err
}
