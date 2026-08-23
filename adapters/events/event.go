package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Type        string
	OccurredAt  time.Time
	AggregateID uuid.UUID
	Payload     json.RawMessage
}
