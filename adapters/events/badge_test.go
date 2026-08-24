package events

import (
	"encoding/json"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBadgeUnlockedPayload_JSON(t *testing.T) {
	userID := pgtype.UUID{
		Bytes: [16]byte{
			0x3c, 0x24, 0x0a, 0xd4,
			0x38, 0xe0, 0x4e, 0xab,
			0x85, 0x50, 0xf7, 0x4b,
			0xb6, 0x4a, 0xe9, 0x99,
		},
		Valid: true,
	}

	payload := BadgeUnlockedPayload{
		BadgeName: "Bronze Badge",
		BadgeCode: "bronze",
		User:      userID,
	}

	data, err := json.Marshal(payload)

	require.NoError(t, err)

	assert.JSONEq(
		t,
		`{
			"badge_name": "Bronze Badge",
			"badge_code": "bronze",
			"user": "3c240ad4-38e0-4eab-8550-f74bb64ae999"
		}`,
		string(data),
	)
}

func TestBadgeUnlockedPayload_JSON_Unmarshal(t *testing.T) {
	data := []byte(`{
		"badge_name": "Silver Badge",
		"badge_code": "silver",
		"user": "3c240ad4-38e0-4eab-8550-f74bb64ae999"
	}`)

	var payload BadgeUnlockedPayload

	err := json.Unmarshal(data, &payload)

	require.NoError(t, err)

	assert.Equal(t, "Silver Badge", payload.BadgeName)
	assert.Equal(t, "silver", payload.BadgeCode)
	assert.True(t, payload.User.Valid)

	expectedUUID := [16]byte{
		0x3c, 0x24, 0x0a, 0xd4,
		0x38, 0xe0, 0x4e, 0xab,
		0x85, 0x50, 0xf7, 0x4b,
		0xb6, 0x4a, 0xe9, 0x99,
	}

	assert.Equal(t, expectedUUID, payload.User.Bytes)
}
