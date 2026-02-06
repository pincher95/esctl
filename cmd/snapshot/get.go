package snapshot

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/es/snapshots"
	"github.com/pincher95/esctl/output"
)

func HandleSnapshotGet(ctx context.Context, repository, snapshotName string) error {
	result, err := snapshots.ListSnapshots(ctx, repository)
	if err != nil {
		return err
	}

	for _, snap := range result.Snapshots {
		if snap.Snapshot == snapshotName {
			return output.Render(snap)
		}
	}

	return fmt.Errorf("snapshot not found: %s/%s", repository, snapshotName)
}
