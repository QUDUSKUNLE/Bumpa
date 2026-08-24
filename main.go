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
	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/adapters/handlers"
	"github.com/QUDUSKUNLE/Bumpa/adapters/payments"
	"github.com/QUDUSKUNLE/Bumpa/adapters/repositories"
	"github.com/QUDUSKUNLE/Bumpa/adapters/routes"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/QUDUSKUNLE/Bumpa/core/services"
	"github.com/QUDUSKUNLE/Bumpa/core/services/badges"
	"github.com/QUDUSKUNLE/Bumpa/core/services/cashback"
	"github.com/QUDUSKUNLE/Bumpa/core/services/outboxprocessor"
	"github.com/QUDUSKUNLE/Bumpa/core/utils"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

// func seedUser(ctx context.Context, repo *repositories.Repository) error {
// 	user := db.CreateUserParams{
// 		Name:  "Jane Doe",
// 		Email: "jane@example.com",
// 		Phone: pgtype.Text{String: "+2348000000000", Valid: true},
// 		PaymentAccount: pgtype.Text{
// 			String: "acct_123456",
// 			Valid:  true,
// 		},
// 	}

// 	_, err := repo.CreateUser(ctx, user)
// 	return err
// }

func main() {
	// Load configuration
	cfg, err := config.LoadEnvironmentVariables()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Application context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database with custom configuration
	dbConfig := config.DBConfig()

	store, conn, err := db.DatabaseConnection(
		ctx,
		cfg.DATABASE_URL,
		dbConfig,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	repo := repositories.NewRepository(store, conn)

	// if err := seedUser(context.Background(), repo); err != nil {
	// 	log.Printf("seed user: %v", err)
	// }

	// Event bus
	bus := events.NewEventBus()

	service := services.NewServiceAdapter(repo, *bus)
	httpHandler := handlers.NewHttpAdapter(*service)

	badgeService := badges.NewBadgeService(repo, badges.BadgeDefinition(), *bus)

	paymentProvider := &payments.PaymentHandler{
		BaseURL:    "https://api.example.com",
		APIKey:     "demo-key",
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
	cashbackSvc := cashback.NewCashbackService(repo, paymentProvider)

	bus.Subscribe(
		"AchievementUnlocked",
		func(ctx context.Context, evt domain.Event) error {
			utils.LogInfo("Triggered Subscribe AchievementUnlocked: %v", nil)
			return badgeService.HandleAchievementUnlocked(ctx, events.Event{
				ID:          evt.ID,
				UserID:      evt.UserID,
				Type:        evt.Type,
				OccurredAt:  evt.OccurredAt,
				AggregateID: evt.AggregateID,
				Payload:     evt.Payload,
			})
		})

	bus.Subscribe(
		"BadgeUnlocked",
		func(ctx context.Context, evt domain.Event) error {
			utils.LogInfo("Triggered Subscribe BadgeUnlocked: %v", nil)
			return cashbackSvc.HandleBadgeUnlocked(ctx, evt)
		},
	)

	// --------------------------------------------------
	// Outbox processor
	// --------------------------------------------------
	outboxProcessor := outboxprocessor.NewOutboxProcessor(repo, *bus)
	go outboxProcessor.Run(ctx)

	// --------------------------------------------------
	// HTTP server
	// --------------------------------------------------

	e := echo.New()

	e.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodOptions,
			},
			AllowHeaders: []string{
				echo.HeaderContentType,
				echo.HeaderAuthorization,
			},
		},
	))

	routes.PublicRoutesAdaptor(
		e.Group(""),
		httpHandler,
	)

	if cfg.HTTP_PORT == "" {
		cfg.HTTP_PORT = "8081"
	}

	go func() {
		if err := e.Start(":" + cfg.HTTP_PORT); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("start server: %v", err)
		}
	}()

	// --------------------------------------------------
	// Wait for shutdown signal
	// --------------------------------------------------

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutdown signal received")

	cancel()

	// Give HTTP server time to finish requests
	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer shutdownCancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown server: %v", err)
	}

	log.Println("Application stopped")
}
