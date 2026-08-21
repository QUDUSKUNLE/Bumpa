package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// This route is the initial API boundary for the achievements feature.
	// The domain and repository implementation can be attached without
	// changing the server bootstrap.
	e.GET("/users/:user/achievements", func(c echo.Context) error {
		if c.Param("user") == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "user is required")
		}

		return c.JSON(http.StatusOK, map[string]any{
			"unlocked_achievements":          []string{},
			"next_available_achievements":    []string{},
			"current_badge":                  nil,
			"next_badge":                     nil,
			"remaining_to_unlock_next_badge": 0,
		})
	})

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		if err := e.Start(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("start server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Printf("shutdown server: %v", err)
	}

}
