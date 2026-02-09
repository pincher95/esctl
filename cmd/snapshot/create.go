package snapshot

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/snapshots"
)

func HandleSnapshotCreate(ctx context.Context, repository, snapshot string, request snapshots.CreateSnapshotRequest, wait bool) error {
	if err := utils.ValidateIndexPatternsCSV(request.Indices); err != nil {
		return err
	}

	if err := snapshots.CreateSnapshot(ctx, repository, snapshot, request, wait); err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}

	if wait {
		fmt.Printf("Snapshot '%s' created successfully in repository '%s'\n", snapshot, repository)
	} else {
		fmt.Printf("Snapshot '%s' creation started in repository '%s'\n", snapshot, repository)
	}

	return nil
}
