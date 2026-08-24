package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBConfig(t *testing.T) {
	cfg := DBConfig()

	assert.Equal(t, int32(20), cfg.MaxConns)
	assert.Equal(t, int32(5), cfg.MinConns)
	assert.Equal(t, 15*time.Second, cfg.ConnTimeout)
	assert.Equal(t, 30*time.Minute, cfg.MaxConnLifetime)
	assert.Equal(t, 5*time.Minute, cfg.MaxConnIdleTime)
	assert.Equal(t, time.Minute, cfg.HealthCheckPeriod)
}

func TestLoadEnvironmentVariables(t *testing.T) {
	t.Setenv("HTTP_PORT", "8080")
	t.Setenv("DATABASE_URL", "postgres://localhost/bumpa")
	t.Setenv("PAYSTACK_BASE_URL", "https://api.paystack.co")
	t.Setenv("PAYSTACK_PUBLIC_KEY", "pk_test_example")
	t.Setenv("PAYSTACK_SECRET_KEY", "sk_test_example")
	t.Setenv("REDIS_URL", "redis://localhost:6379")

	cfg, err := LoadEnvironmentVariables()

	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "8080", cfg.HTTP_PORT)
	assert.Equal(t, "postgres://localhost/bumpa", cfg.DATABASE_URL)
	assert.Equal(t, "https://api.paystack.co", cfg.PAYSTACK_BASE_URL)
	assert.Equal(t, "pk_test_example", cfg.PAYSTACK_PUBLIC_KEY)
	assert.Equal(t, "sk_test_example", cfg.PAYSTACK_SECRET_KEY)
	assert.Equal(t, "redis://localhost:6379", cfg.REDIS_URL)
}

func TestLoadEnvironmentVariables_MissingHTTPPort(t *testing.T) {
	t.Setenv("HTTP_PORT", "")

	cfg, err := LoadEnvironmentVariables()

	assert.Nil(t, cfg)
	require.Error(t, err)

	assert.EqualError(
		t,
		err,
		"missing required environment variables: PORT",
	)
}
