package cashback

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type mockCashbackRepository struct {
	getUserFn                  func(context.Context, pgtype.UUID) (db.GetUserRow, error)
	getPaymentByUserAndBadgeFn func(
		context.Context,
		pgtype.UUID,
		*string,
	) (*domain.Payment, error)
	createPaymentIfNewFn func(context.Context, domain.Payment) (*bool, error)

	getUserCalls       int
	getPaymentCalls    int
	createPaymentCalls int
}

func (m *mockCashbackRepository) GetUser(
	ctx context.Context,
	id pgtype.UUID,
) (db.GetUserRow, error) {
	m.getUserCalls++

	if m.getUserFn != nil {
		return m.getUserFn(ctx, id)
	}

	return db.GetUserRow{}, nil
}

func (m *mockCashbackRepository) GetPaymentByUserAndBadge(
	ctx context.Context,
	userID pgtype.UUID,
	badgeCode *string,
) (*domain.Payment, error) {
	m.getPaymentCalls++

	if m.getPaymentByUserAndBadgeFn != nil {
		return m.getPaymentByUserAndBadgeFn(
			ctx,
			userID,
			badgeCode,
		)
	}

	return nil, pgx.ErrNoRows
}

func (m *mockCashbackRepository) CreatePaymentIfNew(
	ctx context.Context,
	payment domain.Payment,
) (*bool, error) {
	m.createPaymentCalls++

	if m.createPaymentIfNewFn != nil {
		return m.createPaymentIfNewFn(ctx, payment)
	}

	value := true
	return &value, nil
}

// -----------------------------------------------------------------------------
// RepositoryPorts methods not used by CashbackService.
// -----------------------------------------------------------------------------

func (m *mockCashbackRepository) WithTx(
	ctx context.Context,
	fn func(tx interface{}) error,
) error {
	panic("not implemented")
}
