package cat

import (
	"fmt"

	"github.com/pincher95/esctl/shared"
)

type Plugin struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Component   string `json:"component,omitempty"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
}

func CatPlugins(endpoint *string) ([]Plugin, error) {
	if endpoint == nil {
		endpoint = new(string)
		*endpoint = "_cat/plugins?format=json&h=id,name,component,version,description"
	}

	plugins := make([]Plugin, 0)

	resp, err := shared.Client.R().SetHeader("Content-Type", "application/json").SetResult(&plugins).Get(*endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get nodes plugins: %s", resp.Status())
	}

	return plugins, nil
}
