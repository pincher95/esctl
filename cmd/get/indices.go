package get

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	cat "github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getIndicesCmd = &cobra.Command{
	Use:                   "indices [--index index] [--bytes unit]",
	DisableFlagsInUseLine: true,
	Short:                 "Get Elasticsearch indices",
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
		indicesClient := cat.NewCat()
		conf := config.ParseConfigFile()

		ctx := cmd.Context()

		if !flagRefresh {
			handleIndicesLogic(ctx, indicesClient, *conf)
			return
		}

		utils.WatchLoop(flagRefreshInterval, func() error {
			handleIndicesLogic(ctx, indicesClient, *conf)
			return nil
		})
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

func handleIndicesLogic(ctx context.Context, client cat.Cat, conf config.Config) {
	indices, err := client.CatIndices(ctx, "", flagIndex, flagBytes)
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
			"PRIMARY":        strconv.Itoa(index.Primary),
			"REPLICAS":       strconv.Itoa(index.Replica),
			"DOCS-COUNT":     strconv.Itoa(index.DocsCount),
			"DOCS-DELETED":   strconv.Itoa(index.DocDeleted),
			"STORE-SIZE":     index.StoreSize,
			"PRI-STORE-SIZE": index.PrimaryStoreSize,
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
