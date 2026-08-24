package domain

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Payment struct {
	ID          pgtype.UUID
	UserID      pgtype.UUID
	BadgeCode   string
	AmountKobo  int64
	Status      string
	CreatedAt   time.Time
	ProviderRef string
}

type PaymentResult struct {
	ProviderReference string
	Status            string
}
