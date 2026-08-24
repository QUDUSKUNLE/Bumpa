package repositories

import (
	"context"
	"errors"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *Repository) InsertPurchaseIfNew(ctx context.Context, req events.Purchase) (bool, error) {
	utils.LogInfo(
		"=== INSERT PURCHASE START ID=%v ===",
		req.ID,
	)

	_, err := r.queries.CreatePurchase(ctx, db.CreatePurchaseParams{
		ID:         pgtype.UUID{Bytes: req.ID, Valid: true},
		UserID:     pgtype.UUID{Bytes: req.User, Valid: true},
		ExternalID: req.ExternalID,
		AmountKobo: int64(req.AmountKobo)})

	utils.LogInfo(
		"=== INSERT PURCHASE RESULT ID=%v ERR=%v ===",
		req.ID,
		err,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		utils.LogInfo("Purchase already exists, skipping: %v", req.ID)
		return false, nil
	}
	if err != nil {
		utils.LogError("CreatePurchase Repositories Error: %v", err)
		return false, err
	}
	return true, nil
}

func (r *Repository) GetPurchaseStats(ctx context.Context, userID pgtype.UUID) (events.PurchaseStats, error) {
	// count purchases for this user
	count, err := r.queries.CountPurchasesByUser(ctx, userID)
	if err != nil {
		utils.LogError("CountPurchasesByUser Repositories: %v", err)
		return events.PurchaseStats{}, err
	}
	return events.PurchaseStats{TotalPurchases: int(count)}, nil
}
