package events

type User struct {
	ID             string
	Name           string
	Email          string
	Phone          string
	PaymentAccount string // e.g. Paystack recipient/account reference
}
