package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QUDUSKUNLE/Bumpa/core/services"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeResponseMessage_Success(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	response := Response{
		Status: http.StatusOK,
		Data: map[string]string{
			"message": "Hello Bumpa",
		},
		Context: ctx,
	}

	err := ComputeResponseMessage(response)

	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	var body map[string]any

	err = json.Unmarshal(rec.Body.Bytes(), &body)
	require.NoError(t, err)

	assert.Equal(t, true, body["status"])

	assert.Equal(
		t,
		map[string]any{
			"message": "Hello Bumpa",
		},
		body["data"],
	)
}

func TestNewHttpAdapter(t *testing.T) {
	service := services.ServicesHandler{}

	handler := NewHttpAdapter(service)

	require.NotNil(t, handler)
	assert.Equal(t, service, handler.servicesAdapter)
}
