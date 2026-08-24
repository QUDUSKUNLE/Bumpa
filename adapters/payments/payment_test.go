package payments

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QUDUSKUNLE/Bumpa/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPaystackAdapter(t *testing.T) {
	secretKey := "sk_test_xxxxx"
	baseURL := "https://api.paystack.co"

	config := &PaystackConfig{
		SecretKey:  secretKey,
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}

	adapter := NewPaystackAdapter(config)

	require.NotNil(t, adapter)
	assert.Equal(t, config, adapter.config)
}

func TestPaystackAdapter_SendCashback_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/transfer", r.URL.Path)

		assert.Equal(
			t,
			"Bearer sk_test_xxxxx",
			r.Header.Get("Authorization"),
		)

		assert.Equal(
			t,
			"application/json",
			r.Header.Get("Content-Type"),
		)

		var payload map[string]interface{}

		err := json.NewDecoder(r.Body).Decode(&payload)
		require.NoError(t, err)

		assert.Equal(t, "balance", payload["source"])
		assert.Equal(t, float64(30000), payload["amount"])
		assert.Equal(t, "RCP_xxxxx", payload["recipient"])
		assert.Equal(t, "CASHBACK-001", payload["reference"])
		assert.Equal(t, "Cashback reward", payload["reason"])
		assert.Equal(t, "NGN", payload["currency"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err = w.Write([]byte(`{
			"status": true,
			"message": "Transfer has been queued",
			"data": {
				"domain": "test",
				"amount": 30000,
				"currency": "NGN",
				"reference": "CASHBACK-001",
				"source": "balance",
				"reason": "Cashback reward",
				"status": "pending",
				"transfer_code": "TRF_xxxxx",
				"id": 12345,
				"integration": 100,
				"request": 200,
				"recipient": 300,
				"created_at": "2026-08-24T10:00:00Z",
				"updated_at": "2026-08-24T10:00:00Z"
			}
		}`))

		require.NoError(t, err)
	}))

	defer server.Close()

	secretKey := "sk_test_xxxxx"

	adapter := NewPaystackAdapter(&PaystackConfig{
		SecretKey:  secretKey,
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	})

	request := domain.CashbackRequest{
		Source:         "balance",
		AmountKobo:     30000,
		PaymentAccount: "RCP_xxxxx",
		Reference:      "CASHBACK-001",
		Reason:         "Cashback reward",
		Currency:       "NGN",
	}

	response, err := adapter.SendCashback(
		context.Background(),
		request,
	)

	require.NoError(t, err)

	assert.NotNil(t, response.Status)
	assert.True(t, response.Status)

	assert.NotNil(t, response.Message)
	assert.Equal(
		t,
		"Transfer has been queued",
		response.Message,
	)

	assert.NotNil(t, response.Data.Amount)
	assert.Equal(t, float64(30000), response.Data.Amount)

	assert.NotNil(t, response.Data.Currency)
	assert.Equal(t, "NGN", response.Data.Currency)

	assert.NotNil(t, response.Data.Reference)
	assert.Equal(t, "CASHBACK-001", response.Data.Reference)

	assert.NotNil(t, response.Data.Status)
	assert.Equal(t, "pending", response.Data.Status)

	assert.NotNil(t, response.Data.TransferCode)
	assert.Equal(t, "TRF_xxxxx", response.Data.TransferCode)

	assert.NotNil(t, response.Data.ID)
	assert.Equal(t, 12345, response.Data.ID)
}

func TestPaystackAdapter_SendCashback_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		_, err := w.Write([]byte(`{
			"status": false,
			"message": "Invalid recipient"
		}`))

		require.NoError(t, err)
	}))

	defer server.Close()

	secretKey := "sk_test_xxxxx"

	adapter := NewPaystackAdapter(&PaystackConfig{
		SecretKey:  secretKey,
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	})

	request := domain.CashbackRequest{
		Source:         "balance",
		AmountKobo:     30000,
		PaymentAccount: "INVALID",
		Reference:      "CASHBACK-001",
		Reason:         "Cashback reward",
		Currency:       "NGN",
	}

	_, err := adapter.SendCashback(
		context.Background(),
		request,
	)

	require.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"paystack returned HTTP 400",
	)

	assert.Contains(
		t,
		err.Error(),
		"Invalid recipient",
	)
}

func TestPaystackAdapter_SendCashback_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	secretKey := "sk_test_xxxxx"

	adapter := NewPaystackAdapter(&PaystackConfig{
		SecretKey:  secretKey,
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	})

	request := domain.CashbackRequest{
		Source:         "balance",
		AmountKobo:     30000,
		PaymentAccount: "RCP_xxxxx",
		Reference:      "CASHBACK-001",
		Reason:         "Cashback reward",
		Currency:       "NGN",
	}

	_, err := adapter.SendCashback(
		context.Background(),
		request,
	)

	require.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"empty response body",
	)
}

func TestPaystackAdapter_SendCashback_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(`not-valid-json`))

		require.NoError(t, err)
	}))

	defer server.Close()

	secretKey := "sk_test_xxxxx"

	adapter := NewPaystackAdapter(&PaystackConfig{
		SecretKey:  secretKey,
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	})

	request := domain.CashbackRequest{
		Source:         "balance",
		AmountKobo:     30000,
		PaymentAccount: "RCP_xxxxx",
		Reference:      "CASHBACK-001",
		Reason:         "Cashback reward",
		Currency:       "NGN",
	}

	_, err := adapter.SendCashback(
		context.Background(),
		request,
	)

	require.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"decode paystack response",
	)
}

func TestPaystackAdapter_SendCashback_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(`Paystack internal error`))

		require.NoError(t, err)
	}))

	defer server.Close()

	secretKey := "sk_test_xxxxx"

	adapter := NewPaystackAdapter(&PaystackConfig{
		SecretKey:  secretKey,
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	})

	request := domain.CashbackRequest{
		Source:         "balance",
		AmountKobo:     30000,
		PaymentAccount: "RCP_xxxxx",
		Reference:      "CASHBACK-001",
		Reason:         "Cashback reward",
		Currency:       "NGN",
	}

	_, err := adapter.SendCashback(
		context.Background(),
		request,
	)

	require.Error(t, err)

	assert.Contains(
		t,
		err.Error(),
		"paystack returned HTTP 500",
	)
}
