package snapshot

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/es/snapshots"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
)

func HandleSnapshotGet(ctx context.Context, repository, snapshotName string) error {
	result, err := snapshots.ListSnapshots(ctx, repository)
	if err != nil {
		return err
	}

	matches := make([]snapshots.SnapshotInfo, 0)
	for _, snap := range result.Snapshots {
		if strings.Contains(snap.Snapshot, snapshotName) {
			matches = append(matches, snap)
		}
	}

	if len(matches) == 0 {
		return fmt.Errorf("snapshot not found: %s/%s", repository, snapshotName)
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		if len(matches) == 1 {
			return output.Render(matches[0])
		}
		return output.Render(matches)
	}

	if len(matches) == 1 {
		// Single snapshot: key-value table
		snap := matches[0]
		columnDefs := []output.ColumnDefaults{
			{Header: "SETTING", Type: output.Text},
			{Header: "VALUE", Type: output.Text},
		}
		data := [][]string{
			{"snapshot", snap.Snapshot},
			{"uuid", snap.UUID},
			{"state", snap.State},
			{"version", snap.Version},
			{"start_time", snap.StartTime},
			{"end_time", snap.EndTime},
			{"duration_ms", fmt.Sprintf("%d", snap.DurationInMillis)},
			{"indices", fmt.Sprintf("%d", len(snap.Indices))},
			{"shards.total", fmt.Sprintf("%d", snap.Shards.Total)},
			{"shards.successful", fmt.Sprintf("%d", snap.Shards.Successful)},
			{"shards.failed", fmt.Sprintf("%d", snap.Shards.Failed)},
		}
		if len(snap.Failures) > 0 {
			data = append(data, []string{"failures", fmt.Sprintf("%d", len(snap.Failures))})
		}
		return output.PrintTable(columnDefs, data, nil)
	}

	// Multiple matches: list table
	columnDefs := []output.ColumnDefaults{
		{Header: "SNAPSHOT", Type: output.Text},
		{Header: "STATE", Type: output.Text},
		{Header: "START_TIME", Type: output.Text},
		{Header: "END_TIME", Type: output.Text},
		{Header: "DURATION_MS", Type: output.Number},
		{Header: "INDICES", Type: output.Number},
		{Header: "SHARDS_OK", Type: output.Number},
		{Header: "SHARDS_FAIL", Type: output.Number},
	}

	data := make([][]string, 0, len(matches))
	for _, snap := range matches {
		data = append(data, []string{
			snap.Snapshot,
			snap.State,
			snap.StartTime,
			snap.EndTime,
			fmt.Sprintf("%d", snap.DurationInMillis),
			fmt.Sprintf("%d", len(snap.Indices)),
			fmt.Sprintf("%d", snap.Shards.Successful),
			fmt.Sprintf("%d", snap.Shards.Failed),
		})
	}

	return output.PrintTable(columnDefs, data, nil)
}
