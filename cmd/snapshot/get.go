package snapshot

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/es/snapshots"
	"github.com/pincher95/esctl/output"
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
	if len(matches) == 1 {
		return output.Render(matches[0])
	}
	return output.Render(matches)
}
