package snapshot

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/es/snapshots"
	"github.com/pincher95/esctl/internal/validation"
)

func HandleSnapshotCreate(ctx context.Context, repository, snapshot string, request snapshots.CreateSnapshotRequest, wait bool) error {
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
