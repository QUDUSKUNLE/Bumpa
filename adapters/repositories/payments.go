package repositories

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/google/uuid"
)

func (r *Repository) CreatePaymentIfNew(ctx context.Context, payment domain.Payment) (bool, error) {
	// insert into user_achievements with ON CONFLICT DO NOTHING
	return true, nil
}

func (r *Repository) MarkPaymentFailed(ctx context.Context, id uuid.UUID, providerRef string) error {
	return nil
}

func (r *Repository) MarkPaymentSuccessful(ctx context.Context, id uuid.UUID, reason string) error {
	return nil
}
