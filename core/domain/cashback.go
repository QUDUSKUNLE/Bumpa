package domain

import "github.com/google/uuid"

type CashbackRequest struct {
	UserID         uuid.UUID
	PaymentAccount string
	AmountKobo     int64
	Reference      string
	Reason         string
}
