package get

import (
	"fmt"
	"os"
	"time"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"

	cat "github.com/pincher95/esctl/es/cat"
)

var getNodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Get all nodes in the Elasticsearch cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		// ctx := cmd.Context()
		config := config.ParseConfigFile()

		// select
		// If --watch is NOT set, just run once
		if !flagRefresh {
			handleNodeLogic(*config)
			return nil
		}

		// If --watch is set, run in a loop
		for {
			clearScreen() // optional, to mimic "watch" clearing
			handleNodeLogic(*config)
			time.Sleep(flagRefreshInterval)
		}
	},
}

func init() {
	getNodesCmd.Flags().StringVar(&flagNode, "filter", "", "Filter by node name or substring of node name e.g. 'data-node-1', 'node', 'data'")
	getNodesCmd.Flags().StringVar(&flagBytes, "bytes", "", "The unit in which to display byte values. Valid values are: 'b', 'kb', 'mb', 'gb', 'tb', 'pb'.")
	getNodesCmd.Flags().StringVar(&flagTime, "time", "", "Specifies the time units, for example, 5d or 7h. Valid values are: nanos, micros, ms, s, m, h, d.")
}

var nodeColumns = []output.ColumnDefaults{
	{Header: "IP", Type: output.Text},
	{Header: "HEAP-PERCENT", Type: output.Percent},
	{Header: "RAM-PERCENT", Type: output.Percent},
	{Header: "CPU", Type: output.Percent},
	{Header: "LOAD-1M", Type: output.Number},
	{Header: "LOAD-5M", Type: output.Number},
	{Header: "LOAD-15M", Type: output.Number},
	{Header: "NODE-ROLE", Type: output.Text},
	{Header: "NODE-ROLES", Type: output.Text},
	{Header: "MASTER", Type: output.Text},
	{Header: "NAME", Type: output.Text},
}

func handleNodeLogic(conf config.Config) {
	nodes, err := cat.CatNodes(nil, &flagNode, &flagBytes, &flagTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to retrieve nodes: %v\n", err)
		os.Exit(1)
	}

	columnDefs, err := getColumnDefs(conf, "node", nodeColumns)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get column definitions:", err)
		os.Exit(1)
	}

	data := [][]string{}

	for _, node := range nodes {
		rowData := map[string]string{
			"IP":           node.IP,
			"HEAP-PERCENT": fmt.Sprintf("%d%%", utils.SafeInt(node.HeapPercent)),
			"RAM-PERCENT":  fmt.Sprintf("%d%%", utils.SafeInt(node.RAMPercent)),
			"CPU":          fmt.Sprintf("%d%%", utils.SafeInt(node.CPU)),
			"LOAD-1M":      utils.SafeString(node.Load1M),
			"LOAD-5M":      utils.SafeString(node.Load5M),
			"LOAD-15M":     utils.SafeString(node.Load15M),
			"NODE-ROLE":    node.Role,
			"NODE-ROLES":   node.Roles,
			"MASTER":       node.Master,
			"NAME":         node.Name,
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
		sortCols := output.ParseSortColumns("NAME")
		output.PrintTable(columnDefs, data, sortCols)
	}
}
