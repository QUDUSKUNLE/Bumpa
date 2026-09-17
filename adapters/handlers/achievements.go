package handlers

import (
	"github.com/labstack/echo/v4"
)

func (handler *HttpHandler) GetUserAchievements(ctx echo.Context) error {
	return handler.servicesAdapter.AchievementService.GetUserAchievements(ctx)
}
