package snapshot

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/es/snapshots"
	"github.com/pincher95/esctl/output"
)

func HandleSnapshotStatus(ctx context.Context, repository, snapshot string) error {
	result, err := snapshots.SnapshotStatus(ctx, repository, snapshot)
	if err != nil {
		return fmt.Errorf("failed to get snapshot status: %w", err)
	}

	return output.Render(result)
}
