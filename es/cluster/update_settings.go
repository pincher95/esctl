package cluster

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/shared"
)

type UpdateSettingsResponse map[string]any

func UpdateClusterSettings(ctx context.Context, body map[string]any) (*UpdateSettingsResponse, error) {
	endpoint := "_cluster/settings?pretty&flat_settings=true"
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
		return nil, fmt.Errorf("failed to update settings: %s", resp.Status())
	}
	return &out, nil
}
