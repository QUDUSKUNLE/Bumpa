package routes

import (
	"github.com/QUDUSKUNLE/Bumpa/adapters/handlers"
	"github.com/labstack/echo/v4"
)

func PublicRoutesAdaptor(public *echo.Group, handler *handlers.HttpHandler) *echo.Group {
	public.GET("/", handler.Home)
	public.GET("/health", handler.Health)
	public.GET("/users/:user/achievements", handler.GetUserAchievements)
	public.POST("/users/purchases", handler.CreatePurchase)
	return public
}
