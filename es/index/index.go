package index

import (
	"fmt"

	"github.com/pincher95/esctl/shared"
)

type IndexSettingsResponse map[string]any

type Index interface {
	UpdateIndexSettings(endpoint, index *string, body *map[string]any, flat bool) (*IndexSettingsResponse, error)
}

type index struct {
	Index
}

func NewIndex() Index {
	return &index{}
}

func (i *index) UpdateIndexSettings(endpoint, index *string, body *map[string]any, flat bool) (*IndexSettingsResponse, error) {
	if endpoint == nil {
		endpoint = new(string)
		*endpoint = fmt.Sprintf("%s/_settings?format=json", *index)
	}

	if !flat {
		*endpoint += fmt.Sprintf("&%s", "flat_settings")
	}

	var settings IndexSettingsResponse

	resp, err := shared.Client.R().SetHeader("Content-Type", "application/json").SetBody(*body).SetResult(&settings).Put(*endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to update index settings: %s", resp.Status())
	}

	return &settings, nil
}
