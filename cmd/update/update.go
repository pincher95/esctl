package update

import (
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	flagMertic       string
	flagDryRun       bool
	flagExplain      bool
	flagRetryFailed  bool
	flagFlatBody     string
	flagSettings     string
	flagFlatSettings bool
	flagIndex        string
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Elasticsearch entities",
	Long: utils.Trim(`
The 'update' command allows you to update Elasticsearch entities.

Available Entities:
  - reroute: Changes the allocation of shards in a cluster.
  - alias: Move an alias between indices.
  - pipeline: Simulate pipeline execution.
  - snapshot: Restore a snapshot.`),
	Example: utils.TrimAndIndent(`
# Reroute the shards in the cluster.
esctl update reroute
	`),
}

func init() {
	// updateCmd.PersistentFlags().StringVarP(&flagSortBy, "sort-by", "s", "", "Columns to sort by (comma-separated), e.g. 'NAME:desc,HEAP-PERCENT:asc'")
	// updateCmd.PersistentFlags().StringSliceVarP(&flagColumns, "columns", "c", []string{}, "Columns to display (comma-separated) or 'all'")
	// updateCmd.PersistentFlags().BoolVarP(&flagRefresh, "watch", "w", false, "Continuously watch the output")
	// updateCmd.PersistentFlags().DurationVar(&flagRefreshInterval, "interval", 5*time.Second, "Interval between consecutive fetches")

	updateCmd.AddCommand(updateRerouteCmd)
	updateCmd.AddCommand(updateIndexCmd)
	updateCmd.AddCommand(updateCacheClearCmd)
	updateCmd.AddCommand(updateAliasCmd)
	updateCmd.AddCommand(updateSnapshotCmd)
	updateCmd.AddCommand(updatePipelineCmd)
	updateCmd.AddCommand(updateDataStreamCmd)
	updateCmd.AddCommand(updateSearchTemplateCmd)
}

func Cmd() *cobra.Command {
	return updateCmd
}
