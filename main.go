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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

func seedUser(ctx context.Context, repo *repositories.Repository) error {
	user := db.CreateUserParams{
		Name:  "Qudus Yekeen Adekunle335",
		Email: "qudus.adekunle@example.com",
		Phone: pgtype.Text{String: "+23480000000001", Valid: true},
		PaymentAccount: pgtype.Text{
			String: "RCP_m7ljkv8leesep7p",
			Valid:  true,
		},
	}

	_, err := repo.CreateUser(ctx, user)
	return err
}

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

	if err := seedUser(ctx, repo); err != nil {
		log.Printf("seed user: %v", err)
	}

	// Event bus
	bus := events.NewEventBus()

	service := services.NewServiceAdapter(repo, *bus)
	httpHandler := handlers.NewHttpAdapter(*service)

	badgeService := badges.NewBadgeService(repo, badges.BadgeDefinition(), *bus)

	// Paystack Service
	paymentService := payments.NewPaystackAdapter(&payments.PaystackConfig{
		SecretKey:  cfg.PAYSTACK_SECRET_KEY,
		BaseURL:    cfg.PAYSTACK_BASE_URL,
		HTTPClient: &http.Client{},
	})

	// paymentProvider := payments.NewPaystackAdapter()
	cashbackSvc := cashback.NewCashbackService(repo, *paymentService)

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
		func(
			ctx context.Context,
			evt domain.Event,
		) error {
			utils.LogInfo(
				"Triggered Subscribe BadgeUnlocked: eventID=%s userID=%s",
				evt.ID,
				evt.UserID,
			)
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
