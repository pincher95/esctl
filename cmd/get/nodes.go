package get

import (
	"fmt"
	"os"
	"time"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/es"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getNodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Get all nodes in the Elasticsearch cluster",
	Run: func(cmd *cobra.Command, args []string) {
		config := config.ParseConfigFile()

		// If --watch is NOT set, just run once
		if !flagRefresh {
			handleNodeLogic(*config)
			return
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
	getNodesCmd.Flags().StringVar(&flagNode, "node", "", "Filter shards by node")
}

var nodeColumns = []output.ColumnDef{
	{Header: "NAME", Type: output.Text},
	{Header: "IP", Type: output.Text},
	{Header: "NODE-ROLE", Type: output.Text},
	{Header: "MASTER", Type: output.Text},
	{Header: "HEAP-MAX", Type: output.DataSize},
	{Header: "HEAP-CURRENT", Type: output.DataSize},
	{Header: "HEAP-PERCENT", Type: output.Percent},
	{Header: "RAM-MAX", Type: output.DataSize},
	{Header: "RAM-CURRENT", Type: output.DataSize},
	{Header: "RAM-PERCENT", Type: output.Percent},
	{Header: "CPU", Type: output.Percent},
	{Header: "LOAD-1M", Type: output.Number},
	{Header: "DISK-TOTAL", Type: output.DataSize},
	{Header: "DISK-USED", Type: output.DataSize},
	{Header: "DISK-AVAILABLE", Type: output.DataSize},
	{Header: "UPTIME", Type: output.Text},
}

func handleNodeLogic(conf config.Config) {
	nodes, err := es.GetNodes(flagNode)
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
			"NAME":           node.Name,
			"IP":             node.IP,
			"NODE-ROLE":      node.NodeRole,
			"MASTER":         node.Master,
			"HEAP-MAX":       node.HeapMax,
			"HEAP-CURRENT":   node.HeapCurrent,
			"HEAP-PERCENT":   node.HeapPercent + "%",
			"RAM-MAX":        node.RAMMax,
			"RAM-CURRENT":    node.RAMCurrent,
			"RAM-PERCENT":    node.RAMPercent + "%",
			"CPU":            node.CPU + "%",
			"LOAD-1M":        node.Load1m,
			"DISK-TOTAL":     node.DiskTotal,
			"DISK-USED":      node.DiskUsed,
			"DISK-AVAILABLE": node.DiskAvail,
			"UPTIME":         node.Uptime,
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
