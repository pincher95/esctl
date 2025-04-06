package tasks

import (
	"fmt"
	"strings"

	"github.com/pincher95/esctl/shared"
)

type Tasks interface {
	// GetTask returns the task with the given ID.
	GetTasks(endpoint, taskID *string, actions *[]string) (*TasksResponse, error)
	// ListTasks returns a list of tasks.
	// ListTasks() ([]Task, error)
	// // CancelTask cancels the task with the given ID.
	// CancelTask(id string) error
	// // CancelTasks cancels the tasks with the given IDs.
	// CancelTasks(ids []string) error
	// // CancelAllTasks cancels all tasks.
	// CancelAllTasks() error
}

type tasks struct {
	Tasks
}

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

func (t *tasks) GetTasks(endpoint, taskID *string, actions *[]string) (*TasksResponse, error) {
	if endpoint == nil {
		endpoint = new(string)
		*endpoint = "_tasks?detailed&format=json"

		if taskID != nil {
			*endpoint = fmt.Sprintf("_tasks/%s?format=json", *taskID)
		}
	}

	if len(*actions) > 0 {
		actionsParam := strings.Join(*actions, ",")
		*endpoint = fmt.Sprintf("%s&actions=%s", *endpoint, actionsParam)
	}

	var tasks TasksResponse

	resp, err := shared.Client.R().SetHeader("Content-Type", "application/json").SetResult(&tasks).Get(*endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get tasks: %s", resp.String())
	}

	return &tasks, nil
}
