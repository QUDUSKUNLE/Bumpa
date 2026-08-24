package events

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAchievementUnlockedPayload_JSON(t *testing.T) {
	payload := AchievementUnlockedPayload{
		AchievementName: "First Purchase",
	}

	data, err := json.Marshal(payload)

	require.NoError(t, err)

	assert.JSONEq(
		t,
		`{"achievement_name":"First Purchase"}`,
		string(data),
	)
}

func TestAchievementUnlockedPayload_JSON_Unmarshal(t *testing.T) {
	data := []byte(`{"achievement_name":"Three Purchases"}`)

	var payload AchievementUnlockedPayload

	err := json.Unmarshal(data, &payload)

	require.NoError(t, err)

	assert.Equal(t, "Three Purchases", payload.AchievementName)
}
