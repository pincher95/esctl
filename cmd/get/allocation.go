package get

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"

	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/output"
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
		conf := config.ParseConfigFile()

		ctx := cmd.Context()

		if !flagRefresh {
			return handleAllocationLogic(ctx, allocationClient, *conf)
		}

		return utils.WatchLoop(flagRefreshInterval, func() error {
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
	{Header: "HOST", Type: output.Number},
	{Header: "IP", Type: output.Date},
	{Header: "NODE", Type: output.DataSize},
}

func handleAllocationLogic(ctx context.Context, client cat.Cat, conf config.Config) error {
	allocations, err := client.CatAllocation(ctx, "", flagNodeID, flagBytes)
	if err != nil {
		return fmt.Errorf("Failed to retrieve allocation: %v", err)
	}

	columnDefs, err := getColumnDefs(conf, "shards", allocationColumns)
	if err != nil {
		return fmt.Errorf("Failed to get column definitions: %v", err)
	}

	data := [][]string{}

	for _, allocation := range allocations {
		rowData := map[string]string{
			"SHARDS":       strconv.Itoa(allocation.Shards),
			"DISK-INDICES": utils.SafeString(allocation.DiskIndices),
			"DISK-USED":    utils.SafeString(allocation.DiskUsed),
			"DISK-AVAIL":   utils.SafeString(allocation.DiskAvail),
			"DISK-TOTAL":   utils.SafeString(allocation.DiskTotal),
			"DISK-PERCENT": fmt.Sprintf("%d%%", utils.SafeInt(allocation.DiskPercent)),
			"HOST":         utils.SafeString(allocation.Host),
			"IP":           utils.SafeString(allocation.IP),
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
		output.PrintTable(columnDefs, data, sortCols)
	} else {
		sortCols := output.ParseSortColumns("SHARDS")
		output.PrintTable(columnDefs, data, sortCols)
	}

	return nil
}
