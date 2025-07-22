package reindex

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

type ReindexRequest struct {
	Source    ReindexSource  `json:"source"`
	Dest      ReindexDest    `json:"dest"`
	Script    *ReindexScript `json:"script,omitempty"`
	Size      *int           `json:"size,omitempty"`
	Conflicts string         `json:"conflicts,omitempty"`
}

type ReindexSource struct {
	Index  interface{}              `json:"index"`
	Type   string                   `json:"type,omitempty"`
	Query  map[string]interface{}   `json:"query,omitempty"`
	Sort   []map[string]interface{} `json:"sort,omitempty"`
	Size   *int                     `json:"size,omitempty"`
	Remote *ReindexRemote           `json:"remote,omitempty"`
}

type ReindexDest struct {
	Index       string `json:"index"`
	Type        string `json:"type,omitempty"`
	VersionType string `json:"version_type,omitempty"`
	OpType      string `json:"op_type,omitempty"`
	Pipeline    string `json:"pipeline,omitempty"`
}

type ReindexScript struct {
	Source string                 `json:"source"`
	Lang   string                 `json:"lang,omitempty"`
	Params map[string]interface{} `json:"params,omitempty"`
}

type ReindexRemote struct {
	Host           string            `json:"host"`
	Username       string            `json:"username,omitempty"`
	Password       string            `json:"password,omitempty"`
	SocketTimeout  string            `json:"socket_timeout,omitempty"`
	ConnectTimeout string            `json:"connect_timeout,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
}

type ReindexResponse struct {
	Task                 string         `json:"task,omitempty"`
	Took                 int            `json:"took,omitempty"`
	TimedOut             bool           `json:"timed_out,omitempty"`
	Total                int            `json:"total,omitempty"`
	Updated              int            `json:"updated,omitempty"`
	Created              int            `json:"created,omitempty"`
	Deleted              int            `json:"deleted,omitempty"`
	Batches              int            `json:"batches,omitempty"`
	VersionConflicts     int            `json:"version_conflicts,omitempty"`
	Noops                int            `json:"noops,omitempty"`
	Retries              ReindexRetries `json:"retries,omitempty"`
	ThrottledMillis      int            `json:"throttled_millis,omitempty"`
	RequestsPerSecond    float64        `json:"requests_per_second,omitempty"`
	ThrottledUntilMillis int            `json:"throttled_until_millis,omitempty"`
	Failures             []interface{}  `json:"failures,omitempty"`
}

type ReindexRetries struct {
	Bulk   int `json:"bulk"`
	Search int `json:"search"`
}

// StartReindex starts a reindex operation
func StartReindex(ctx context.Context, request ReindexRequest, waitForCompletion bool, requestsPerSecond float64, timeout string, refresh bool, slices interface{}) (ReindexResponse, error) {
	var result ReindexResponse

	u := url.URL{Path: "_reindex"}
	q := u.Query()
	if waitForCompletion {
		q.Set("wait_for_completion", "true")
	}
	if requestsPerSecond > 0 {
		q.Set("requests_per_second", fmt.Sprintf("%.2f", requestsPerSecond))
	}
	if timeout != "" {
		q.Set("timeout", timeout)
	}
	if refresh {
		q.Set("refresh", "true")
	}
	if slices != nil {
		q.Set("slices", fmt.Sprintf("%v", slices))
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
		return result, fmt.Errorf("failed to start reindex: %s", resp.Status())
	}

	return result, nil
}

// GetReindexTaskStatus gets the status of a reindex task
func GetReindexTaskStatus(ctx context.Context, taskID string) (map[string]interface{}, error) {
	var result map[string]interface{}

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&result).
		Get(fmt.Sprintf("_tasks/%s", taskID))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get task status: %s", resp.Status())
	}

	return result, nil
}

// CancelReindexTask cancels a reindex task
func CancelReindexTask(ctx context.Context, taskID string) error {
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Post(fmt.Sprintf("_tasks/%s/_cancel", taskID))
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to cancel task: %s", resp.Status())
	}

	return nil
}

// UpdateReindexByQuery updates the requests_per_second for a running reindex task
func UpdateReindexByQuery(ctx context.Context, taskID string, requestsPerSecond float64) error {
	u := url.URL{Path: fmt.Sprintf("_reindex/%s/_rethrottle", taskID)}
	q := u.Query()
	q.Set("requests_per_second", fmt.Sprintf("%.2f", requestsPerSecond))
	u.RawQuery = q.Encode()

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Post(u.String())
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to update reindex throttle: %s", resp.Status())
	}

	return nil
}
