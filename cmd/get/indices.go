package get

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	cat "github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getIndicesCmd = &cobra.Command{
	Use:   "indices",
	Short: "Get Elasticsearch indices",
	Long: utils.Trim(`
	Get Elasticsearch indices. You can filter the results using the index flag.
	`),
	Example: utils.TrimAndIndent(`
	# Retrieve all indices.
	esctl get indices

	# Retrieve indices for a specific index.
	esctl get indices --index my_index
	`),
	Run: func(cmd *cobra.Command, args []string) {
		config := config.ParseConfigFile()

		// If --watch is NOT set, just run once
		if !flagRefresh {
			handleIndicesLogic(*config)
			return
		}

		// If --watch is set, run in a loop
		for {
			clearScreen() // optional, to mimic "watch" clearing
			handleIndicesLogic(*config)
			time.Sleep(flagRefreshInterval)
		}
	},
}

func init() {
	getIndicesCmd.Flags().StringVarP(&flagIndex, "index", "i", "", "Name of the index")
	getIndicesCmd.Flags().StringVar(&flagBytes, "bytes", "", "The unit in which to display byte values. Valid values are: 'b', 'kb', 'mb', 'gb', 'tb', 'pb'.")
}

var indexColumns = []output.ColumnDefaults{
	{Header: "HEALTH", Type: output.Text},
	{Header: "STATUS", Type: output.Text},
	{Header: "INDEX", Type: output.Text},
	{Header: "UUID", Type: output.Text},
	{Header: "PRIMARY", Type: output.Number},
	{Header: "REPLICAS", Type: output.Number},
	{Header: "DOCS-COUNT", Type: output.Number},
	{Header: "DOCS-DELETED", Type: output.Number},
	{Header: "STORE-SIZE", Type: output.DataSize},
	{Header: "PRI-STORE-SIZE", Type: output.DataSize},
}

func handleIndicesLogic(conf config.Config) {
	indices, err := cat.CatIndices(nil, &flagIndex, &flagBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve indices:", err)
		os.Exit(1)
	}

	columnDefs, err := getColumnDefs(conf, "index", indexColumns)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get column definitions:", err)
		os.Exit(1)
	}

	data := [][]string{}

	for _, index := range indices {
		rowData := map[string]string{
			"HEALTH":         index.Health,
			"STATUS":         index.Status,
			"INDEX":          index.Index,
			"UUID":           index.UUID,
			"PRIMARY":        strconv.Itoa(utils.SafeInt(index.Primary)),
			"REPLICAS":       strconv.Itoa(utils.SafeInt(index.Replica)),
			"DOCS-COUNT":     strconv.Itoa(utils.SafeInt(index.DocsCount)),
			"DOCS-DELETED":   strconv.Itoa(utils.SafeInt(index.DocDeleted)),
			"STORE-SIZE":     utils.SafeString(index.StoreSize),
			"PRI-STORE-SIZE": utils.SafeString(index.PrimaryStoreSize),
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
		sortCols := output.ParseSortColumns("INDEX")
		output.PrintTable(columnDefs, data, sortCols)
	}
}
