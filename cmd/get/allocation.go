package get

import (
	"fmt"
	"os"

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
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.ParseConfigFile()
		handleAllocationLogic(*conf)
	},
}

func init() {
	getAllocationCmd.Flags().StringVar(&flagNode, "node", "", "Name of the node")
}

var allocationColumns = []output.ColumnDef{
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

func handleAllocationLogic(conf config.Config) {
	allocations, err := es.GetAllocation(flagIndex)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve allocation:", err)
		os.Exit(1)
	}

	columnDefs, err := getColumnDefs(conf, "shards", allocationColumns)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get column definitions:", err)
		os.Exit(1)
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
}
