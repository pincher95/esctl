package index

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/es/common"
	"github.com/pincher95/esctl/shared"
)

// IndexOpResponse wraps the common _shards response for index maintenance APIs.
type IndexOpResponse struct {
	Shards common.ResponseShards `json:"_shards"`
}

// Refresh triggers a refresh of one or all indices.
func (i *index) Refresh(ctx context.Context, indexName string) (*IndexOpResponse, error) {
	u := url.URL{}
	if indexName != "" {
		u.Path = fmt.Sprintf("%s/_refresh", indexName)
	} else {
		u.Path = "_refresh"
	}
	q := u.Query()
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	endpoint := u.String()

	var out IndexOpResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&out).
		Post(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to refresh index: %s", resp.Status())
	}
	return &out, nil
}

// Flush triggers a flush of one or all indices.
func (i *index) Flush(ctx context.Context, indexName string) (*IndexOpResponse, error) {
	u := url.URL{}
	if indexName != "" {
		u.Path = fmt.Sprintf("%s/_flush", indexName)
	} else {
		u.Path = "_flush"
	}
	q := u.Query()
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	endpoint := u.String()

	var out IndexOpResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&out).
		Post(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to flush index: %s", resp.Status())
	}
	return &out, nil
}

// Forcemerge optimizes one or all indices. Optional flags mirror Elasticsearch parameters.
func (i *index) Forcemerge(ctx context.Context, indexName string, maxNumSegments int, onlyExpungeDeletes bool, flush bool) (*IndexOpResponse, error) {
	u := url.URL{}
	if indexName != "" {
		u.Path = fmt.Sprintf("%s/_forcemerge", indexName)
	} else {
		u.Path = "_forcemerge"
	}
	q := u.Query()
	q.Set("format", "json")
	if maxNumSegments > 0 {
		q.Set("max_num_segments", fmt.Sprintf("%d", maxNumSegments))
	}
	if onlyExpungeDeletes {
		q.Set("only_expunge_deletes", "true")
	}
	if flush {
		q.Set("flush", "true")
	}
	u.RawQuery = q.Encode()

	endpoint := u.String()

	var out IndexOpResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&out).
		Post(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to forcemerge index: %s", resp.Status())
	}
	return &out, nil
}
