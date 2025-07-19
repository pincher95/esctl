package cat

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/shared"
)

type CatHealthResponse struct {
	Epoch               string `json:"epoch"`
	Timestamp           string `json:"timestamp"`
	Cluster             string `json:"cluster"`
	Status              string `json:"status"`
	NodeTotal           int    `json:"node.total,string"`
	NodeData            int    `json:"node.data,string"`
	Shards              int    `json:"shards,string"`
	Pri                 int    `json:"pri,string"`
	Relo                int    `json:"relo,string"`
	Init                int    `json:"init,string"`
	Unassign            int    `json:"unassign,string"`
	PendingTasks        int    `json:"pending_tasks,string"`
	ActiveShardsPercent string `json:"active_shards_percent"`
}

func (c *cat) CatHealth(ctx context.Context) (*CatHealthResponse, error) {
	endpoint := "_cat/health?format=json"
	rows := make([]CatHealthResponse, 0)
	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&rows).
		Get(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get health: %s", resp.Status())
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("empty health response")
	}
	return &rows[0], nil
}
