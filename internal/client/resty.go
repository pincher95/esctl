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

	r.SetTimeout(cfg.Timeout)
	r.SetDebug(cfg.Debug)

	if cfg.RetryCount > 0 {
		r.
			SetRetryCount(cfg.RetryCount).
			SetRetryWaitTime(cfg.RetryWaitTime)
	}

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
