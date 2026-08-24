package events

import "github.com/google/uuid"

type PurchaseRequest struct {
	ID             string `json:"id"`
	User           string `json:"user"`
	ExternalID     string `json:"external_id"`
	AmountKobo     int64  `json:"amount_kobo"`
	PaymentAccount string `json:"payment_account"`
}

type Purchase struct {
	ID             uuid.UUID `json:"id"`
	User           uuid.UUID `json:"user"`
	ExternalID     string    `json:"external_id"`
	AmountKobo     int64     `json:"amount_kobo"`
	PaymentAccount string    `json:"payment_account"`
}
