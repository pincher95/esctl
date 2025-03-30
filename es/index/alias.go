package index

import (
	"fmt"

	"github.com/pincher95/esctl/shared"
)

type IndexAliasResponse map[string]AliasDetail

type AliasDetail struct {
	Aliases map[string]AliasInfo `json:"aliases"`
}

type AliasInfo struct {
	Filter        map[string]any `json:"filter"`
	IndexRouting  string         `json:"index_routing"`
	SearchRouting string         `json:"search_routing"`
	Routing       string         `json:"routing"`
	IsWriteIndex  bool           `json:"is_write_index"`
	IsHidden      bool           `json:"is_hidden"`
}

func (i *index) GetAliases(endpoint, index *string) (*IndexAliasResponse, error) {
	if endpoint == nil {
		endpoint = new(string)
		*endpoint = "_alias/_all?format=json"

		if index != nil {
			*endpoint = fmt.Sprintf("_alias/%s?format=json", *index)
		}
	}

	// alias := make([]IndexAliasResponse, 0)

	var aliases IndexAliasResponse

	resp, err := shared.Client.R().SetHeader("Content-Type", "application/json").SetResult(&aliases).Get(*endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get alias: %s", resp.Status())
	}

	return &aliases, nil
}
