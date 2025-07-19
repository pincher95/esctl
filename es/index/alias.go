package index

import (
	"context"
	"fmt"
	"net/url"

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

func (i *index) GetAliases(ctx context.Context, indexName string) (*IndexAliasResponse, error) {
	u := url.URL{}
	if indexName != "" {
		u.Path = fmt.Sprintf("%s/_alias", indexName)
	} else {
		u.Path = "_alias/_all"
	}
	q := u.Query()
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	endpoint := u.String()

	var aliases IndexAliasResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&aliases).
		Get(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get alias: %s", resp.Status())
	}
	return &aliases, nil
}
