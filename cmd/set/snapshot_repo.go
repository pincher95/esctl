package set

import (
	"github.com/pincher95/esctl/cmd/snapshot"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	setSnapshotRepoType     string
	setSnapshotRepoSettings string
)

var setSnapshotRepoCmd = &cobra.Command{
	Use:   "snapshot-repo <repository>",
	Short: "Create a snapshot repository",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Create a filesystem repository
	esctl set snapshot-repo my-repo --type=fs --settings="location:/backup/snapshots"
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return snapshot.HandleRepoCreate(cmd.Context(), args[0], setSnapshotRepoType, setSnapshotRepoSettings)
	},
}

func init() {
	setSnapshotRepoCmd.Flags().StringVar(&setSnapshotRepoType, "type", "fs", "Repository type (fs, s3, azure, gcs, etc.)")
	setSnapshotRepoCmd.Flags().StringVar(&setSnapshotRepoSettings, "settings", "", "Repository settings as key:value pairs separated by commas")
	setSnapshotRepoCmd.MarkFlagRequired("settings")
}
