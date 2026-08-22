package cashback

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/QUDUSKUNLE/Bumpa/core/ports"
	"github.com/google/uuid"
)

type PaymentProvider interface {
	SendCashback(ctx context.Context, req domain.CashbackRequest) (domain.PaymentResult, error)
}

type CashbackService struct {
	repo       ports.Repository
	provider   PaymentProvider
	amountKobo int64
}

func NewCashbackService(repo ports.Repository, provider PaymentProvider) *CashbackService {
	return &CashbackService{
		repo:       repo,
		provider:   provider,
		amountKobo: 30000, // ₦300.00
	}
}

func (s *CashbackService) HandleBadgeUnlocked(ctx context.Context, event domain.Event) error {
	var payload events.BadgeUnlockedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}
	user, err := s.repo.GetUser(ctx, event.AggregateID)
	if err != nil {
		return err
	}

	if user.PaymentAccount == "" {
		return fmt.Errorf("user %s has no payment account configured", event.AggregateID)
	}

	payment := domain.Payment{
		ID:         uuid.New(),
		UserID:     event.AggregateID,
		BadgeCode:  payload.BadgeCode,
		AmountKobo: s.amountKobo,
		Status:     "pending",
		CreatedAt:  time.Now().UTC(),
	}

	created, err := s.repo.CreatePaymentIfNew(ctx, payment)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}

	result, err := s.provider.SendCashback(ctx, domain.CashbackRequest{
		UserID:         event.AggregateID,
		PaymentAccount: user.PaymentAccount,
		AmountKobo:     s.amountKobo,
		Reference:      fmt.Sprintf("cashback-%s-%s", payload.BadgeCode, event.AggregateID),
		Reason:         "Badge cashback",
	})
	if err != nil {
		return s.repo.MarkPaymentFailed(ctx, payment.ID, err.Error())
	}

	return s.repo.MarkPaymentSuccessful(ctx, payment.ID, result.ProviderReference)
}
