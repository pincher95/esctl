package node

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

// UpdateSettingsResponse mirrors the cluster/index variants.
type UpdateSettingsResponse map[string]any

// UpdateNodeSettings updates settings for one or more nodes.
// If nodeID is empty the settings apply to all nodes.
func UpdateNodeSettings(ctx context.Context, nodeID string, body map[string]any, flat bool) (*UpdateSettingsResponse, error) {
	path := "_nodes"
	if nodeID != "" {
		path = fmt.Sprintf("_nodes/%s", nodeID)
	}
	path += "/settings"

	u := url.URL{Path: path}
	q := u.Query()
	if flat {
		q.Set("flat_settings", "true")
	}
	u.RawQuery = q.Encode()

	endpoint := u.String()

	var out UpdateSettingsResponse
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
		return nil, fmt.Errorf("failed to update node settings: %s", resp.Status())
	}
	return &out, nil
}
