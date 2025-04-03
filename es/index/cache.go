package index

import (
	"fmt"

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

func (i *index) CacheClear(endpoint, index *string) (*IndexCacheClearResponse, error) {
	if endpoint == nil {
		endpoint = new(string)
		*endpoint = "_cache/clear?format=json"

		if index != nil {
			*endpoint = fmt.Sprintf("_alias/%s?format=json", *index)
		}
	}

	var cacheClear IndexCacheClearResponse

	resp, err := shared.Client.R().SetHeader("Content-Type", "application/json").SetResult(&cacheClear).Post(*endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get alias: %s", resp.Status())
	}

	return &cacheClear, nil
}
