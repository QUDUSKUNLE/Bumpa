package domain

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type CashbackRequest struct {
	UserID         pgtype.UUID
	PaymentAccount string
	AmountKobo     int64
	Reference      string
	Reason         string
}
