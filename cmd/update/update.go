package update

import (
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Elasticsearch entities",
	Long: utils.Trim(`
The 'update' command allows you to update Elasticsearch entities.

Available Entities:
  - reroute: Changes the allocation of shards in a cluster.`),
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

}

func Cmd() *cobra.Command {
	return updateCmd
}
