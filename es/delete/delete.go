package delete

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/pincher95/esctl/shared"
)

type DeleteByQueryRequest struct {
	Query               map[string]any   `json:"query"`
	MaxDocs             *int             `json:"max_docs,omitempty"`
	Conflicts           string           `json:"conflicts,omitempty"`
	Refresh             bool             `json:"refresh,omitempty"`
	Timeout             string           `json:"timeout,omitempty"`
	WaitForActiveShards string           `json:"wait_for_active_shards,omitempty"`
	WaitForCompletion   bool             `json:"wait_for_completion,omitempty"`
	RequestsPerSecond   *float64         `json:"requests_per_second,omitempty"`
	Scroll              string           `json:"scroll,omitempty"`
	ScrollSize          *int             `json:"scroll_size,omitempty"`
	Sort                []map[string]any `json:"sort,omitempty"`
	SearchType          string           `json:"search_type,omitempty"`
	SearchTimeout       string           `json:"search_timeout,omitempty"`
	Slices              any              `json:"slices,omitempty"`
}

type DeleteByQueryResponse struct {
	Took                 int           `json:"took"`
	TimedOut             bool          `json:"timed_out"`
	Total                int           `json:"total"`
	Deleted              int           `json:"deleted"`
	Batches              int           `json:"batches"`
	VersionConflicts     int           `json:"version_conflicts"`
	Noops                int           `json:"noops"`
	Retries              DeleteRetries `json:"retries"`
	ThrottledMillis      int           `json:"throttled_millis"`
	RequestsPerSecond    float64       `json:"requests_per_second"`
	ThrottledUntilMillis int           `json:"throttled_until_millis"`
	Failures             []any         `json:"failures"`
	Task                 string        `json:"task,omitempty"`
}

type DeleteRetries struct {
	Bulk   int `json:"bulk"`
	Search int `json:"search"`
}

type IndexDeleteResponse struct {
	Acknowledged bool `json:"acknowledged"`
}

type AliasAction struct {
	Remove *AliasRemove `json:"remove,omitempty"`
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

// DeleteIndex deletes one or more indices
func DeleteIndex(ctx context.Context, indices []string, ignoreUnavailable bool, allowNoIndices bool, expandWildcards string) error {
	if len(indices) == 0 {
		return fmt.Errorf("no indices specified")
	}

	indexPath := strings.Join(indices, ",")

	u := url.URL{Path: indexPath}
	q := u.Query()
	if ignoreUnavailable {
		q.Set("ignore_unavailable", "true")
	}
	if allowNoIndices {
		q.Set("allow_no_indices", "true")
	}
	if expandWildcards != "" {
		q.Set("expand_wildcards", expandWildcards)
	}
	u.RawQuery = q.Encode()

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Delete(u.String())
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to delete indices: %s", resp.Status())
	}

	return nil
}

// DeleteByQuery deletes documents matching a query
func DeleteByQuery(ctx context.Context, indices []string, request DeleteByQueryRequest) (DeleteByQueryResponse, error) {
	var result DeleteByQueryResponse

	var indexPath string
	if len(indices) > 0 {
		indexPath = strings.Join(indices, ",") + "/"
	}

	u := url.URL{Path: indexPath + "_delete_by_query"}
	q := u.Query()
	if request.WaitForCompletion {
		q.Set("wait_for_completion", "true")
	}
	if request.Refresh {
		q.Set("refresh", "true")
	}
	if request.Timeout != "" {
		q.Set("timeout", request.Timeout)
	}
	if request.WaitForActiveShards != "" {
		q.Set("wait_for_active_shards", request.WaitForActiveShards)
	}
	if request.RequestsPerSecond != nil {
		q.Set("requests_per_second", fmt.Sprintf("%.2f", *request.RequestsPerSecond))
	}
	if request.Scroll != "" {
		q.Set("scroll", request.Scroll)
	}
	if request.ScrollSize != nil {
		q.Set("scroll_size", fmt.Sprintf("%d", *request.ScrollSize))
	}
	if request.SearchType != "" {
		q.Set("search_type", request.SearchType)
	}
	if request.SearchTimeout != "" {
		q.Set("search_timeout", request.SearchTimeout)
	}
	if request.Slices != nil {
		q.Set("slices", fmt.Sprintf("%v", request.Slices))
	}
	u.RawQuery = q.Encode()

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&result).
		Post(u.String())
	if err != nil {
		return result, err
	}

	if resp.StatusCode() != 200 {
		return result, fmt.Errorf("failed to delete by query: %s", resp.Status())
	}

	return result, nil
}

// DeleteAlias removes an alias from indices
func DeleteAlias(ctx context.Context, indices []string, aliases []string) error {
	if len(indices) == 0 || len(aliases) == 0 {
		return fmt.Errorf("both indices and aliases must be specified")
	}

	var actions []AliasAction
	for _, index := range indices {
		for _, alias := range aliases {
			actions = append(actions, AliasAction{
				Remove: &AliasRemove{
					Index: index,
					Alias: alias,
				},
			})
		}
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
		return fmt.Errorf("failed to delete alias: %s", resp.Status())
	}

	return nil
}
