package repositories

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *Repository) CreateUser(ctx context.Context, user db.CreateUserParams) (db.User, error) {
	return r.queries.CreateUser(ctx, user)
}

func (r *Repository) GetUser(ctx context.Context, id uuid.UUID) (db.GetUserRow, error) {
	pgID := pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
	return r.queries.GetUser(ctx, pgID)
}
