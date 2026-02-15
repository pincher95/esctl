package get

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"

	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

var getAllocationCmd = &cobra.Command{
	Use:                   "allocation [--node-id node-id] [--bytes bytes]",
	DisableFlagsInUseLine: true,
	Short:                 "Get Elasticsearch allocation for the Elasticsearch cluster",
	Long: utils.Trim(`
	Get Elasticsearch allocation. You can filter the results using the node flag.
	`),
	Example: utils.TrimAndIndent(`
	# Retrieve all allocation.
	esctl get allocation

	# Retrieve allocation for a specific node.
	esctl get allocation --node my_node

	# Retrieve allocation in kilobytes.
	esctl get allocation --bytes kb
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		allocationClient := cat.NewCat()
		conf, err := config.ParseConfigFile()
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		if !flagRefresh {
			return handleAllocationLogic(ctx, allocationClient, *conf)
		}

		return utils.WatchLoopContext(ctx, flagRefreshInterval, func() error {
			return handleAllocationLogic(ctx, allocationClient, *conf)
		})
	},
}

func init() {
	getAllocationCmd.Flags().StringVar(&flagNodeID, "node-id", "", "A comma-separated list of node identifiers or names used to limit the returned information.")
	getAllocationCmd.Flags().StringVar(&flagBytes, "bytes", "", "The unit in which to display byte values. Valid values are: 'b', 'kb', 'mb', 'gb', 'tb', 'pb'.")
}

var allocationColumns = []output.ColumnDefaults{
	{Header: "SHARDS", Type: output.Number},
	{Header: "DISK-INDICES", Type: output.Text},
	{Header: "DISK-USED", Type: output.Text},
	{Header: "DISK-AVAIL", Type: output.Text},
	{Header: "DISK-TOTAL", Type: output.Number},
	{Header: "DISK-PERCENT", Type: output.Number},
	{Header: "HOST", Type: output.Text},
	{Header: "IP", Type: output.Text},
	{Header: "NODE", Type: output.Text},
}

func handleAllocationLogic(ctx context.Context, client cat.Cat, conf config.Config) error {
	allocations, err := client.CatAllocation(ctx, "", flagNodeID, flagBytes)
	if err != nil {
		return fmt.Errorf("failed to retrieve allocation: %w", err)
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(allocations)
	}

	columnDefs, err := getColumnDefs(conf, "shards", allocationColumns)
	if err != nil {
		return fmt.Errorf("failed to get column definitions: %w", err)
	}

	data := make([][]string, 0, len(allocations))

	for _, allocation := range allocations {
		rowData := map[string]string{
			"SHARDS":       strconv.Itoa(allocation.Shards),
			"DISK-INDICES": allocation.DiskIndices,
			"DISK-USED":    allocation.DiskUsed,
			"DISK-AVAIL":   allocation.DiskAvail,
			"DISK-TOTAL":   allocation.DiskTotal,
			"DISK-PERCENT": fmt.Sprintf("%d%%", allocation.DiskPercent),
			"HOST":         allocation.Host,
			"IP":           allocation.IP,
			"NODE":         allocation.Node,
		}

		row := make([]string, len(columnDefs))
		for i, colDef := range columnDefs {
			row[i] = rowData[colDef.Header]
		}
		data = append(data, row)
	}

	if len(flagSortBy) > 0 {
		sortCols := output.ParseSortColumns(flagSortBy)
		return output.PrintTable(columnDefs, data, sortCols)
	}
	sortCols := output.ParseSortColumns("SHARDS")
	return output.PrintTable(columnDefs, data, sortCols)
}
