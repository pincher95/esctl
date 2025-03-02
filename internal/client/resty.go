package client

import (
	"time"

	"github.com/go-resty/resty/v2"
)

// Config holds any configuration you want to expose for the HTTP client.
// Add more fields as needed (e.g. base URL, timeouts, etc.).
type Config struct {
	RetryWaitTime time.Duration
	Timeout       time.Duration
	BaseURL       string
	RetryCount    int
	Debug         bool
}

// Client is a wrapper around *resty.Client that you can customize.
type Client struct {
	*resty.Client
}

// NewClient returns a configured resty client based on the given Config.
func NewClient(cfg *Config) *Client {
	r := resty.New()

	// Set base URL if provided
	if cfg.BaseURL != "" {
		r.SetBaseURL(cfg.BaseURL)
	}

	// Set request timeout
	r.SetTimeout(cfg.Timeout)

	// Enable debug if needed
	r.SetDebug(cfg.Debug)

	// Set retry count and wait time
	if cfg.RetryCount > 0 {
		r.
			SetRetryCount(cfg.RetryCount).
			SetRetryWaitTime(cfg.RetryWaitTime)
	}

	return &Client{r}
}

func (c *Client) WithAuthToken(token string) *Client {
	c.SetAuthToken(token)
	return c
}

func (c *Client) WithHeader(key, value string) *Client {
	c.SetHeader(key, value)
	return c
}
