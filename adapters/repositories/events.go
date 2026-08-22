package repositories

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
)

func (r *Repository) AddOutboxEvent(ctx context.Context, event events.Event) error {
	// insert event into outbox_events
	return nil
}
