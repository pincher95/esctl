package cat

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

type CatPluginResponse struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Component   string `json:"component,omitempty"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
}

func (c *cat) CatPlugins(ctx context.Context, endpoint string) ([]CatPluginResponse, error) {
	if endpoint == "" {
		u := url.URL{Path: "_cat/plugins"}
		q := u.Query()
		q.Set("format", "json")
		q.Set("h", "id,name,component,version,description")
		u.RawQuery = q.Encode()
		endpoint = u.String()
	}

	plugins := make([]CatPluginResponse, 0)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&plugins).
		Get(endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get nodes plugins: %s", resp.Status())
	}

	return plugins, nil
}
