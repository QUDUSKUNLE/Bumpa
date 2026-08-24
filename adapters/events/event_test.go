package events

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvent_JSON(t *testing.T) {
	id := uuid.MustParse("3c240ad4-38e0-4eab-8550-f74bb64ae999")
	userID := uuid.MustParse("7c240ad4-38e0-4eab-8550-f74bb64ae999")
	aggregateID := uuid.MustParse("9c240ad4-38e0-4eab-8550-f74bb64ae999")

	eventType := "AchievementUnlocked"
	paymentAccount := "RCP_xxxxx"
	occurredAt := time.Date(
		2026,
		time.August,
		24,
		17,
		0,
		0,
		0,
		time.UTC,
	)

	event := Event{
		ID:             id,
		UserID:         userID,
		Type:           eventType,
		OccurredAt:     occurredAt,
		AggregateID:    aggregateID,
		Payload:        []byte(`{"achievement_name":"First Purchase"}`),
		PaymentAccount: paymentAccount,
	}

	data, err := json.Marshal(event)

	require.NoError(t, err)

	assert.JSONEq(
		t,
		`{
			"id": "3c240ad4-38e0-4eab-8550-f74bb64ae999",
			"user_id": "7c240ad4-38e0-4eab-8550-f74bb64ae999",
			"type": "AchievementUnlocked",
			"occurred_at": "2026-08-24T17:00:00Z",
			"aggregate_id": "9c240ad4-38e0-4eab-8550-f74bb64ae999",
			"payload": "eyJhY2hpZXZlbWVudF9uYW1lIjoiRmlyc3QgUHVyY2hhc2UifQ==",
			"payment_account": "RCP_xxxxx"
		}`,
		string(data),
	)
}
