package client

import (
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
}

// RestyClient is a thin wrapper around *resty.Client that implements ESClient.
type RestyClient struct {
	*resty.Client
}

// NewClient returns an ESClient backed by Resty and configured via Config.
func NewClient(cfg *Config) ESClient {
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

	// Configure basic auth if provided
	if cfg.Username != "" {
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

	return &RestyClient{r}
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
