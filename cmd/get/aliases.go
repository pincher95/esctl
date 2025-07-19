package get

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getAliasesCmd = &cobra.Command{
	Use:                   "aliases [--index index]",
	DisableFlagsInUseLine: true,
	Short:                 "Retrieves information for one or more data stream or index aliases.",
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
		aliasClient := index.NewIndex()
		conf := config.ParseConfigFile()

		ctx := cmd.Context()

		if !flagRefresh {
			handleAliasLogic(ctx, aliasClient, *conf)
			return
		}

		for {
			clearScreen()
			handleAliasLogic(ctx, aliasClient, *conf)
			time.Sleep(flagRefreshInterval)
		}
	},
}

func init() {
	getAliasesCmd.Flags().StringVarP(&flagIndex, "index", "i", "", "Name of the index")
	getAliasesCmd.Flags().BoolVar(&flagWritable, "writable", true, "Filter by writable index")
}

var aliasColumns = []output.ColumnDefaults{
	{Header: "ALIAS", Type: output.Text},
	{Header: "INDEX", Type: output.Text},
	{Header: "FILTER", Type: output.Text},
	{Header: "ROUTING", Type: output.Text},
	{Header: "INDEX-ROUTING", Type: output.Text},
	{Header: "SEARCH-ROUTING", Type: output.Text},
	{Header: "IS_WRITE_INDEX", Type: output.Boolean},
}

func handleAliasLogic(ctx context.Context, client index.Index, conf config.Config) {
	aliases, err := client.GetAliases(ctx, flagIndex)
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

	for idx, detail := range *aliases {
		for alias := range detail.Aliases {
			if includeIndexByWriteIndex(detail.Aliases[alias]) {
				rowData := map[string]string{
					"ALIAS":          alias,
					"INDEX":          idx,
					"INDEX-ROUTING":  detail.Aliases[alias].IndexRouting,
					"SEARCH-ROUTING": detail.Aliases[alias].SearchRouting,
					"ROUTING":        detail.Aliases[alias].Routing,
				}

				row := make([]string, len(columnDefs))
				for i, colDef := range columnDefs {
					row[i] = rowData[colDef.Header]
				}
				data = append(data, row)
			}
		}
	}

	sortCols := output.ParseSortColumns(flagSortBy)

	output.PrintTable(columnDefs, data, sortCols)
}

func includeIndexByWriteIndex(aliasInfo index.AliasInfo) bool {
	switch {
	case flagWritable && aliasInfo.IsWriteIndex:
		return true
	case !flagWritable:
		return true
	}
	return false
}
