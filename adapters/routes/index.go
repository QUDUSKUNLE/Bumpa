package routes

import (
	"github.com/QUDUSKUNLE/Bumpa/adapters/handlers"
	"github.com/labstack/echo/v4"
)

func PublicRoutesAdaptor(public *echo.Group, handler *handlers.HttpHandler) *echo.Group {
	public.GET("/", handler.Home)
	public.GET("/health", handler.Health)
	public.GET("/users/:user_id/achievement", handler.GetAchievement)
	public = public.Group("/v1")
	return public
}
