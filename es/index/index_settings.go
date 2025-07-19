package index

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

func (i *index) UpdateIndexSettings(ctx context.Context, indexName string, body map[string]any, flat bool) (*IndexSettingsResponse, error) {
	if indexName == "" {
		return nil, fmt.Errorf("index name is required")
	}

	u := url.URL{Path: fmt.Sprintf("%s/_settings", indexName)}
	q := u.Query()
	q.Set("format", "json")
	if !flat {
		q.Set("flat_settings", "")
	}
	u.RawQuery = q.Encode()

	endpoint := u.String()

	var out IndexSettingsResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&out).
		Put(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to update index settings: %s", resp.Status())
	}
	return &out, nil
}
