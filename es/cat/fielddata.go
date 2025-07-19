package cat

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

type CatFielddataResponse struct {
	ID    string `json:"id"`
	Host  string `json:"host"`
	IP    string `json:"ip"`
	Node  string `json:"node"`
	Field string `json:"field"`
	Size  string `json:"size"`
}

func (c *cat) CatFielddata(ctx context.Context, endpoint, fields, bytes string) ([]CatFielddataResponse, error) {
	if endpoint == "" {
		path := "_cat/fielddata"
		if fields != "" {
			path = fmt.Sprintf("_cat/fielddata/%s", fields)
		}

		u := url.URL{Path: path}
		q := u.Query()
		q.Set("format", "json")
		if bytes != "" {
			q.Set("bytes", bytes)
		}
		u.RawQuery = q.Encode()
		endpoint = u.String()
	}

	data := make([]CatFielddataResponse, 0)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&data).
		Get(endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get field data: %s", resp.Status())
	}

	return data, nil
}
