package cat

import (
	"fmt"
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

func (c *cat) CatSnapshots(endpoint, repository, snapshotName *string) (*[]CatSnapshotResponse, error) {
	if endpoint == nil {
		endpoint = new(string)
		*endpoint = "_cat/snapshots?format=json"
	}

	if *repository != "" {
		*endpoint = fmt.Sprintf("_cat/snapshots/%s?format=json", *repository)
	}

	snapshots := make([]CatSnapshotResponse, 0)

	resp, err := shared.Client.R().SetHeader("Content-Type", "application/json").SetResult(&snapshots).Get(*endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get snapshots: %s", resp.Status())
	}

	if snapshotName != nil {
		filtered := make([]CatSnapshotResponse, 0, len(snapshots))

		for _, snapshot := range snapshots {
			if strings.Contains(snapshot.ID, *snapshotName) {
				filtered = append(filtered, snapshot)
			}
		}

		// If no matches found, return error
		if len(filtered) == 0 {
			return nil, fmt.Errorf("node not found: %s", *snapshotName)
		}

		// Replace the original slice with the filtered one
		snapshots = filtered
	}

	return &snapshots, nil
}
