package payments

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/QUDUSKUNLE/Bumpa/core/domain"
)

type (
	PaymentHandler struct {
		BaseURL    string
		APIKey     string
		HTTPClient *http.Client
	}
	PaystackConfig struct {
		SecretKey string
		BaseURL   string
	}
	PaystackTransactionResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    struct {
			AuthorizationURL string `json:"authorization_url"`
			AccessCode       string `json:"access_code"`
			Reference        string `json:"reference"`
		} `json:"data"`
	}
	Authorization struct {
		AuthorizationCode string `json:"authorization_code"`
		CardType          string `json:"card_type"`
		Last4             string `json:"last4"`
		ExpMonth          string `json:"exp_month"`
		ExpYear           string `json:"exp_year"`
		Bank              string `json:"bank"`
		Brand             string `json:"brand"`
	}
	Data struct {
		Status        string                 `json:"status"`
		Reference     string                 `json:"reference"`
		Amount        float64                `json:"amount"`
		Channel       string                 `json:"channel"`
		Currency      string                 `json:"currency"`
		PaidAt        string                 `json:"paid_at"`
		CreatedAt     string                 `json:"created_at"`
		Customer      Customer               `json:"customer"`
		Authorization Authorization          `json:"authorization"`
		Fees          float64                `json:"fees"`
		Metadata      map[string]interface{} `json:"metadata"`
	}
	PaystackVerificationResponse struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
		Data    Data   `json:"data"`
	}
	Customer struct {
		ID           int    `json:"id"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Email        string `json:"email"`
		CustomerCode string `json:"customer_code"`
		Phone        string `json:"phone"`
	}
	PaystackAdapter struct {
		config *PaystackConfig
	}
)

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

func NewPaystackAdapter(con *PaystackConfig) *PaystackAdapter {
	return &PaystackAdapter{
		config: con,
	}
}

func (p *PaystackAdapter) InitializeTransaction(email, reference string, amount float64, metadata map[string]interface{}) (*PaystackTransactionResponse, error) {
	url := fmt.Sprintf("%s/transaction/initialize", p.config.BaseURL)

	// Paystack expects amount in kobo (multiply by 100)
	amountInKobo := int64(amount * 100)

	payload := map[string]interface{}{
		"email":     email,
		"amount":    amountInKobo,
		"reference": reference,
		"metadata":  metadata,
		"currency":  "NGN",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.SecretKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result PaystackTransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// VerifyTransaction verifies a transaction using its reference
func (p *PaystackAdapter) VerifyTransaction(reference string) (*PaystackVerificationResponse, error) {
	url := fmt.Sprintf("%s/transaction/verify/%s", p.config.BaseURL, reference)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.SecretKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result PaystackVerificationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}
