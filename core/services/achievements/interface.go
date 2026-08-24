package achievements

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

type Service interface {
	ProcessPurchase(
		ctx context.Context,
		userID pgtype.UUID,
		purchase events.Purchase,
	) error

	Process(ctx context.Context) error

	GetUserAchievements(ctx echo.Context) error
}
