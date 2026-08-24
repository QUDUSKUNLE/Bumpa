package repositories

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *Repository) CreateUser(ctx context.Context, user db.CreateUserParams) (db.User, error) {
	return r.queries.CreateUser(ctx, user)
}

func (r *Repository) GetUser(ctx context.Context, id pgtype.UUID) (db.GetUserRow, error) {
	return r.queries.GetUser(ctx, id)
}
