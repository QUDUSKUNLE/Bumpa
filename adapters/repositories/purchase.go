package repositories

import (
	"context"

  "github.com/google/uuid"
	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
)



func (r *Repository) InsertPurchaseIfNew(ctx context.Context, purchase events.Purchase) (bool, error) {
    // convert purchase to db.CreatePurchaseParams
    // insert with ON CONFLICT DO NOTHING
    // return true if inserted
    return true, nil
}

func (r *Repository) GetPurchaseStats(ctx context.Context, userID uuid.UUID) (events.PurchaseStats, error) {
    // count purchases for this user
    return events.PurchaseStats{TotalPurchases: 3}, nil
}
