package events

import (
	"github.com/google/uuid"
)

type Purchase struct {
	ID             uuid.UUID
	user           uuid.UUID
	ExternalID     string
	AmountKobo     float32
	PaymentAccount string //
}
