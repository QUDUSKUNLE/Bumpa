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

	"github.com/QUDUSKUNLE/Bumpa/adapters/config"
	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/adapters/handlers"
	"github.com/QUDUSKUNLE/Bumpa/adapters/routes"
	"github.com/QUDUSKUNLE/Bumpa/adapters/repositories"
	"github.com/QUDUSKUNLE/Bumpa/core/services"

	"github.com/labstack/echo/v4"
)

func main() {

	// Load configuration
	cfg, err := config.LoadEnvironmentVariables()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize database with custom configuration
	dbConfig := config.DBConfig()

	store, conn, err := db.DatabaseConnection(
		context.Background(),
		cfg.DATABASE_URL,
		dbConfig,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	repo := repositories.NewRepository(store, conn)
	service := services.NewServiceAdapter(repo)
	httpHandler := handlers.NewHttpAdapter(*service)

	e := echo.New()

	// Plug echo into PublicRoutesAdaptor
	public := e.Group("")

	routes.PublicRoutesAdaptor(public, httpHandler)

	if cfg.HTTP_PORT == "" {
		cfg.HTTP_PORT = "8080"
	}

	go func() {
		if err := e.Start(":" + cfg.HTTP_PORT); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
