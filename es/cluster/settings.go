package cluster

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

type Settings map[string]any

func ClusterSettings(ctx context.Context, flatSettings, includeDefaults bool) (*Settings, error) {
	u := url.URL{Path: "_cluster/settings"}
	q := u.Query()
	q.Set("format", "json")
	if !flatSettings {
		q.Set("flat_settings", "")
	}
	if includeDefaults {
		q.Set("include_defaults", "")
	}
	u.RawQuery = q.Encode()

	endpoint := u.String()

	var out Settings
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&out).
		Get(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get cluster settings: %s", resp.Status())
	}
	return &out, nil
}
