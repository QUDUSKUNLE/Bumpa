package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (handler HttpHandler) Health(context echo.Context) error {
	return ComputeResponseMessage(Response{http.StatusOK, map[string]string{"status": "Ok"}, context})
}
