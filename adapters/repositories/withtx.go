package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
)

func (r *Repository) WithTx(ctx context.Context, fn func(tx Repository) error) error {
	_, err := db.WithTx(ctx, r.pool, func(tx pgx.Tx) (Repository, error) {
		txQueries := db.New(tx)

		return Repository{
			queries: txQueries,
			pool: r.pool,
		}, fn(Repository{
			queries: txQueries,
			pool: r.pool,
		})
	})
	return err
}
