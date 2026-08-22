package payments

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/QUDUSKUNLE/Bumpa/core/domain"
)

type PaymentHandler struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func (p *PaymentHandler) SendCashback(ctx context.Context, req domain.CashbackRequest) (domain.PaymentResult, error) {
	payload := map[string]any{
		"account":     req.PaymentAccount,
		"amount_kobo": req.AmountKobo,
		"reference":   req.Reference,
		"reason":      req.Reason,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return domain.PaymentResult{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.BaseURL+"/transfers", bytes.NewReader(body))
	if err != nil {
		return domain.PaymentResult{}, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.HTTPClient.Do(httpReq)
	if err != nil {
		return domain.PaymentResult{}, err
	}
	defer resp.Body.Close()

	var out struct {
		Reference string `json:"reference"`
		Status    string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return domain.PaymentResult{}, err
	}

	return domain.PaymentResult{
		ProviderReference: out.Reference,
		Status:            out.Status,
	}, nil
}
