package repositories

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/jackc/pgx/v5/pgtype"
	// "github.com/google/uuid"
	// "github.com/jackc/pgx/v5"
)

func (r *Repository) CreateUser(ctx context.Context, user db.CreateUserParams) (db.CreateUserRow, error) {
	return r.CreateUser(ctx, user)
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (db.GetUserByEmailRow, error) {
	// map db row to domain.User
	return r.GetUserByEmail(ctx, email)
}

func (r *Repository) GetUser(ctx context.Context, id pgtype.UUID) (db.GetUserRow, error) {
	return r.GetUser(ctx, id)
}

func (r *Repository) ListUsers(ctx context.Context) ([]db.ListUsersRow, error) {
	return r.ListUsers(ctx)
}
