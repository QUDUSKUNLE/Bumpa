package repositories

import (
	"context"
	"errors"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/QUDUSKUNLE/Bumpa/core/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *Repository) CreatePaymentIfNew(ctx context.Context, payment domain.Payment) (bool, error) {
	_, err := r.queries.CreatePayment(ctx, db.CreatePaymentParams{
		ID:         pgtype.UUID{Bytes: uuid.New(), Valid: true},
		UserID:     payment.UserID,
		BadgeCode:  payment.BadgeCode,
		AmountKobo: payment.AmountKobo,
		Status:     payment.Status,
		CreatedAt: pgtype.Timestamptz{
			Time:  payment.CreatedAt,
			Valid: true,
		},
	})
	if errors.Is(err, pgx.ErrNoRows) {
		// Payment Already Exist
		return false, nil
	}
	if err != nil {
		utils.LogError("CreatePayment Repositories Error: %v", err)
		return false, err
	}
	return true, nil
}

func (r *Repository) MarkPaymentFailed(ctx context.Context, id pgtype.UUID, providerRef string) error {
	if err := r.queries.MarkPaymentFailed(ctx, db.MarkPaymentFailedParams{
		ID:                id,
		ProviderReference: pgtype.Text{String: providerRef, Valid: true},
	}); err != nil {
		utils.LogError("MarkPaymentFailed Repositories Error: %v", err)
		return err
	}
	return nil
}

func (r *Repository) MarkPaymentSuccessful(ctx context.Context, id pgtype.UUID, reason string) error {
	if err := r.queries.MarkPaymentSuccessful(ctx, db.MarkPaymentSuccessfulParams{
		ID:                id,
		ProviderReference: pgtype.Text{String: reason, Valid: true},
	}); err != nil {
		utils.LogError("MarkPaymentSuccessful Repositories Error: %v", err)
		return err
	}
	return nil
}

func (r *Repository) GetPaymentByUserAndBadge(ctx context.Context, userID pgtype.UUID, badgeCode string) (*domain.Payment, error) {
	payment, err := r.queries.GetPaymentByUserAndBadge(ctx, db.GetPaymentByUserAndBadgeParams{UserID: userID, BadgeCode: badgeCode})
	if err != nil {
		return nil, err
	}
	return &domain.Payment{
		ID:          payment.ID,
		UserID:      payment.UserID,
		BadgeCode:   payment.BadgeCode,
		AmountKobo:  payment.AmountKobo,
		Status:      payment.Status,
		ProviderRef: payment.ProviderReference.String,
		CreatedAt:   payment.CreatedAt.Time,
	}, nil
}
