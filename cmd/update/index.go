package update

import (
	"github.com/pincher95/esctl/cmd/update/index"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var updateIndexCmd = &cobra.Command{
	Use:   "index",
	Short: "Updates the index settings",
	Long: utils.Trim(`
	Updates the settings of an index. You can update the number of replicas, the number of shards, and other settings.
	`),
	Example: utils.TrimAndIndent(`
	# Update the number of replicas for an index.
	esctl update index --index my_index --replicas 2
	`),
}

func init() {
	updateIndexCmd.PersistentFlags().BoolVar(&flagFlatSettings, "no-flat-settings", true, "If set, print settings in a none flat format (Default is false)")
	updateIndexCmd.PersistentFlags().StringVarP(&flagIndex, "index", "i", "", "Name of the index")
	// updateIndexCmd.Flags().StringVar(&flagFlatBody, "flat-settings", "", "Index settings")

	// Mark name as required
	_ = updateIndexCmd.MarkPersistentFlagRequired("index")

	updateIndexCmd.AddCommand(index.SettingsCmd)
}
