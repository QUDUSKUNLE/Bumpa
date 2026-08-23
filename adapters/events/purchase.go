package events

import (
	"github.com/google/uuid"
)

type Purchase struct {
	ID             uuid.UUID `json:"id"`
	User           uuid.UUID `json:"user"`
	ExternalID     string    `json:"external_id"`
	AmountKobo     float32   `json:"amount_kobo"`
	PaymentAccount string    `json:"payment_account"`
}
