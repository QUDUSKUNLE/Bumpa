package routes

import (
	"testing"

	"github.com/QUDUSKUNLE/Bumpa/adapters/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublicRoutesAdaptor(t *testing.T) {
	e := echo.New()

	public := e.Group("/api")

	handler := &handlers.HttpHandler{}

	result := PublicRoutesAdaptor(public, handler)

	require.NotNil(t, result)
	assert.Same(t, public, result)

	routes := e.Routes()

	expectedRoutes := map[string]string{
		"GET /api/":                         "Home",
		"GET /api/health":                   "Health",
		"GET /api/users/:user/achievements": "GetUserAchievements",
		"POST /api/users/purchases":         "CreatePurchase",
	}

	for _, route := range routes {
		key := route.Method + " " + route.Path

		if expectedHandler, ok := expectedRoutes[key]; ok {
			assert.Contains(
				t,
				route.Name,
				expectedHandler,
				"route %s should use handler %s",
				key,
				expectedHandler,
			)
		}
	}

	assert.Len(t, routes, len(expectedRoutes))
}
