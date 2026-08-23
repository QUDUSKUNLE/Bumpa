package ports

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/google/uuid"
)

type RepositoryPorts interface {
	WithTx(ctx context.Context, fn func(tx RepositoryPorts) error) error

	CreateUser(ctx context.Context, user db.CreateUserParams) (db.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (db.GetUserRow, error)

	InsertPurchaseIfNew(ctx context.Context, purchase events.Purchase) (bool, error)

	GetPurchaseStats(ctx context.Context, userID uuid.UUID) (events.PurchaseStats, error)
	UnlockAchievementIfNew(ctx context.Context, userID uuid.UUID, achievementCode string) (bool, error)
	GetUnlockedAchievementCount(ctx context.Context, userID uuid.UUID) (int, error)
	UnlockBadgeIfNew(ctx context.Context, userID uuid.UUID, badgeCode string) (bool, error)

	AddOutboxEvent(ctx context.Context, event domain.Event) error

	CreatePaymentIfNew(ctx context.Context, payment domain.Payment) (bool, error)
	MarkPaymentSuccessful(ctx context.Context, id uuid.UUID, providerRef string) error
	MarkPaymentFailed(ctx context.Context, id uuid.UUID, reason string) error
}
