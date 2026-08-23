package handlers

import (
	"github.com/QUDUSKUNLE/Bumpa/core/services"
	"github.com/labstack/echo/v4"
)

type HttpHandler struct {
	servicesAdapter services.ServicesHandler
	// bus             events.EventBus
}

func NewHttpAdapter(service services.ServicesHandler) *HttpHandler {
	return &HttpHandler{
		servicesAdapter: service,
		// bus:             *events.NewEventBus(),
	}
}

func ComputeResponseMessage(response Response) error {
	return response.Context.JSON(
		response.Status, echo.Map{
			"data":   response.Data,
			"status": true,
		})
}

type ValidationStruct struct {
	Context   echo.Context `json:"context"`
	Interface interface{}  `json:"interface"`
}

type Response struct {
	Status  int          `json:"status"`
	Data    interface{}  `json:"data"`
	Context echo.Context `json:"context"`
}
