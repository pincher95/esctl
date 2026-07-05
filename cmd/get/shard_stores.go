package get

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

var (
	flagShardStoresIndices string
	flagShardStoresStatus  string
)

var getShardStoresCmd = &cobra.Command{
	Use:                   "shard-stores [--indices index] [--status green,yellow,red,all]",
	DisableFlagsInUseLine: true,
	Short:                 "Get shard store information to diagnose unassigned shards",
	Long: utils.Trim(`
The 'shard-stores' command reports, for every on-disk copy of a shard, which node holds it,
its allocation id, whether the copy is currently used as a primary or replica (or is unused),
and any store-level exception such as corruption.

This is the primary tool for diagnosing why a shard cannot be allocated: it shows which nodes
still physically hold a copy of an unassigned shard, so you can decide whether to wait for a
node to rejoin, or force allocation with 'esctl update reroute'.

By default only shards with at least one unassigned copy are returned (status 'yellow,red').
Use --status all to list stores for every shard. Empty columns (e.g. STORE-EXCEPTION when no
copy is corrupt) are hidden automatically.`),
	Example: utils.TrimAndIndent(`
# Show stores for shards that are yellow or red (the default).
esctl get shard-stores

# Show stores for specific indices.
esctl get shard-stores --indices my-index,other-index

# Only shards with an unassigned primary (red).
esctl get shard-stores --status red

# Show stores for all shards, sorted by node.
esctl get shard-stores --status all --sort-by NODE

# Full detail as JSON, including store exceptions.
esctl get shard-stores --indices my-index -o json
`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		idxClient := index.NewIndex()
		conf, err := config.ParseConfigFile()
		if err != nil {
			return err
		}

		if !flagRefresh {
			return handleShardStoresLogic(ctx, idxClient, *conf)
		}
		return utils.WatchLoopContext(ctx, flagRefreshInterval, func() error {
			return handleShardStoresLogic(ctx, idxClient, *conf)
		})
	},
}

func init() {
	getShardStoresCmd.Flags().StringVar(&flagShardStoresIndices, "indices", "", "Comma-separated list of indices (default: all)")
	getShardStoresCmd.Flags().StringVar(&flagShardStoresStatus, "status", "", "Filter shards by health: comma-separated subset of green,yellow,red,all (default: yellow,red)")
}

var shardStoresColumns = []output.ColumnDefaults{
	{Header: "INDEX", Type: output.Text},
	{Header: "SHARD", Type: output.Number},
	{Header: "ALLOCATION", Type: output.Text},
	{Header: "NODE", Type: output.Text},
	{Header: "NODE-ID", Type: output.Text},
	{Header: "ALLOCATION-ID", Type: output.Text},
	{Header: "STORE-EXCEPTION", Type: output.Text},
}

func handleShardStoresLogic(ctx context.Context, client index.Index, conf config.Config) error {
	var indices []string
	if flagShardStoresIndices != "" {
		var err error
		indices, err = utils.ParseIndexPatternsCSV(flagShardStoresIndices, false)
		if err != nil {
			return err
		}
	}

	stores, err := client.GetShardStores(ctx, indices, flagShardStoresStatus)
	if err != nil {
		return fmt.Errorf("failed to retrieve shard stores: %w", err)
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(stores)
	}

	columnDefs, err := getColumnDefs(conf, "shard-stores", shardStoresColumns)
	if err != nil {
		return fmt.Errorf("failed to get column definitions: %w", err)
	}

	data := make([][]string, 0)
	for indexName, idx := range stores.Indices {
		for shardID, shard := range idx.Shards {
			for _, store := range shard.Stores {
				rowData := map[string]string{
					"INDEX":           indexName,
					"SHARD":           shardID,
					"ALLOCATION":      store.Allocation,
					"NODE":            shardStoreNode(store),
					"NODE-ID":         store.NodeID,
					"ALLOCATION-ID":   store.AllocationID,
					"STORE-EXCEPTION": shardStoreException(store),
				}

				row := make([]string, len(columnDefs))
				for i, colDef := range columnDefs {
					row[i] = rowData[colDef.Header]
				}
				data = append(data, row)
			}
		}
	}

	sortBy := flagSortBy
	if sortBy == "" {
		// Map iteration order is non-deterministic; default to a stable ordering.
		sortBy = "INDEX,SHARD"
	}
	return output.PrintTable(columnDefs, data, output.ParseSortColumns(sortBy))
}

// shardStoreNode returns a human-friendly node label, falling back to the node id.
func shardStoreNode(store index.ShardStore) string {
	if store.NodeName != "" {
		return store.NodeName
	}
	return store.NodeID
}

// shardStoreException returns a concise store-exception label (empty when healthy).
func shardStoreException(store index.ShardStore) string {
	if store.StoreException == nil {
		return ""
	}
	if store.StoreException.Type != "" {
		return store.StoreException.Type
	}
	return store.StoreException.Reason
}
