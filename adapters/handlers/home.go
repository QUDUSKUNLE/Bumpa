package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (handler HttpHandler) Home(context echo.Context) error {
	return ComputeResponseMessage(Response{http.StatusOK, map[string]string{"home": "Bumpa Assessment"}, context})
}
