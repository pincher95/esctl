package get

import (
	"fmt"
	"strings"
	"time"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagActions             []string
	flagColumns             []string
	flagIndex               string
	flagNode                string
	flagNodeID              string
	flagSortBy              string
	flagBytes               string
	flagTime                string
	flagRepository          string
	flagFields              string
	flagStatus              string
	flagFilter              string
	flagTasksID             string
	flagRefreshInterval     time.Duration
	flagShard               int
	flagInitializing        bool
	flagPrimary             bool
	flagRelocating          bool
	flagReplica             bool
	flagStarted             bool
	flagUnassigned          bool
	flagRefresh             bool
	flagIncludeDiskInfo     bool
	flagIncludeYesDecisions bool
	flagWritable            bool
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Elasticsearch entities",
	Long: utils.Trim(`
The 'get' command allows you to retrieve information about Elasticsearch entities.

Available Entities:
  - nodes: List all nodes in the Elasticsearch cluster.
  - indices: List all indices in the Elasticsearch cluster.
  - shards: List detailed information about shards, including their sizes and placement.
  - aliases: List all aliases in the Elasticsearch cluster.
  - alias: Get details of a specific alias.
  - tasks: List all tasks in the Elasticsearch cluster.
	- allocation: List allocation in the Elasticsearch cluster.
	- plugins: List all plugins in the Elasticsearch cluster.
	- explain: List allocation explain in the Elasticsearch cluster.
  - snapshot: Get snapshot details or list snapshots.
  - snapshot-status: Show snapshot status.
  - snapshot-repo: Get snapshot repository details.
	- fielddata: List fielddata in the Elasticsearch cluster.
	- health: Cluster health overview.
  - pipeline: Get ingest pipeline details.
  - user: Get user details.
  - role: Get role details.
  - reindex: Get reindex task status.
	`),
	// 	Example: utils.TrimAndIndent(`
	// #Retrieve a list of all nodes in the Elasticsearch cluster.
	// esctl get nodes

	// #Retrieve a list of all indices in the Elasticsearch cluster.
	// esctl get indices

	// #Retrieve detailed information about shards in the Elasticsearch cluster.
	// esctl get shards

	// #Retrieve shard information for an index.
	// esctl get shards --index my_index

	// #Retrieve shard information filtered by state.
	// esctl get shards --started --relocating

	// #Retrieve all aliases.
	// esctl get aliases

	// #Retrieve tasks filtered by actions using wildcard patterns.
	// esctl get tasks --actions 'index*' --actions '*search*'

	// #Retrieve all tasks.
	// esctl get tasks`),
}

func init() {
	getCmd.PersistentFlags().StringVarP(&flagSortBy, "sort-by", "s", "", "Columns to sort by (comma-separated), e.g. 'NAME:desc,HEAP-PERCENT:asc'")
	getCmd.PersistentFlags().StringSliceVarP(&flagColumns, "columns", "c", []string{}, "Columns to display (comma-separated) or 'all'")
	getCmd.PersistentFlags().BoolVarP(&flagRefresh, "watch", "w", false, "Continuously watch the output")
	getCmd.PersistentFlags().DurationVar(&flagRefreshInterval, "interval", 5*time.Second, "Interval between consecutive fetches")

	getCmd.AddCommand(getAliasesCmd)
	getCmd.AddCommand(getIndicesCmd)
	getCmd.AddCommand(getNodesCmd)
	getCmd.AddCommand(getShardsCmd)
	getCmd.AddCommand(getTasksCmd)
	getCmd.AddCommand(getAllocationCmd)
	getCmd.AddCommand(getPluginsCmd)
	getCmd.AddCommand(getAllocationExplainCmd)
	getCmd.AddCommand(getFielddataCmd)
	getCmd.AddCommand(getHealthCmd)
	getCmd.AddCommand(getAliasCmd)
	getCmd.AddCommand(getPipelineCmd)
	getCmd.AddCommand(getPipelinesCmd)
	getCmd.AddCommand(getSnapshotCmd)
	getCmd.AddCommand(getSnapshotStatusCmd)
	getCmd.AddCommand(getSnapshotRepoCmd)
	getCmd.AddCommand(getSnapshotReposCmd)
	getCmd.AddCommand(getUserCmd)
	getCmd.AddCommand(getUsersCmd)
	getCmd.AddCommand(getRoleCmd)
	getCmd.AddCommand(getRolesCmd)
	getCmd.AddCommand(getReindexCmd)
}

func Cmd() *cobra.Command {
	return getCmd
}

func buildColumnDefs(columns []string, defaultColumns []output.ColumnDefaults) ([]output.ColumnDefaults, error) {
	columnDefs := make([]output.ColumnDefaults, 0, len(columns))
	for _, column := range columns {
		var found bool
		for _, defaultColumn := range defaultColumns {
			if strings.EqualFold(defaultColumn.Header, column) {
				columnDefs = append(columnDefs, defaultColumn)
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("unknown column: %s", column)
		}
	}
	return columnDefs, nil
}

func getColumnDefs(conf config.Config, entity string, defaultColumns []output.ColumnDefaults) ([]output.ColumnDefaults, error) {
	if len(flagColumns) > 0 {
		for _, column := range flagColumns {
			if strings.EqualFold(column, "all") {
				return defaultColumns, nil
			}
		}
		return buildColumnDefs(flagColumns, defaultColumns)
	} else {
		entityConfig, ok := conf.Entities[entity]
		if !ok || len(entityConfig.Columns) == 0 {
			return defaultColumns, nil
		}
		return buildColumnDefs(entityConfig.Columns, defaultColumns)
	}
}

func clearScreen() {
	// Move cursor to top-left and clear screen
	fmt.Print("\033[?1049h\033[H\033[?25l")
}
