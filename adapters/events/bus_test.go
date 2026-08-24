package events

import (
	"context"
	"errors"
	"testing"

	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEventBus(t *testing.T) {
	bus := NewEventBus()

	require.NotNil(t, bus)
	require.NotNil(t, bus.handlers)
	assert.Empty(t, bus.handlers)
}

func TestEventBus_Subscribe(t *testing.T) {
	bus := NewEventBus()

	handler := func(ctx context.Context, event domain.Event) error {
		return nil
	}

	bus.Subscribe("AchievementUnlocked", handler)

	require.Len(t, bus.handlers["AchievementUnlocked"], 1)
	assert.NotNil(t, bus.handlers["AchievementUnlocked"][0])
}

func TestEventBus_Publish_CallsSubscribedHandler(t *testing.T) {
	bus := NewEventBus()

	called := false

	bus.Subscribe(
		"AchievementUnlocked",
		func(ctx context.Context, event domain.Event) error {
			called = true
			return nil
		},
	)

	event := domain.Event{
		Type: "AchievementUnlocked",
	}

	err := bus.Publish(context.Background(), event)

	require.NoError(t, err)
	assert.True(t, called)
}

func TestEventBus_Publish_PassesContextAndEvent(t *testing.T) {
	bus := NewEventBus()

	expectedContext := context.WithValue(
		context.Background(),
		"test-key",
		"test-value",
	)

	expectedEvent := domain.Event{
		Type: "AchievementUnlocked",
	}

	var receivedContext context.Context
	var receivedEvent domain.Event

	bus.Subscribe(
		"AchievementUnlocked",
		func(ctx context.Context, event domain.Event) error {
			receivedContext = ctx
			receivedEvent = event
			return nil
		},
	)

	err := bus.Publish(expectedContext, expectedEvent)

	require.NoError(t, err)

	assert.Equal(t, expectedContext, receivedContext)
	assert.Equal(t, expectedEvent, receivedEvent)
}

func TestEventBus_Publish_ReturnsHandlerError(t *testing.T) {
	bus := NewEventBus()

	expectedErr := errors.New("handler failed")

	bus.Subscribe(
		"AchievementUnlocked",
		func(ctx context.Context, event domain.Event) error {
			return expectedErr
		},
	)

	event := domain.Event{
		Type: "AchievementUnlocked",
	}

	err := bus.Publish(context.Background(), event)

	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
}

func TestEventBus_Publish_DoesNotCallHandlerForDifferentEventType(t *testing.T) {
	bus := NewEventBus()

	called := false

	bus.Subscribe(
		"AchievementUnlocked",
		func(ctx context.Context, event domain.Event) error {
			called = true
			return nil
		},
	)

	event := domain.Event{
		Type: "BadgeUnlocked",
	}

	err := bus.Publish(context.Background(), event)

	require.NoError(t, err)
	assert.False(t, called)
}

func TestEventBus_Publish_CallsMultipleHandlers(t *testing.T) {
	bus := NewEventBus()

	firstCalled := false
	secondCalled := false

	bus.Subscribe(
		"AchievementUnlocked",
		func(ctx context.Context, event domain.Event) error {
			firstCalled = true
			return nil
		},
	)

	bus.Subscribe(
		"AchievementUnlocked",
		func(ctx context.Context, event domain.Event) error {
			secondCalled = true
			return nil
		},
	)

	event := domain.Event{
		Type: "AchievementUnlocked",
	}

	err := bus.Publish(context.Background(), event)

	require.NoError(t, err)

	assert.True(t, firstCalled)
	assert.True(t, secondCalled)
}

func TestEventBus_Publish_StopsOnHandlerError(t *testing.T) {
	bus := NewEventBus()

	expectedErr := errors.New("first handler failed")
	secondCalled := false

	bus.Subscribe(
		"AchievementUnlocked",
		func(ctx context.Context, event domain.Event) error {
			return expectedErr
		},
	)

	bus.Subscribe(
		"AchievementUnlocked",
		func(ctx context.Context, event domain.Event) error {
			secondCalled = true
			return nil
		},
	)

	event := domain.Event{
		Type: "AchievementUnlocked",
	}

	err := bus.Publish(context.Background(), event)

	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.False(t, secondCalled)
}
