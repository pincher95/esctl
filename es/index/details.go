package index

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/shared"
)

type IndexMappings struct {
	Mappings any `json:"mappings"`
}

type IndexSettings struct {
	Settings any `json:"settings"`
}

type MappingsResponse map[string]IndexMappings

type SettingsResponse map[string]IndexSettings

type IndexDetails struct {
	Settings any `json:"settings,omitempty"`
	Mappings any `json:"mappings,omitempty"`
}

type IndexDetailsResponse map[string]IndexDetails

// GetIndexDetails fetches mappings and/or settings for the given index.
func GetIndexDetails(ctx context.Context, indexName string, wantMappings, wantSettings, flat, includeDefaults bool) (IndexDetailsResponse, error) {
	var mappingsResp MappingsResponse
	var settingsResp SettingsResponse

	if wantMappings {
		if err := requestJSON(ctx, fmt.Sprintf("%s/_mappings", indexName), &mappingsResp); err != nil {
			return nil, fmt.Errorf("failed to get mappings: %w", err)
		}
	}

	if wantSettings {
		endpoint := fmt.Sprintf("%s/_settings", indexName)
		sep := "?"
		if flat {
			endpoint += sep + "flat_settings=true"
			sep = "&"
		}
		if includeDefaults {
			endpoint += sep + "include_defaults=true"
		}
		if err := requestJSON(ctx, endpoint, &settingsResp); err != nil {
			return nil, fmt.Errorf("failed to get settings: %w", err)
		}
	}

	merged := make(IndexDetailsResponse)
	for idx, m := range mappingsResp {
		d := IndexDetails{Mappings: m.Mappings}
		if s, ok := settingsResp[idx]; ok {
			d.Settings = s.Settings
		}
		merged[idx] = d
	}
	for idx, s := range settingsResp {
		if _, ok := merged[idx]; !ok {
			merged[idx] = IndexDetails{Settings: s.Settings}
		}
	}
	return merged, nil
}

// requestJSON is a lightweight wrapper around shared.Client for GET requests.
func requestJSON(ctx context.Context, endpoint string, target any) error {
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(target).
		Get(endpoint)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("unexpected status %d", resp.StatusCode())
	}
	return nil
}
