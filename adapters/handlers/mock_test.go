package handlers

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

type MockAchievementService struct {
	ProcessPurchaseFunc func(
		ctx context.Context,
		userID pgtype.UUID,
		purchase events.Purchase,
	) error

	ProcessFunc func(
		ctx context.Context,
	) error

	GetUserAchievementsFunc func(
		ctx echo.Context,
	) error
}

func (m *MockAchievementService) GetUserAchievements(
	ctx echo.Context,
) error {
	if m.GetUserAchievementsFunc != nil {
		return m.GetUserAchievementsFunc(ctx)
	}

	return nil
}

func (m *MockAchievementService) ProcessPurchase(
	ctx context.Context,
	userID pgtype.UUID,
	purchase events.Purchase,
) error {
	if m.ProcessPurchaseFunc != nil {
		return m.ProcessPurchaseFunc(ctx, userID, purchase)
	}

	return nil
}

func (m *MockAchievementService) Process(
	ctx context.Context,
) error {
	if m.ProcessFunc != nil {
		return m.ProcessFunc(ctx)
	}

	return nil
}
