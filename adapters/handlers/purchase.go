package handlers

import (
	"net/http"

	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func (h *HttpHandler) CreatePurchase(c echo.Context) error {
	var req events.PurchaseRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"status": false,
			"error":  err.Error(),
		})
	}

	utils.LogInfo(
		"HTTP CreatePurchase REQUEST ID=%s ExternalID=%s",
		req.ID,
		req.ExternalID,
	)

	purchaseID, err := uuid.Parse(req.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid purchase UUID")
	}

	userID, err := uuid.Parse(req.User)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid UUID format")
	}

	purchaseEvent := events.Purchase{
		ID:             purchaseID,
		User:           userID,
		ExternalID:     req.ExternalID,
		AmountKobo:     req.AmountKobo,
		PaymentAccount: req.PaymentAccount,
	}

	utils.LogInfo(
		"BEFORE ProcessPurchase ID=%s ExternalID=%s",
		purchaseEvent.ID,
		purchaseEvent.ExternalID,
	)

	err = h.servicesAdapter.AchievementService.ProcessPurchase(
		c.Request().Context(),
		pgtype.UUID{Bytes: purchaseEvent.User, Valid: true},
		purchaseEvent,
	)

	utils.LogInfo(
		"AFTER ProcessPurchase ID=%s ExternalID=%s ERROR=%v",
		purchaseEvent.ID,
		purchaseEvent.ExternalID,
		err,
	)

	if err != nil {
		return c.JSON(http.StatusConflict, map[string]any{
			"status": false,
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status": "Purchase made successfully",
	})
}
