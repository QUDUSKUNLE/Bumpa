package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QUDUSKUNLE/Bumpa/core/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserAchievements_Success(t *testing.T) {
	e := echo.New()

	userID := uuid.New()
	called := false

	mockAchievement := &MockAchievementService{
		GetUserAchievementsFunc: func(ctx echo.Context) error {
			called = true

			return ctx.JSON(http.StatusOK, map[string]any{
				"unlocked_achievements": []string{
					"First Purchase",
					"Three Purchases",
				},
				"next_available_achievements": []string{
					"Five Purchases",
				},
			})
		},
	}

	service := services.ServicesHandler{
		AchievementService: mockAchievement,
	}

	handler := NewHttpAdapter(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/"+userID.String()+"/achievements",
		nil,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("user")
	c.SetParamValues(userID.String())

	err := handler.GetUserAchievements(c)

	require.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))

	assert.Equal(
		t,
		[]any{"First Purchase", "Three Purchases"},
		response["unlocked_achievements"],
	)

	assert.Equal(
		t,
		[]any{"Five Purchases"},
		response["next_available_achievements"],
	)
}
