package events

import (
	"context"

	"github.com/QUDUSKUNLE/Bumpa/core/domain"
)

type Handler func(context.Context, domain.Event) error

type EventPublisher interface {
	Publish(ctx context.Context, event domain.Event) error
}

type EventBus struct {
	handlers map[string][]Handler
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]Handler),
	}
}

func (b *EventBus) Subscribe(eventType string, handler Handler) {
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *EventBus) Publish(ctx context.Context, event domain.Event) error {
	for _, handler := range b.handlers[event.Type] {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}
	return nil
}
