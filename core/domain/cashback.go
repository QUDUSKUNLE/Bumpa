package domain

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type CashbackRequest struct {
	Source           string
	UserID           pgtype.UUID
	PaymentAccount   string
	RecipientAccount string
	Currency         string
	AmountKobo       int64
	Reference        string
	AccountReference string
	Reason           string
}

type FinaliseCashBackRequest struct {
	TransferCode string
	Otp          string
}
