package ports

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

type RepositoryPorts interface {
	WithTx(ctx context.Context, fn func(tx RepositoryPorts) error) error

	CreateUser(ctx context.Context, user db.CreateUserParams) (db.User, error)
	GetUser(ctx context.Context, id pgtype.UUID) (db.GetUserRow, error)

	InsertPurchaseIfNew(ctx context.Context, purchase events.Purchase) (bool, error)

	GetPurchaseStats(ctx context.Context, userID pgtype.UUID) (events.PurchaseStats, error)
	UnlockAchievementIfNew(ctx context.Context, userID pgtype.UUID, achievementCode string) (bool, error)
	GetUnlockedAchievementCount(ctx context.Context, userID pgtype.UUID) (int, error)
	UnlockBadgeIfNew(ctx context.Context, userID pgtype.UUID, badgeCode string) (bool, error)

	AddOutboxEvent(ctx context.Context, event domain.Event) error

	CreatePaymentIfNew(ctx context.Context, payment domain.Payment) (bool, error)
	MarkPaymentSuccessful(ctx context.Context, id pgtype.UUID, providerRef string) error
	GetPaymentByUserAndBadge(
		ctx context.Context,
		userID pgtype.UUID,
		badgeCode string,
	) (*domain.Payment, error)
	MarkPaymentFailed(ctx context.Context, id pgtype.UUID, reason string) error
	GetPendingOutboxEvents(ctx context.Context, batchSize int) ([]db.OutboxEvent, error)
	GetUserAchievements(
		ctx context.Context,
		userID pgtype.UUID,
	) (db.GetUserAchievementsRow, error)

	MarkOutboxEventProcessed(
		ctx context.Context,
		id pgtype.UUID,
	) error
}
