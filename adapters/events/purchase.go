package events

import (
	"github.com/google/uuid"
)

type Purchase struct {
	ID             string
	user           uuid.UUID
	ExternalID     string
	AmountKobo     float32
	PaymentAccount string //
}
