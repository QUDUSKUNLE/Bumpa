package domain

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID             uuid.UUID `json:"id"`
	PurchaseID     uuid.UUID `json:"purchase_id"`
	Type           string    `json:"type"`
	UserID         uuid.UUID `json:"user_id"`
	OccurredAt     time.Time `json:"occurred_at"`
	AggregateID    uuid.UUID `json:"aggregate_id"`
	PaymentAccount string    `json:"payment_account"`
	Payload        []byte    `json:"payload"`
}

type Purchase struct {
	ID             uuid.UUID `json:"id"`
	User           uuid.UUID `json:"user"`
	ExternalID     string    `json:"external_id"`
	AmountKobo     int64     `json:"amount_kobo"`
	PaymentAccount string    `json:"payment_account"`
}
