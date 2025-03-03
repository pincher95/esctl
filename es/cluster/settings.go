package cluster

import (
	"fmt"

	"github.com/pincher95/esctl/shared"
)

type Settings map[string]any

func ClusterSettings(endpoint *string, flat, defaults bool) (*Settings, error) {
	if endpoint == nil {
		endpoint = new(string)
		*endpoint = "_cluster/settings?format=json"
	}

	if !flat {
		*endpoint += fmt.Sprintf("&%s", "flat_settings")
	}

	if defaults {
		*endpoint += fmt.Sprintf("&%s", "include_defaults")
	}

	var settings Settings

	resp, err := shared.Client.R().SetHeader("Content-Type", "application/json").SetResult(&settings).Get(*endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get cluster settings: %s", resp.Status())
	}

	return &settings, nil
}
