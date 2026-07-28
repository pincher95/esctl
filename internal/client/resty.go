package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	eserrors "github.com/pincher95/esctl/internal/errors"
)

// DefaultTimeout is the per-request timeout used when none is configured. It is
// deliberately forgiving: several read-only ES endpoints (_cat/snapshots,
// _cluster/stats on large clusters) can take tens of seconds to answer.
const DefaultTimeout = 60 * time.Second

// silentLogger discards resty's internal logging when debug mode is off.
type silentLogger struct{}

func (silentLogger) Errorf(string, ...any) {}
func (silentLogger) Warnf(string, ...any)  {}
func (silentLogger) Debugf(string, ...any) {}

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

	// Set default timeout if not specified. Some read-only Elasticsearch endpoints
	// are legitimately slow (listing snapshots in a large repository reads metadata
	// from remote storage), so the default is generous; raise it with --timeout.
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	r.SetTimeout(timeout)
	r.SetDebug(cfg.Debug)

	// resty logs retries and transport errors on its own logger. esctl reports
	// errors itself, so keep resty's output for --debug only.
	if !cfg.Debug {
		r.SetLogger(silentLogger{})
	}

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
		AddRetryCondition(func(resp *resty.Response, err error) bool {
			// Only retry requests that are safe to repeat. Re-sending a timed-out
			// mutation can duplicate work the server already accepted — e.g. a
			// snapshot restore that had already started, whose retry then fails
			// with "an open index with same name already exists".
			if resp == nil || resp.Request == nil || !isRetryableMethod(resp.Request.Method) {
				return false
			}
			if err != nil {
				// A timeout means the server is slow, not that the connection
				// blipped: re-sending re-runs the same expensive query and adds
				// load to an already-struggling cluster, so surface it instead
				// and let the caller raise --timeout.
				return !eserrors.IsTimeout(err)
			}
			return resp.StatusCode() >= 500
		})

	return &RestyClient{r}, nil
}

// isRetryableMethod reports whether a request with this method can be safely
// re-sent. Only read-only methods qualify: Elasticsearch write endpoints are not
// reliably idempotent (creating a snapshot or restoring one twice errors out), so
// POST/PUT/PATCH/DELETE are never retried automatically.
func isRetryableMethod(method string) bool {
	switch strings.ToUpper(method) {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
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
