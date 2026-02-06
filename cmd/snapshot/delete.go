package snapshot

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/es/snapshots"
)

func HandleSnapshotDelete(ctx context.Context, repository, snapshot string) error {
	if err := snapshots.DeleteSnapshot(ctx, repository, snapshot); err != nil {
		return fmt.Errorf("failed to delete snapshot: %w", err)
	}

	fmt.Printf("Snapshot '%s' deleted successfully from repository '%s'\n", snapshot, repository)
	return nil
}
