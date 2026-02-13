package index

import (
	"context"
	"fmt"
	"maps"
	"net/url"
	"strings"

	"github.com/pincher95/esctl/shared"
)

// IndexOperationResponse wraps the common response for index lifecycle operations.
type IndexOperationResponse struct {
	Acknowledged       bool   `json:"acknowledged"`
	ShardsAcknowledged bool   `json:"shards_acknowledged"`
	Index              string `json:"index,omitempty"`
}

// Open opens one or more closed indices.
func (i *index) Open(ctx context.Context, indices []string) (*IndexOperationResponse, error) {
	if len(indices) == 0 {
		return nil, fmt.Errorf("at least one index must be specified")
	}

	indexPath := strings.Join(indices, ",")
	u := url.URL{
		Path: fmt.Sprintf("%s/_open", indexPath),
	}

	var out IndexOperationResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&out).
		Post(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to open index: %s", resp.Status())
	}
	return &out, nil
}

// Close closes one or more open indices.
func (i *index) Close(ctx context.Context, indices []string) (*IndexOperationResponse, error) {
	if len(indices) == 0 {
		return nil, fmt.Errorf("at least one index must be specified")
	}

	indexPath := strings.Join(indices, ",")
	u := url.URL{
		Path: fmt.Sprintf("%s/_close", indexPath),
	}

	var out IndexOperationResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&out).
		Post(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to close index: %s", resp.Status())
	}
	return &out, nil
}

// Clone clones an index to a new index with the same data and settings.
func (i *index) Clone(ctx context.Context, source, target string, settings map[string]any) (*IndexOperationResponse, error) {
	if source == "" {
		return nil, fmt.Errorf("source index must be specified")
	}
	if target == "" {
		return nil, fmt.Errorf("target index must be specified")
	}

	u := url.URL{
		Path: fmt.Sprintf("%s/_clone/%s", source, target),
	}

	body := make(map[string]any)
	if settings != nil && len(settings) > 0 {
		body["settings"] = settings
	}

	var out IndexOperationResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&out).
		Post(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to clone index: %s", resp.Status())
	}
	return &out, nil
}

// Split splits an index into a new index with more primary shards.
func (i *index) Split(ctx context.Context, source, target string, shards int, settings map[string]any) (*IndexOperationResponse, error) {
	if source == "" {
		return nil, fmt.Errorf("source index must be specified")
	}
	if target == "" {
		return nil, fmt.Errorf("target index must be specified")
	}
	if shards <= 0 {
		return nil, fmt.Errorf("number of shards must be positive")
	}

	u := url.URL{
		Path: fmt.Sprintf("%s/_split/%s", source, target),
	}

	body := map[string]any{
		"settings": map[string]any{
			"index.number_of_shards": shards,
		},
	}

	// Merge additional settings if provided
	if settings != nil && len(settings) > 0 {
		if settingsMap, ok := body["settings"].(map[string]any); ok {
			maps.Copy(settingsMap, settings)
		}
	}

	var out IndexOperationResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&out).
		Post(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to split index: %s", resp.Status())
	}
	return &out, nil
}

// Shrink shrinks an index into a new index with fewer primary shards.
func (i *index) Shrink(ctx context.Context, source, target string, shards int, settings map[string]any) (*IndexOperationResponse, error) {
	if source == "" {
		return nil, fmt.Errorf("source index must be specified")
	}
	if target == "" {
		return nil, fmt.Errorf("target index must be specified")
	}
	if shards <= 0 {
		return nil, fmt.Errorf("number of shards must be positive")
	}

	u := url.URL{
		Path: fmt.Sprintf("%s/_shrink/%s", source, target),
	}

	body := map[string]any{
		"settings": map[string]any{
			"index.number_of_shards": shards,
		},
	}

	// Merge additional settings if provided
	if settings != nil && len(settings) > 0 {
		if settingsMap, ok := body["settings"].(map[string]any); ok {
			maps.Copy(settingsMap, settings)
		}
	}

	var out IndexOperationResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&out).
		Post(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to shrink index: %s", resp.Status())
	}
	return &out, nil
}
