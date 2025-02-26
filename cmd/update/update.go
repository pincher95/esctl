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

}

func Cmd() *cobra.Command {
	return updateCmd
}

// func buildColumnDefs(columns []string, defaultColumns []output.ColumnDefaults) ([]output.ColumnDefaults, error) {
// 	columnDefs := make([]output.ColumnDefaults, 0, len(columns))
// 	for _, column := range columns {
// 		var found bool
// 		for _, defaultColumn := range defaultColumns {
// 			if strings.EqualFold(defaultColumn.Header, column) {
// 				columnDefs = append(columnDefs, defaultColumn)
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			return nil, fmt.Errorf("unknown column: %s", column)
// 		}
// 	}
// 	return columnDefs, nil
// }

// func getColumnDefs(conf config.Config, entity string, defaultColumns []output.ColumnDefaults) ([]output.ColumnDefaults, error) {
// 	if len(flagColumns) > 0 {
// 		for _, column := range flagColumns {
// 			if strings.EqualFold(column, "all") {
// 				return defaultColumns, nil
// 			}
// 		}
// 		return buildColumnDefs(flagColumns, defaultColumns)
// 	} else {
// 		entityConfig, ok := conf.Entities[entity]
// 		if !ok || len(entityConfig.Columns) == 0 {
// 			return defaultColumns, nil
// 		}
// 		return buildColumnDefs(entityConfig.Columns, defaultColumns)
// 	}
// }

// func clearScreen() {
// 	// Move cursor to top-left and clear screen
// 	fmt.Print("\033[?1049h\033[H\033[?25l")
// }
