package tasks

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/pincher95/esctl/shared"
)

type Tasks interface {
	// GetTasks returns running tasks optionally filtered by taskID and actions.
	GetTasks(ctx context.Context, taskID string, actions []string) (*TasksResponse, error)
	// CancelTasks stops running tasks optionally filtered by taskID and actions.
	CancelTasks(ctx context.Context, taskID string, actions []string) (*TasksResponse, error)
}

// tasks is a concrete implementation of the Tasks interface.
// It carries no internal state at the moment.
type tasks struct{}

func NewTasks() Tasks {
	return &tasks{}
}

type TasksResponse struct {
	Nodes map[string]TaskNode `json:"nodes"`
}

type TaskNode struct {
	Name             string            `json:"name"`
	TransportAddress string            `json:"transport_address"`
	Host             string            `json:"host"`
	IP               string            `json:"ip"`
	Roles            []string          `json:"roles"`
	Attributes       map[string]string `json:"attributes"`
	Tasks            map[string]Task   `json:"tasks"`
}

type Task struct {
	Node               string         `json:"node"`
	ID                 int64          `json:"id"`
	Type               string         `json:"type"`
	Action             string         `json:"action"`
	Description        string         `json:"description"`
	StartTimeInMillis  int64          `json:"start_time_in_millis"`
	RunningTimeInNanos int64          `json:"running_time_in_nanos"`
	Cancellable        bool           `json:"cancellable"`
	Cancelled          bool           `json:"cancelled"`
	ParentTaskID       string         `json:"parent_task_id"`
	Headers            map[string]any `json:"headers"`
}

func (t *tasks) GetTasks(ctx context.Context, taskID string, actions []string) (*TasksResponse, error) {
	// Build endpoint using url.URL for safety.
	u := url.URL{Path: "_tasks"}
	if taskID != "" {
		u.Path = fmt.Sprintf("_tasks/%s", taskID)
	}

	q := u.Query()
	q.Set("format", "json")
	q.Set("detailed", "")
	if len(actions) > 0 {
		q.Set("actions", strings.Join(actions, ","))
	}
	u.RawQuery = q.Encode()

	endpoint := u.String()

	var respBody TasksResponse

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&respBody).
		Get(endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get tasks: %s", resp.String())
	}

	return &respBody, nil
}

func (t *tasks) CancelTasks(ctx context.Context, taskID string, actions []string) (*TasksResponse, error) {
	u := url.URL{Path: "_tasks/_cancel"}
	if taskID != "" {
		u.Path = fmt.Sprintf("_tasks/%s/_cancel", taskID)
	}

	q := u.Query()
	q.Set("format", "json")
	q.Set("detailed", "")
	if len(actions) > 0 {
		q.Set("actions", strings.Join(actions, ","))
	}
	u.RawQuery = q.Encode()

	endpoint := u.String()

	var respBody TasksResponse

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&respBody).
		Post(endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to cancel tasks: %s", resp.String())
	}

	return &respBody, nil
}
