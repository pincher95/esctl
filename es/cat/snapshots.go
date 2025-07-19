package cat

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/pincher95/esctl/shared"
)

type CatSnapshotResponse struct {
	ID               string `json:"id"`
	Status           string `json:"status"`
	StartEpoch       int    `json:"start_epoch,string"`
	StartTime        string `json:"start_time"`
	EndEpoch         int    `json:"end_epoch,string"`
	EndTime          string `json:"end_time"`
	Duration         string `json:"duration"`
	Indices          int    `json:"indices,string"`
	SuccessfulShards int    `json:"successful_shards,string"`
	FailedShards     int    `json:"failed_shards,string"`
	TotalShards      int    `json:"total_shards,string"`
	Reason           string `json:"reason"`
}

func (c *cat) CatSnapshots(ctx context.Context, endpoint, repository, snapshotName string) ([]CatSnapshotResponse, error) {
	if endpoint == "" {
		path := "_cat/snapshots"
		if repository != "" {
			path = fmt.Sprintf("_cat/snapshots/%s", repository)
		}

		u := url.URL{Path: path}
		q := u.Query()
		q.Set("format", "json")
		u.RawQuery = q.Encode()
		endpoint = u.String()
	}

	snapshots := make([]CatSnapshotResponse, 0)

	resp, err := shared.Client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetResult(&snapshots).
		Get(endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get snapshots: %s", resp.Status())
	}

	if snapshotName != "" {
		filtered := make([]CatSnapshotResponse, 0, len(snapshots))
		for _, s := range snapshots {
			if strings.Contains(s.ID, snapshotName) {
				filtered = append(filtered, s)
			}
		}
		if len(filtered) == 0 {
			return nil, fmt.Errorf("snapshot not found: %s", snapshotName)
		}
		snapshots = filtered
	}

	return snapshots, nil
}
