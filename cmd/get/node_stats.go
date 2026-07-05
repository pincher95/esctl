package get

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/node"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

var flagNodeStatsNode string

var getNodeStatsCmd = &cobra.Command{
	Use:   "node-stats",
	Short: "Get node health stats (heap, GC, thread-pool rejections, breakers, disk)",
	Long: utils.Trim(`
Show the health-critical node statistics that _cat/nodes omits: JVM heap pressure, garbage-collection
counts, thread-pool rejections (search/write), tripped circuit breakers, and available disk. Rising
rejections or tripped breakers are the earliest signal of an overloaded node.`),
	Example: utils.TrimAndIndent(`
# Stats for all nodes.
esctl get node-stats

# Stats for a single node, sorted by heap usage.
esctl get node-stats --node es-data-0 --sort-by HEAP%:desc

# Watch rejections build up.
esctl get node-stats --watch --interval 5s
`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if !flagRefresh {
			return handleNodeStatsLogic(ctx)
		}
		return utils.WatchLoopContext(ctx, flagRefreshInterval, func() error {
			return handleNodeStatsLogic(ctx)
		})
	},
}

func init() {
	getNodeStatsCmd.Flags().StringVar(&flagNodeStatsNode, "node", "", "Node id/name to filter by (default: all nodes)")
	getCmd.AddCommand(getNodeStatsCmd)
}

var nodeStatsColumns = []output.ColumnDefaults{
	{Header: "NODE", Type: output.Text},
	{Header: "HEAP%", Type: output.Percent},
	{Header: "HEAP-USED", Type: output.DataSize},
	{Header: "GC-YOUNG", Type: output.Number},
	{Header: "GC-OLD", Type: output.Number},
	{Header: "SEARCH-REJECTED", Type: output.Number},
	{Header: "WRITE-REJECTED", Type: output.Number},
	{Header: "BREAKERS-TRIPPED", Type: output.Number},
	{Header: "DISK-AVAIL", Type: output.DataSize},
}

func handleNodeStatsLogic(ctx context.Context) error {
	resp, err := node.GetNodeStats(ctx, flagNodeStatsNode)
	if err != nil {
		return fmt.Errorf("failed to retrieve node stats: %w", err)
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(resp)
	}

	data := make([][]string, 0, len(resp.Nodes))
	for _, n := range resp.Nodes {
		var tripped int64
		for _, b := range n.Breakers {
			tripped += b.Tripped
		}

		rowData := map[string]string{
			"NODE":             n.Name,
			"HEAP%":            fmt.Sprintf("%d%%", n.JVM.Mem.HeapUsedPercent),
			"HEAP-USED":        humanizeBytes(n.JVM.Mem.HeapUsedInBytes),
			"GC-YOUNG":         strconv.FormatInt(n.JVM.GC.Collectors.Young.CollectionCount, 10),
			"GC-OLD":           strconv.FormatInt(n.JVM.GC.Collectors.Old.CollectionCount, 10),
			"SEARCH-REJECTED":  strconv.FormatInt(threadPoolRejected(n, "search"), 10),
			"WRITE-REJECTED":   strconv.FormatInt(threadPoolRejected(n, "write", "bulk"), 10),
			"BREAKERS-TRIPPED": strconv.FormatInt(tripped, 10),
			"DISK-AVAIL":       humanizeBytes(n.FS.Total.AvailableInBytes),
		}
		row := make([]string, len(nodeStatsColumns))
		for i, colDef := range nodeStatsColumns {
			row[i] = rowData[colDef.Header]
		}
		data = append(data, row)
	}

	sortBy := flagSortBy
	if sortBy == "" {
		sortBy = "NODE"
	}
	return output.PrintTable(nodeStatsColumns, data, output.ParseSortColumns(sortBy))
}

// threadPoolRejected returns the rejected count for the first thread pool present
// among the given names (e.g. "write" then legacy "bulk").
func threadPoolRejected(n node.NodeStats, names ...string) int64 {
	for _, name := range names {
		if pool, ok := n.ThreadPool[name]; ok {
			return pool.Rejected
		}
	}
	return 0
}

// humanizeBytes formats a byte count using the same lowercase unit suffixes as
// the _cat APIs (b, kb, mb, gb, tb, pb).
func humanizeBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%db", b)
	}
	suffixes := []string{"kb", "mb", "gb", "tb", "pb"}
	value := float64(b)
	i := -1
	for value >= unit && i < len(suffixes)-1 {
		value /= unit
		i++
	}
	return fmt.Sprintf("%.1f%s", value, suffixes[i])
}
