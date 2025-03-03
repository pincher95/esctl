package get

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/constants"
	cat "github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

var getShardsCmd = &cobra.Command{
	Use:   "shards",
	Short: "Get shards in Elasticsearch cluster",
	Long: utils.Trim(`
The 'shards' command provides detailed information about each shard in the Elasticsearch cluster.

This includes:
  - Shard number
  - State of the shard (e.g., whether it's started, relocating, initializing, or unassigned)
  - Whether the shard is a primary or a replica
  - Size of the shard
  - Node on which the shard is located

Filters can be applied to only show shards in certain states, with a specific number, located on a particular node, or designated as primary or replica.`),
	Example: utils.TrimAndIndent(`
# Retrieve detailed information about shards in the Elasticsearch cluster.
esctl get shards

# Retrieve shard information for an index.
esctl get shards --index my_index

# Retrieve shard information filtered by state.
esctl get shards --started --relocating`),
	Run: func(cmd *cobra.Command, args []string) {
		config := config.ParseConfigFile()

		// If --watch is NOT set, just run once
		if !flagRefresh {
			handleShardLogic(*config)
			return
		}

		// If --watch is set, run in a loop
		for {
			clearScreen() // optional, to mimic "watch" clearing
			handleShardLogic(*config)
			time.Sleep(flagRefreshInterval)
		}
	},
}

func init() {
	getShardsCmd.Flags().StringVarP(&flagIndex, "index", "i", "", "Name of the index")
	getShardsCmd.Flags().StringVar(&flagNode, "node", "", "Filter shards by node")
	getShardsCmd.Flags().IntVar(&flagShard, "shard", -1, "Filter shards by shard number")
	getShardsCmd.Flags().BoolVar(&flagPrimary, "primary", false, "Filter primary shards")
	getShardsCmd.Flags().BoolVar(&flagReplica, "replica", false, "Filter replica shards")
	getShardsCmd.Flags().BoolVar(&flagStarted, "started", false, "Filter shards in STARTED state")
	getShardsCmd.Flags().BoolVar(&flagRelocating, "relocating", false, "Filter shards in RELOCATING state")
	getShardsCmd.Flags().BoolVar(&flagInitializing, "initializing", false, "Filter shards in INITIALIZING state")
	getShardsCmd.Flags().BoolVar(&flagUnassigned, "unassigned", false, "Filter shards in UNASSIGNED state")
}

func includeShardByState(shard cat.Shard) bool {
	switch {
	case flagStarted && shard.State == constants.ShardStateStarted:
		return true
	case flagRelocating && shard.State == constants.ShardStateRelocating:
		return true
	case flagInitializing && shard.State == constants.ShardStateInitializing:
		return true
	case flagUnassigned && shard.State == constants.ShardStateUnassigned:
		return true
	case !flagStarted && !flagRelocating && !flagInitializing && !flagUnassigned:
		return true
	}
	return false
}

func includeShardByNumber(shard cat.Shard) bool {
	return flagShard == -1 || flagShard == shard.Shard
}

func includeShardByPriRep(shard cat.Shard) bool {
	return (flagPrimary && shard.Prirep == constants.ShardPrimary) ||
		(flagReplica && shard.Prirep == constants.ShardReplica) ||
		(!flagPrimary && !flagReplica)
}

func includeShardByNode(shard cat.Shard) bool {
	if flagNode == "" {
		return true
	}

	return utils.SafeString(shard.Node) == flagNode
}

func humanizePriRep(priRep string) string {
	switch priRep {
	case constants.ShardPrimary:
		return "primary"
	case constants.ShardReplica:
		return "replica"
	default:
		return priRep
	}
}

var shardColumns = []output.ColumnDefaults{
	{Header: "INDEX", Type: output.Text},
	{Header: "SHARD", Type: output.Number},
	{Header: "PRI-REP", Type: output.Text},
	{Header: "STATE", Type: output.Text},
	{Header: "DOCS", Type: output.Number},
	{Header: "STORE", Type: output.DataSize},
	{Header: "IP", Type: output.Text},
	{Header: "NODE", Type: output.Text},
}

func handleShardLogic(conf config.Config) {
	shards, err := cat.CatShards(nil, &flagIndex, nil, nil, shared.Debug)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve shards:", err)
		os.Exit(1)
	}

	columnDefs, err := getColumnDefs(conf, "shard", shardColumns)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get column definitions:", err)
		os.Exit(1)
	}

	data := [][]string{}

	for _, shard := range shards {
		if includeShardByState(shard) && includeShardByNumber(shard) &&
			includeShardByPriRep(shard) && includeShardByNode(shard) {

			rowData := map[string]string{
				"INDEX":   shard.Index,
				"SHARD":   strconv.Itoa(shard.Shard),
				"PRI-REP": humanizePriRep(shard.Prirep),
				"STATE":   shard.State,
				"DOCS":    utils.SafeString(shard.Docs),
				"STORE":   utils.SafeString(shard.Store),
				"IP":      utils.SafeString(shard.IP),
				"NODE":    utils.SafeString(shard.Node),
			}

			row := make([]string, len(columnDefs))
			for i, colDef := range columnDefs {
				row[i] = rowData[colDef.Header]
			}
			data = append(data, row)
		}
	}

	if len(flagSortBy) > 0 {
		sortCols := output.ParseSortColumns(flagSortBy)
		output.PrintTable(columnDefs, data, sortCols)
	} else {
		sortCols := output.ParseSortColumns("SHARD")
		output.PrintTable(columnDefs, data, sortCols)
	}
}
