package index

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

type IndexCacheClearResponse struct {
	Shards struct {
		Total      int             `json:"total"`
		Successful int             `json:"successful"`
		Failed     int             `json:"failed"`
		Failures   []FailuresShard `json:"failures"`
	} `json:"_shards"`
}

func (i *index) CacheClear(ctx context.Context, indexName string) (*IndexCacheClearResponse, error) {
	u := url.URL{}
	if indexName != "" {
		u.Path = fmt.Sprintf("%s/_cache/clear", indexName)
	} else {
		u.Path = "_cache/clear"
	}
	q := u.Query()
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	endpoint := u.String()

	var out IndexCacheClearResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&out).
		Post(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to clear cache: %s", resp.Status())
	}
	return &out, nil
}
