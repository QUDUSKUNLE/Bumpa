package repositories

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/QUDUSKUNLE/Bumpa/core/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *Repository) CreatePaymentIfNew(ctx context.Context, payment domain.Payment) (bool, error) {
	_, err := r.queries.CreatePayment(ctx, db.CreatePaymentParams{
		ID:         pgtype.UUID{Bytes: uuid.New(), Valid: true},
		UserID:     pgtype.UUID{Bytes: payment.UserID, Valid: true},
		BadgeCode:  payment.BadgeCode,
		AmountKobo: payment.AmountKobo,
		Status:     payment.Status,
		CreatedAt: pgtype.Timestamptz{
			Time:  payment.CreatedAt,
			Valid: true,
		},
	})
	if err != nil {
		utils.LogError("CreatePayment error", err)
		return false, err
	}
	return true, nil
}

func (r *Repository) MarkPaymentFailed(ctx context.Context, id uuid.UUID, providerRef string) error {
	if err := r.queries.MarkPaymentFailed(ctx, db.MarkPaymentFailedParams{
		ID:                pgtype.UUID{Bytes: id, Valid: true},
		ProviderReference: pgtype.Text{String: providerRef, Valid: true},
	}); err != nil {
		utils.LogError("MarkPaymentFailed error", err)
		return err
	}
	return nil
}

func (r *Repository) MarkPaymentSuccessful(ctx context.Context, id uuid.UUID, reason string) error {
	if err := r.queries.MarkPaymentSuccessful(ctx, db.MarkPaymentSuccessfulParams{
		ID:                pgtype.UUID{Bytes: id, Valid: true},
		ProviderReference: pgtype.Text{String: reason, Valid: true},
	}); err != nil {
		utils.LogError("MarkPaymentSuccessful error", err)
		return err
	}
	return nil
}
