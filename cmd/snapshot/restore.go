package snapshot

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/snapshots"
)

func HandleSnapshotRestore(ctx context.Context, repository, snapshot string, request snapshots.RestoreSnapshotRequest, wait bool) error {
	if err := utils.ValidateIndexPatternsCSV(request.Indices); err != nil {
		return err
	}

	if err := snapshots.RestoreSnapshot(ctx, repository, snapshot, request, wait); err != nil {
		return fmt.Errorf("failed to restore snapshot: %w", err)
	}

	if wait {
		fmt.Printf("Snapshot '%s' restored successfully from repository '%s'\n", snapshot, repository)
	} else {
		fmt.Printf("Snapshot '%s' restore started from repository '%s'\n", snapshot, repository)
	}

	return nil
}
