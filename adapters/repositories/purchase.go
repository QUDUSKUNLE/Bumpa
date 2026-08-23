package repositories

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/adapters/db"
	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *Repository) InsertPurchaseIfNew(ctx context.Context, req events.Purchase) (bool, error) {
	// convert purchase to db.CreatePurchaseParams
	// insert with ON CONFLICT DO NOTHING
	// return true if inserted
	_, err := r.queries.CreatePurchase(ctx, db.CreatePurchaseParams{
		ID:         pgtype.UUID{Bytes: uuid.New(), Valid: true},
		UserID:     pgtype.UUID{Bytes: req.User, Valid: true},
		ExternalID: req.ExternalID,
		AmountKobo: int64(req.AmountKobo)})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Repository) GetPurchaseStats(ctx context.Context, userID uuid.UUID) (events.PurchaseStats, error) {
	// count purchases for this user
	count, err := r.queries.CountPurchasesByUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return events.PurchaseStats{}, err
	}
	return events.PurchaseStats{TotalPurchases: int(count)}, nil
}
