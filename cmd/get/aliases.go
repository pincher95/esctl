package get

import (
	"fmt"
	"os"
	"time"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getAliasesCmd = &cobra.Command{
	Use:   "aliases",
	Short: "Get Elasticsearch aliases",
	Long: utils.Trim(`
	Get Elasticsearch aliases. You can filter the results using the index flag.
	`),
	Example: utils.TrimAndIndent(`
	# Retrieve all aliases.
	esctl get aliases

	# Retrieve aliases for a specific index.
	esctl get aliases --index my_index
	`),
	Run: func(cmd *cobra.Command, args []string) {
		config := config.ParseConfigFile()

		// If --watch is NOT set, just run once
		if !flagRefresh {
			handleAliasLogic(*config)
			return
		}

		// If --watch is set, run in a loop
		for {
			clearScreen() // optional, to mimic "watch" clearing
			handleAliasLogic(*config)
			time.Sleep(flagRefreshInterval)
		}
	},
}

func init() {
	getAliasesCmd.Flags().StringVarP(&flagIndex, "index", "i", "", "Name of the index")
}

var aliasColumns = []output.ColumnDef{
	{Header: "ALIAS", Type: output.Text},
	{Header: "INDEX", Type: output.Text},
}

func handleAliasLogic(conf config.Config) {
	aliases, err := es.GetAliases(flagIndex)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve aliases:", err)
		os.Exit(1)
	}

	columnDefs, err := getColumnDefs(conf, "alias", aliasColumns)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get column definitions:", err)
		os.Exit(1)
	}

	data := [][]string{}

	for alias, index := range aliases {
		rowData := map[string]string{
			"ALIAS": alias,
			"INDEX": index,
		}

		row := make([]string, len(columnDefs))
		for i, colDef := range columnDefs {
			row[i] = rowData[colDef.Header]
		}
		data = append(data, row)
	}

	sortCols := output.ParseSortColumns(flagSortBy)

	output.PrintTable(columnDefs, data, sortCols)
}
