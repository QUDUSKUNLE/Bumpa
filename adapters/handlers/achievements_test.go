package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/services"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePurchase_Success(t *testing.T) {
	e := echo.New()

	userID := uuid.New()
	purchaseID := uuid.New()

	reqBody := map[string]any{
		"id":              purchaseID.String(),
		"user":            userID.String(),
		"external_id":     "ORDER-001",
		"amount_kobo":     30000,
		"payment_account": "RCP_xxxxx",
	}

	reqJSON, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/users/purchases",
		bytes.NewReader(reqJSON),
	)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var receivedPurchase events.Purchase
	var receivedUserID pgtype.UUID

	mockAchievement := &MockAchievementService{
		ProcessPurchaseFunc: func(
			ctx context.Context,
			id pgtype.UUID,
			purchase events.Purchase,
		) error {
			receivedUserID = id
			receivedPurchase = purchase
			return nil
		},
	}

	service := services.ServicesHandler{
		AchievementService: mockAchievement,
	}

	handler := NewHttpAdapter(service)

	err = handler.CreatePurchase(c)

	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	assert.Equal(t, userID, uuid.UUID(receivedUserID.Bytes))
	assert.True(t, receivedUserID.Valid)
	assert.Equal(t, purchaseID, receivedPurchase.ID)
	assert.Equal(t, userID, receivedPurchase.User)
	assert.Equal(t, "ORDER-001", receivedPurchase.ExternalID)
	assert.Equal(t, int64(30000), receivedPurchase.AmountKobo)
	assert.Equal(t, "RCP_xxxxx", receivedPurchase.PaymentAccount)

	var response map[string]any

	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(
		t,
		"Purchase made successfully",
		response["status"],
	)
}

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
