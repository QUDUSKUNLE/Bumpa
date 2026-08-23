package handlers

import (
	"net/http"

	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/labstack/echo/v4"
)

func (handler *HttpHandler) GetAchievement(ctx echo.Context) error {
	if ctx.Param("user_id") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id is required")
	}
	return ComputeResponseMessage(Response{http.StatusOK, map[string]any{
		"unlocked_achievements":          []string{},
		"next_available_achievements":    []string{},
		"current_badge":                  nil,
		"next_badge":                     nil,
		"remaining_to_unlock_next_badge": 0,
	}, ctx})
}

func (h *HttpHandler) CreatePurchase(c echo.Context) error {
	var req events.Purchase
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.servicesAdapter.AchievementService.ProcessPurchase(c.Request().Context(), req.User, req); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "true",
	})
}
