package snapshot

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/snapshots"
	"github.com/spf13/cobra"
)

var (
	createIndices            string
	createWait               bool
	createIgnoreUnavailable  bool
	createIncludeGlobalState *bool
	createPartial            bool
)

var createCmd = &cobra.Command{
	Use:   "create <repository> <snapshot>",
	Short: "Create a new snapshot",
	Args:  cobra.ExactArgs(2),
	Example: utils.TrimAndIndent(`
	# Create a snapshot of all indices
	esctl snapshot create my-repo my-snapshot

	# Create a snapshot with specific indices
	esctl snapshot create my-repo my-snapshot --indices="index1,index2"

	# Create a snapshot and wait for completion
	esctl snapshot create my-repo my-snapshot --wait

	# Create a snapshot without global state
	esctl snapshot create my-repo my-snapshot --include-global-state=false
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleSnapshotCreate(ctx, args[0], args[1])
	},
}

func init() {
	createCmd.Flags().StringVar(&createIndices, "indices", "", "Comma-separated list of indices to snapshot")
	createCmd.Flags().BoolVar(&createWait, "wait", false, "Wait for snapshot completion")
	createCmd.Flags().BoolVar(&createIgnoreUnavailable, "ignore-unavailable", false, "Ignore unavailable indices")
	createCmd.Flags().BoolVar(&createPartial, "partial", false, "Allow partial snapshots")

	// Use a pointer to bool to distinguish between unset and false
	var includeGlobalState bool
	createCmd.Flags().BoolVar(&includeGlobalState, "include-global-state", true, "Include global cluster state")
	createIncludeGlobalState = &includeGlobalState
}

func handleSnapshotCreate(ctx context.Context, repository, snapshot string) error {
	request := snapshots.CreateSnapshotRequest{
		Indices:            createIndices,
		IgnoreUnavailable:  createIgnoreUnavailable,
		IncludeGlobalState: createIncludeGlobalState,
		Partial:            createPartial,
	}

	err := snapshots.CreateSnapshot(ctx, repository, snapshot, request, createWait)
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}

	if createWait {
		fmt.Printf("Snapshot '%s' created successfully in repository '%s'\n", snapshot, repository)
	} else {
		fmt.Printf("Snapshot '%s' creation started in repository '%s'\n", snapshot, repository)
	}

	return nil
}
