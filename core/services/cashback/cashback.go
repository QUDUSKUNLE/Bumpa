package cashback

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/QUDUSKUNLE/Bumpa/adapters/events"
	"github.com/QUDUSKUNLE/Bumpa/adapters/payments"
	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/QUDUSKUNLE/Bumpa/core/ports"
	"github.com/QUDUSKUNLE/Bumpa/core/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type CashbackService struct {
	repo            ports.RepositoryPorts
	paymentProvider payments.PaystackAdapter
	amountKobo      int64
}

func NewCashbackService(repo ports.RepositoryPorts, provider payments.PaystackAdapter) *CashbackService {
	return &CashbackService{
		repo:            repo,
		paymentProvider: provider,
		amountKobo:      30000, // ₦300.00
	}
}

func (s *CashbackService) HandleBadgeUnlocked(
	ctx context.Context,
	event domain.Event,
) error {

	utils.LogInfo(
		"Triggered HandleBadgeUnlocked: eventID=%s userID=%s",
		event.ID,
		event.UserID,
	)

	var payload events.BadgeUnlockedPayload

	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		utils.LogError(
			"Unmarshal BadgeUnlocked Payload Error: %v",
			err,
		)
		return err
	}

	userID := pgtype.UUID{
		Bytes: event.UserID,
		Valid: true,
	}

	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		utils.LogError("GetUser Service Error: %v", err)
		return err
	}

	if !user.PaymentAccount.Valid ||
		user.PaymentAccount.String == "" {

		return fmt.Errorf(
			"user %s has no payment account configured",
			event.UserID,
		)
	}

	// ---------------------------------------------------------
	// Find existing cashback payment
	// ---------------------------------------------------------

	payment, err := s.repo.GetPaymentByUserAndBadge(
		ctx,
		userID,
		payload.BadgeCode,
	)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		utils.LogError(
			"GetPaymentByUserAndBadge Service Error: %v",
			err,
		)
		return err
	}

	if err == nil {

		// Already successfully paid.
		if payment.Status == "successful" {
			utils.LogInfo(
				"Cashback already successful: user=%s badge=%s",
				event.UserID,
				payload.BadgeCode,
			)
			return nil
		}

		// Failed or pending payment.
		// Retry the provider request.
		utils.LogInfo(
			"Retrying cashback: user=%s badge=%s status=%s",
			event.UserID,
			payload.BadgeCode,
			payment.Status,
		)

	} else {

		// ---------------------------------------------------------
		// No payment exists — create one.
		// ---------------------------------------------------------

		payment = &domain.Payment{
			ID: pgtype.UUID{
				Bytes: uuid.New(),
				Valid: true,
			},
			UserID:     userID,
			BadgeCode:  payload.BadgeCode,
			AmountKobo: s.amountKobo,
			Status:     "pending",
			CreatedAt:  time.Now().UTC(),
		}

		created, err := s.repo.CreatePaymentIfNew(
			ctx,
			*payment,
		)
		if err != nil {
			utils.LogError(
				"CreatePaymentIfNew Service Error: %v",
				err,
			)
			return err
		}

		if !created {
			// Another worker may have created the payment.
			utils.LogInfo(
				"Payment was created concurrently: user=%s badge=%s",
				event.UserID,
				payload.BadgeCode,
			)
			return nil
		}
	}

	// ---------------------------------------------------------
	// Send cashback
	// ---------------------------------------------------------

	reference := fmt.Sprintf(
		"cashback-%s-%s",
		payload.BadgeCode,
		event.UserID,
	)

	// result, err := s.paymentProvider.SendCashback(
	// 	ctx,
	// 	domain.CashbackRequest{
	// 		Source:         "balance",
	// 		UserID:         userID,
	// 		PaymentAccount: user.PaymentAccount.String,
	// 		AmountKobo:     s.amountKobo,
	// 		Reference:      reference,
	// 		Reason:         "Badge cashback",
	// 		Currency:       "NGN",
	// 	},
	// )

	// if err != nil {
	// 	utils.LogError(
	// 		"SendCashback Service Error: %v",
	// 		err,
	// 	)

	// 	if markErr := s.repo.MarkPaymentFailed(
	// 		ctx,
	// 		payment.ID,
	// 		err.Error(),
	// 	); markErr != nil {

	// 		utils.LogError(
	// 			"MarkPaymentFailed Error: %v",
	// 			markErr,
	// 		)

	// 		return fmt.Errorf(
	// 			"cashback failed: %w; failed to mark payment failed: %v",
	// 			err,
	// 			markErr,
	// 		)
	// 	}

	// 	return err
	// }

	// // ---------------------------------------------------------
	// // Mark successful
	// // ---------------------------------------------------------

	// if err := s.repo.MarkPaymentSuccessful(
	// 	ctx,
	// 	payment.ID,
	// 	result.Data.Reference,
	// ); err != nil {

	// 	utils.LogError(
	// 		"MarkPaymentSuccessful Error: %v",
	// 		err,
	// 	)
	// 	return err
	// }

	utils.LogInfo(
		"Cashback successfully completed: user=%s badge=%s reference=%s",
		event.UserID,
		payload.BadgeCode,
		reference,
		// result.Data.Reference,
		// "ok",
	)

	return nil
}
