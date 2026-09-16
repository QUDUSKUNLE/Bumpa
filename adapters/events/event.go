package events

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID             uuid.UUID `json:"id"`
	Type           string    `json:"type"`
	UserID         uuid.UUID `json:"user_id"`
	OccurredAt     time.Time `json:"occurred_at"`
	AggregateID    uuid.UUID `json:"aggregate_id"`
	Payload        []byte    `json:"payload"`
	PaymentAccount string    `json:"payment_account"`
}
