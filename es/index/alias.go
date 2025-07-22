package index

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/shared"
)

type AliasAction struct {
	Add    *AliasAdd    `json:"add,omitempty"`
	Remove *AliasRemove `json:"remove,omitempty"`
}

type AliasAdd struct {
	Index   string                 `json:"index"`
	Alias   string                 `json:"alias"`
	Filter  map[string]interface{} `json:"filter,omitempty"`
	Routing string                 `json:"routing,omitempty"`
}

type AliasRemove struct {
	Index string `json:"index"`
	Alias string `json:"alias"`
}

type AliasRequest struct {
	Actions []AliasAction `json:"actions"`
}

type AliasResponse struct {
	Acknowledged bool `json:"acknowledged"`
}

type AliasInfo struct {
	Aliases map[string]AliasDetails `json:"aliases"`
}

type AliasDetails struct {
	Filter         map[string]interface{} `json:"filter,omitempty"`
	Routing        string                 `json:"routing,omitempty"`
	IndexRouting   string                 `json:"index_routing,omitempty"`
	SearchRouting  string                 `json:"search_routing,omitempty"`
	IsWriteIndex   bool                   `json:"is_write_index,omitempty"`
	IsHidden       bool                   `json:"is_hidden,omitempty"`
}

type AliasListResponse map[string]AliasInfo

// AddAlias adds an alias to one or more indices
func AddAlias(ctx context.Context, indices []string, alias string, filter map[string]interface{}, routing string) error {
	if len(indices) == 0 {
		return fmt.Errorf("no indices specified")
	}

	var actions []AliasAction
	for _, index := range indices {
		action := AliasAction{
			Add: &AliasAdd{
				Index: index,
				Alias: alias,
			},
		}

		if filter != nil {
			action.Add.Filter = filter
		}

		if routing != "" {
			action.Add.Routing = routing
		}

		actions = append(actions, action)
	}

	request := AliasRequest{
		Actions: actions,
	}

	var result AliasResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&result).
		Post("_aliases")
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to add alias: %s", resp.Status())
	}

	return nil
}

// RemoveAlias removes an alias from indices
func RemoveAlias(ctx context.Context, indices []string, alias string) error {
	if len(indices) == 0 {
		return fmt.Errorf("no indices specified")
	}

	var actions []AliasAction
	for _, index := range indices {
		actions = append(actions, AliasAction{
			Remove: &AliasRemove{
				Index: index,
				Alias: alias,
			},
		})
	}

	request := AliasRequest{
		Actions: actions,
	}

	var result AliasResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&result).
		Post("_aliases")
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to remove alias: %s", resp.Status())
	}

	return nil
}

// MoveAlias atomically moves an alias from one index to another
func MoveAlias(ctx context.Context, fromIndex, toIndex, alias string) error {
	actions := []AliasAction{
		{
			Remove: &AliasRemove{
				Index: fromIndex,
				Alias: alias,
			},
		},
		{
			Add: &AliasAdd{
				Index: toIndex,
				Alias: alias,
			},
		},
	}

	request := AliasRequest{
		Actions: actions,
	}

	var result AliasResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&result).
		Post("_aliases")
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to move alias: %s", resp.Status())
	}

	return nil
}

// ListAliases lists aliases for indices matching a pattern
func ListAliases(ctx context.Context, indexPattern, aliasPattern string) (AliasListResponse, error) {
	var result AliasListResponse

	var path strings.Builder
	if indexPattern != "" {
		path.WriteString(indexPattern)
	} else {
		path.WriteString("*")
	}
	path.WriteString("/_alias")
	if aliasPattern != "" {
		path.WriteString("/")
		path.WriteString(aliasPattern)
	}

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(path.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to list aliases: %s", resp.Status())
	}

	return result, nil
}

// GetAlias gets aliases for specific indices and/or alias names
func GetAlias(ctx context.Context, indices []string, aliases []string) (AliasListResponse, error) {
	var result AliasListResponse

	var path strings.Builder
	if len(indices) > 0 {
		path.WriteString(strings.Join(indices, ","))
	} else {
		path.WriteString("*")
	}
	path.WriteString("/_alias")
	if len(aliases) > 0 {
		path.WriteString("/")
		path.WriteString(strings.Join(aliases, ","))
	}

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(path.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get aliases: %s", resp.Status())
	}

	return result, nil
}
