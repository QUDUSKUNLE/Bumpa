package repositories

import (
	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewRepository(q *db.Queries, pool *pgxpool.Pool) *Repository {
	return &Repository{
		queries: q,
		pool:    pool,
	}
}
