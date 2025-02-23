package get

import (
	"fmt"
	"time"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getAllocationCmd = &cobra.Command{
	Use:   "allocation",
	Short: "Get Elasticsearch allocation",
	Long: utils.Trim(`
	Get Elasticsearch allocation. You can filter the results using the node flag.
	`),
	Example: utils.TrimAndIndent(`
	# Retrieve all allocation.
	esctl get allocation

	# Retrieve allocation for a specific node.
	esctl get allocation --node my_node
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		config := config.ParseConfigFile()

		// If --watch is NOT set, just run once
		if !flagRefresh {
			return handleAllocationLogic(*config)
		}

		// If --watch is set, run in a loop
		for {
			clearScreen() // optional, to mimic "watch" clearing
			err := handleAllocationLogic(*config)
			if err != nil {
				return err
			}
			time.Sleep(flagRefreshInterval)
		}
	},
}

func init() {
	getAllocationCmd.Flags().StringVar(&flagNode, "node", "", "Name of the node")
}

var allocationColumns = []output.ColumnDefaults{
	{Header: "SHARDS", Type: output.Number},
	{Header: "DISK-INDICES", Type: output.Text},
	{Header: "DISK-USED", Type: output.Text},
	{Header: "DISK-AVAIL", Type: output.Text},
	{Header: "DICK-TOTAL", Type: output.Number},
	{Header: "DISK-PERCENT", Type: output.Number},
	{Header: "HOST", Type: output.Number},
	{Header: "IP", Type: output.Date},
	{Header: "NODE", Type: output.DataSize},
}

func handleAllocationLogic(conf config.Config) error {
	allocations, err := es.GetAllocation(flagIndex)
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
			"SHARDS":       allocation.Shards,
			"DISK-INDICES": allocation.DiskIndices,
			"DISK-USED":    allocation.DiskUsed,
			"DISK-AVAIL":   allocation.DiskAvail,
			"DICK-TOTAL":   allocation.DiskTotal,
			"DISK-PERCENT": allocation.DiskPercent,
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
		output.PrintTable(columnDefs, data, sortCols)
	} else {
		sortCols := output.ParseSortColumns("SHARDS")
		output.PrintTable(columnDefs, data, sortCols)
	}

	return nil
}
