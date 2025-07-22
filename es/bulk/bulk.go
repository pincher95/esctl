package bulk

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/pincher95/esctl/shared"
)

type BulkResponse struct {
	Took   int                     `json:"took"`
	Errors bool                    `json:"errors"`
	Items  []map[string]BulkResult `json:"items"`
}

type BulkResult struct {
	Index   string                 `json:"_index"`
	Type    string                 `json:"_type"`
	ID      string                 `json:"_id"`
	Version int                    `json:"_version"`
	Result  string                 `json:"result"`
	Status  int                    `json:"status"`
	Error   *BulkError             `json:"error,omitempty"`
	Shards  map[string]interface{} `json:"_shards,omitempty"`
}

type BulkError struct {
	Type   string     `json:"type"`
	Reason string     `json:"reason"`
	Cause  *BulkError `json:"caused_by,omitempty"`
}

// ExecuteBulk executes a bulk request from NDJSON data
func ExecuteBulk(ctx context.Context, data io.Reader, index string, refresh string, timeout string) (BulkResponse, error) {
	var result BulkResponse

	// Read all data from reader
	bulkData, err := io.ReadAll(data)
	if err != nil {
		return result, fmt.Errorf("failed to read bulk data: %w", err)
	}

	// Ensure data ends with newline
	bulkStr := strings.TrimSpace(string(bulkData))
	if !strings.HasSuffix(bulkStr, "\n") {
		bulkStr += "\n"
	}

	// Build URL with query parameters
	endpoint := "_bulk"
	if index != "" {
		endpoint = fmt.Sprintf("%s/_bulk", index)
	}

	req := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-ndjson").
		SetBody(bulkStr).
		SetResult(&result)

	if refresh != "" {
		req.SetQueryParam("refresh", refresh)
	}
	if timeout != "" {
		req.SetQueryParam("timeout", timeout)
	}

	resp, err := req.Post(endpoint)
	if err != nil {
		return result, err
	}

	if resp.StatusCode() != 200 {
		return result, fmt.Errorf("bulk operation failed: %s", resp.Status())
	}

	return result, nil
}

// GenerateBulkTemplate generates a sample bulk NDJSON template
func GenerateBulkTemplate() string {
	return `{"index": {"_index": "my-index", "_id": "1"}}
{"field1": "value1", "field2": "value2"}
{"index": {"_index": "my-index", "_id": "2"}}
{"field1": "value3", "field2": "value4"}
{"update": {"_index": "my-index", "_id": "1"}}
{"doc": {"field3": "updated_value"}}
{"delete": {"_index": "my-index", "_id": "2"}}
`
}
