package payments

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/QUDUSKUNLE/Bumpa/core/utils"
)

type (
	PaystackConfig struct {
		SecretKey  string
		BaseURL    string
		HTTPClient *http.Client
	}
	Data struct {
		TransferSessionID []any                  `json:"transfersessionid"`
		TransferTrials    []any                  `json:"transfertrials"`
		Domain            string                 `json:"domain"`
		Amount            float64                `json:"amount"`
		Currency          string                 `json:"currency"`
		Reference         string                 `json:"reference"`
		Source            string                 `json:"source"`
		Reason            string                 `json:"reason"`
		Status            string                 `json:"status"`
		TransferCode      string                 `json:"transfer_code"`
		ID                int                    `json:"id"`
		Integration       int                    `json:"integration"`
		Request           int                    `json:"request"`
		Recipient         int                    `json:"recipient"`
		CreatedAt         string                 `json:"created_at"`
		UpdatedAt         string                 `json:"updated_at"`
		Metadata          map[string]interface{} `json:"metadata"`
	}
	PaystackResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    Data   `json:"data"`
	}
	PaystackAdapter struct {
		config *PaystackConfig
	}
)

func (p *PaystackAdapter) SendCashback(
	ctx context.Context,
	req domain.CashbackRequest,
) (PaystackResponse, error) {

	payload := map[string]any{
		"source":    req.Source,
		"amount":    req.AmountKobo,
		"recipient": req.PaymentAccount,
		"reference": req.Reference,
		"reason":    req.Reason,
		"currency":  req.Currency,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return PaystackResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		p.config.BaseURL+"/transfer",
		bytes.NewReader(body),
	)
	if err != nil {
		return PaystackResponse{}, err
	}

	httpReq.Header.Set(
		"Authorization",
		"Bearer "+p.config.SecretKey,
	)
	httpReq.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := p.config.HTTPClient.Do(httpReq)
	if err != nil {
		return PaystackResponse{}, err
	}
	defer resp.Body.Close()

	// Read the complete response body first.
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return PaystackResponse{}, err
	}

	utils.LogInfo(
		"Paystack response: status=%d body=%s",
		resp.StatusCode,
		string(responseBody),
	)

	// HTTP error from Paystack.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return PaystackResponse{}, fmt.Errorf(
			"paystack returned HTTP %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	// Successful HTTP response but empty body.
	if len(responseBody) == 0 {
		return PaystackResponse{}, fmt.Errorf(
			"paystack returned HTTP %d with empty response body",
			resp.StatusCode,
		)
	}

	var out PaystackResponse

	if err := json.Unmarshal(responseBody, &out); err != nil {
		return PaystackResponse{}, fmt.Errorf(
			"decode paystack response: %w; body=%s",
			err,
			string(responseBody),
		)
	}

	return out, nil
}

func (p *PaystackAdapter) FinaliseCashBack(ctx context.Context, req domain.FinaliseCashBackRequest) (PaystackResponse, error) {
	return PaystackResponse{}, nil
}

func NewPaystackAdapter(con *PaystackConfig) *PaystackAdapter {
	return &PaystackAdapter{
		config: con,
	}
}
