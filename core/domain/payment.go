package domain

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID          uuid.UUID
	UserID      uuid.UUID
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
