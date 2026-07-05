package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

// ESClient is the minimal contract the rest of the codebase relies on.
// Having this interface makes it easy to replace the implementation in unit tests.
type ESClient interface {
	// R returns a Resty request that callers can further configure (headers, body, etc.).
	R() *resty.Request
}

// Config holds HTTP-client configuration.
type Config struct {
	RetryWaitTime time.Duration
	Timeout       time.Duration
	BaseURL       string
	RetryCount    int
	Debug         bool
	Username      string
	Password      string
	// APIKey, when set, authenticates via "Authorization: ApiKey <key>" and takes
	// precedence over basic auth.
	APIKey string
	// CACertPath is a PEM bundle used to verify the server's TLS certificate.
	CACertPath string
	// TLSInsecure disables TLS certificate verification (use only for testing).
	TLSInsecure bool
}

// RestyClient is a thin wrapper around *resty.Client that implements ESClient.
type RestyClient struct {
	*resty.Client
}

// NewClient returns an ESClient backed by Resty and configured via Config.
// It returns an error only when TLS material (a CA certificate) cannot be loaded.
func NewClient(cfg *Config) (ESClient, error) {
	r := resty.New()

	if cfg.BaseURL != "" {
		r.SetBaseURL(cfg.BaseURL)
	}

	// Set default timeout if not specified
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	r.SetTimeout(timeout)
	r.SetDebug(cfg.Debug)

	// TLS: custom CA bundle and/or skip verification.
	if cfg.TLSInsecure {
		r.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	if cfg.CACertPath != "" {
		pem, err := os.ReadFile(cfg.CACertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate %q: %w", cfg.CACertPath, err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(pem) {
			return nil, fmt.Errorf("no valid certificates found in CA file %q", cfg.CACertPath)
		}
		r.SetTLSClientConfig(&tls.Config{RootCAs: pool, InsecureSkipVerify: cfg.TLSInsecure})
	}

	// Auth: API key takes precedence over basic auth.
	if cfg.APIKey != "" {
		r.SetAuthScheme("ApiKey")
		r.SetAuthToken(cfg.APIKey)
	} else if cfg.Username != "" {
		r.SetBasicAuth(cfg.Username, cfg.Password)
	}

	// Configure retries with sensible defaults
	retryCount := cfg.RetryCount
	if retryCount == 0 {
		retryCount = 3
	}
	retryWaitTime := cfg.RetryWaitTime
	if retryWaitTime == 0 {
		retryWaitTime = 1 * time.Second
	}

	r.SetRetryCount(retryCount).
		SetRetryWaitTime(retryWaitTime).
		SetRetryMaxWaitTime(10 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			// Retry on network errors or 5xx status codes
			if err != nil {
				return true
			}
			return r.StatusCode() >= 500
		})

	return &RestyClient{r}, nil
}

// Compile-time assertion that RestyClient satisfies ESClient.
var _ ESClient = (*RestyClient)(nil)

// WithAuthToken helper for fluent auth setting.
func (c *RestyClient) WithAuthToken(token string) *RestyClient {
	c.SetAuthToken(token)
	return c
}

// WithHeader helper for fluent header setting.
func (c *RestyClient) WithHeader(key, value string) *RestyClient {
	c.SetHeader(key, value)
	return c
}
