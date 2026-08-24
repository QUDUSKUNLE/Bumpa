package cashback

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/QUDUSKUNLE/Bumpa/core/ports"
	"github.com/QUDUSKUNLE/Bumpa/core/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type PaymentProvider interface {
	SendCashback(ctx context.Context, req domain.CashbackRequest) (domain.PaymentResult, error)
}

type CashbackService struct {
	repo       ports.RepositoryPorts
	provider   PaymentProvider
	amountKobo int64
}

func NewCashbackService(repo ports.RepositoryPorts, provider PaymentProvider) *CashbackService {
	return &CashbackService{
		repo:       repo,
		provider:   provider,
		amountKobo: 30000, // ₦300.00
	}
}

func (s *CashbackService) HandleBadgeUnlocked(ctx context.Context, event domain.Event) error {
	utils.LogInfo("Triggered HandleBadgeUnlocked: %v", nil)
	var payload events.BadgeUnlockedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		utils.LogError("Unmarshal Service: %v", err)
		return err
	}
	user, err := s.repo.GetUser(ctx, pgtype.UUID{Bytes: event.UserID, Valid: true})
	if err != nil {
		utils.LogError("GetUser Service Error: %v", err)
		return err
	}

	if user.PaymentAccount.String == "" {
		return fmt.Errorf("user %s has no payment account configured", event.AggregateID)
	}

	payment := domain.Payment{
		ID:         pgtype.UUID{Bytes: uuid.New(), Valid: true},
		UserID:     pgtype.UUID{Bytes: event.AggregateID, Valid: true},
		BadgeCode:  payload.BadgeCode,
		AmountKobo: s.amountKobo,
		Status:     "pending",
		CreatedAt:  time.Now().UTC(),
	}

	created, err := s.repo.CreatePaymentIfNew(ctx, payment)
	if err != nil {
		utils.LogError("CreatePaymentIfNew Service Error: %v", err)
		return err
	}
	if !created {
		return nil
	}

	result, err := s.provider.SendCashback(ctx, domain.CashbackRequest{
		UserID:         pgtype.UUID{Bytes: event.AggregateID, Valid: true},
		PaymentAccount: user.PaymentAccount.String,
		AmountKobo:     s.amountKobo,
		Reference:      fmt.Sprintf("cashback-%s-%s", payload.BadgeCode, event.AggregateID),
		Reason:         "Badge cashback",
	})
	if err != nil {
		utils.LogError("SendCashback Service Error: %v", err)
		return s.repo.MarkPaymentFailed(ctx, payment.ID, err.Error())
	}

	utils.LogInfo("Finished HandleBadgeUnlocked: %v", nil)
	return s.repo.MarkPaymentSuccessful(ctx, payment.ID, result.ProviderReference)
}
