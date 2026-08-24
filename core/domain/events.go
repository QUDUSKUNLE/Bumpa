package domain

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Type           string    `json:"type"`
	OccurredAt     time.Time `json:"occurred_at"`
	AggregateID    uuid.UUID `json:"aggregate_id"`
	PaymentAccount string    `json:"payment_account"`
	Payload        []byte    `json:"payload"`
}
