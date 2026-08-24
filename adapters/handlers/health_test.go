package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealth_Success(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	handler := &HttpHandler{}

	err := handler.Health(c)

	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]any

	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, true, response["status"])
	assert.Equal(
		t,
		map[string]any{
			"status": "Ok",
		},
		response["data"],
	)
}
