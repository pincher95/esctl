package snapshot

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/es/snapshots"
	"github.com/pincher95/esctl/internal/validation"
)

func HandleSnapshotRestore(ctx context.Context, repository, snapshot string, request snapshots.RestoreSnapshotRequest, wait bool) error {
	if request.Indices != "" {
		for _, idx := range strings.Split(request.Indices, ",") {
			clean := strings.TrimSpace(idx)
			if clean == "" {
				continue
			}
			if err := validation.ValidateIndexPattern(clean); err != nil {
				return err
			}
		}
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
