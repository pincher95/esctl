package shared

import (
	"time"

	"github.com/pincher95/esctl/internal/client"
)

var (
	Client  client.ESClient
	Context string

	ElasticsearchProtocol string
	ElasticsearchUsername string
	ElasticsearchPassword string
	ElasticsearchHost     string
	ElasticsearchPort     int
	Debug                 bool
	TimeoutDuration       time.Duration
	OutputFormat          string
)

// SetClient allows tests or callers to swap the underlying HTTP client implementation.
func SetClient(c client.ESClient) {
	Client = c
}
