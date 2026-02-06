package set

import (
	"github.com/pincher95/esctl/cmd/snapshot"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/snapshots"
	"github.com/spf13/cobra"
)

var (
	setSnapshotIndices            string
	setSnapshotWait               bool
	setSnapshotIgnoreUnavailable  bool
	setSnapshotIncludeGlobalState *bool
	setSnapshotPartial            bool
)

var setSnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Create a new snapshot",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Create a snapshot of all indices
	esctl set snapshot --repository my-repo --name my-snapshot

	# Create a snapshot with specific indices
	esctl set snapshot --repository my-repo --name my-snapshot --indices="index1,index2"
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		request := snapshots.CreateSnapshotRequest{
			Indices:            setSnapshotIndices,
			IgnoreUnavailable:  setSnapshotIgnoreUnavailable,
			IncludeGlobalState: setSnapshotIncludeGlobalState,
			Partial:            setSnapshotPartial,
		}
		return snapshot.HandleSnapshotCreate(cmd.Context(), setSnapshotRepo, setSnapshotName, request, setSnapshotWait)
	},
}

func init() {
	setSnapshotCmd.Flags().StringVar(&setSnapshotRepo, "repository", "", "Snapshot repository name")
	setSnapshotCmd.Flags().StringVar(&setSnapshotName, "name", "", "Snapshot name")
	setSnapshotCmd.Flags().StringVar(&setSnapshotIndices, "indices", "", "Comma-separated list of indices to snapshot")
	setSnapshotCmd.Flags().BoolVar(&setSnapshotWait, "wait", false, "Wait for snapshot completion")
	setSnapshotCmd.Flags().BoolVar(&setSnapshotIgnoreUnavailable, "ignore-unavailable", false, "Ignore unavailable indices")
	setSnapshotCmd.Flags().BoolVar(&setSnapshotPartial, "partial", false, "Allow partial snapshots")

	var includeGlobalState bool
	setSnapshotCmd.Flags().BoolVar(&includeGlobalState, "include-global-state", true, "Include global cluster state")
	setSnapshotIncludeGlobalState = &includeGlobalState
	setSnapshotCmd.MarkFlagRequired("repository")
	setSnapshotCmd.MarkFlagRequired("name")
}

var (
	setSnapshotRepo string
	setSnapshotName string
)
