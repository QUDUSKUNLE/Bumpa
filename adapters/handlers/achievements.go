package handlers

import (
	"net/http"

	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func (handler *HttpHandler) GetUserAchievements(ctx echo.Context) error {
	return handler.servicesAdapter.AchievementService.GetUserAchievements(ctx)
}

func (h *HttpHandler) CreatePurchase(c echo.Context) error {
	var req events.PurchaseRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"status": false,
			"error":  err.Error(),
		})
	}

	utils.LogInfo(
		">>> HTTP CreatePurchase REQUEST ID=%s ExternalID=%s",
		req.ID,
		req.ExternalID,
	)

	purchaseID, err := uuid.Parse(req.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid purchase UUID")
	}

	userID, err := uuid.Parse(req.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid UUID format")
	}

	evt := events.Purchase{
		ID:             purchaseID,
		User:           userID,
		ExternalID:     req.ExternalID,
		AmountKobo:     req.AmountKobo,
		PaymentAccount: req.PaymentAccount,
	}

	utils.LogInfo(
		">>> BEFORE ProcessPurchase ID=%s ExternalID=%s",
		evt.ID,
		evt.ExternalID,
	)

	err = h.servicesAdapter.AchievementService.ProcessPurchase(
		c.Request().Context(),
		pgtype.UUID{Bytes: evt.User, Valid: true},
		evt,
	)

	utils.LogInfo(
		">>> AFTER ProcessPurchase ID=%s ExternalID=%s ERR=%v",
		evt.ID,
		evt.ExternalID,
		err,
	)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]any{
			"status": false,
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status": "Purchase made successfully",
	})
}
