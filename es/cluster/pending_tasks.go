package cluster

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pincher95/esctl/shared"
)

// PendingTask is a single cluster-state change waiting to be applied by the master.
type PendingTask struct {
	InsertOrder       int    `json:"insert_order"`
	Priority          string `json:"priority"`
	Source            string `json:"source"`
	Executing         bool   `json:"executing"`
	TimeInQueueMillis int    `json:"time_in_queue_millis"`
	TimeInQueue       string `json:"time_in_queue"`
}

// PendingTasksResponse wraps the _cluster/pending_tasks response.
type PendingTasksResponse struct {
	Tasks []PendingTask `json:"tasks"`
}

// ClusterPendingTasks returns cluster-level changes that have not yet been executed.
// A persistently non-empty queue points at an overloaded or stuck master.
func ClusterPendingTasks(ctx context.Context) (*PendingTasksResponse, error) {
	u := url.URL{Path: "_cluster/pending_tasks"}
	q := u.Query()
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	var out PendingTasksResponse
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&out).
		Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get pending tasks: %s", resp.Status())
	}
	return &out, nil
}
