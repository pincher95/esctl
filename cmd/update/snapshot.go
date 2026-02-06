package update

import (
	"fmt"
	"strings"

	"github.com/pincher95/esctl/cmd/snapshot"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/snapshots"
	"github.com/spf13/cobra"
)

var (
	updateSnapshotIndices             string
	updateSnapshotWait                bool
	updateSnapshotIgnoreUnavailable   bool
	updateSnapshotIncludeGlobalState  bool
	updateSnapshotRenamePattern       string
	updateSnapshotRenameReplacement   string
	updateSnapshotIncludeAliases      bool
	updateSnapshotIndexSettings       string
	updateSnapshotIgnoreIndexSettings string
)

var updateSnapshotCmd = &cobra.Command{
	Use:   "snapshot <repository> <snapshot>",
	Short: "Restore a snapshot",
	Args:  cobra.ExactArgs(2),
	Example: utils.TrimAndIndent(`
	# Restore all indices from a snapshot
	esctl update snapshot my-repo my-snapshot

	# Restore specific indices
	esctl update snapshot my-repo my-snapshot --indices="index1,index2"
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		request := snapshots.RestoreSnapshotRequest{
			Indices:            updateSnapshotIndices,
			IgnoreUnavailable:  updateSnapshotIgnoreUnavailable,
			IncludeGlobalState: updateSnapshotIncludeGlobalState,
			RenamePattern:      updateSnapshotRenamePattern,
			RenameReplacement:  updateSnapshotRenameReplacement,
			IncludeAliases:     updateSnapshotIncludeAliases,
		}

		if updateSnapshotIndexSettings != "" {
			settings, err := snapshot.ParseSettings(updateSnapshotIndexSettings)
			if err != nil {
				return fmt.Errorf("invalid index settings format: %w", err)
			}
			request.IndexSettings = settings
		}

		if updateSnapshotIgnoreIndexSettings != "" {
			request.IgnoreIndexSettings = strings.Split(updateSnapshotIgnoreIndexSettings, ",")
			for i, setting := range request.IgnoreIndexSettings {
				request.IgnoreIndexSettings[i] = strings.TrimSpace(setting)
			}
		}

		return snapshot.HandleSnapshotRestore(cmd.Context(), args[0], args[1], request, updateSnapshotWait)
	},
}

func init() {
	updateSnapshotCmd.Flags().StringVar(&updateSnapshotIndices, "indices", "", "Comma-separated list of indices to restore")
	updateSnapshotCmd.Flags().BoolVar(&updateSnapshotWait, "wait", false, "Wait for restore completion")
	updateSnapshotCmd.Flags().BoolVar(&updateSnapshotIgnoreUnavailable, "ignore-unavailable", false, "Ignore unavailable indices")
	updateSnapshotCmd.Flags().BoolVar(&updateSnapshotIncludeGlobalState, "include-global-state", false, "Include global cluster state")
	updateSnapshotCmd.Flags().StringVar(&updateSnapshotRenamePattern, "rename-pattern", "", "Rename pattern for restored indices")
	updateSnapshotCmd.Flags().StringVar(&updateSnapshotRenameReplacement, "rename-replacement", "", "Rename replacement for restored indices")
	updateSnapshotCmd.Flags().BoolVar(&updateSnapshotIncludeAliases, "include-aliases", true, "Include aliases when restoring")
	updateSnapshotCmd.Flags().StringVar(&updateSnapshotIndexSettings, "index-settings", "", "Index settings to apply during restore (key:value pairs)")
	updateSnapshotCmd.Flags().StringVar(&updateSnapshotIgnoreIndexSettings, "ignore-index-settings", "", "Comma-separated list of index settings to ignore")
}
