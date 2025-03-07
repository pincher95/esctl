package count

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pincher95/esctl/es"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var countCmd = &cobra.Command{
	Use:   "count [--index index] [--group-by field]",
	Short: "Count documents in an index or in all indices matching a pattern",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		handleCount()
	},
}

func Cmd() *cobra.Command {
	return countCmd
}

func handleCount() {
	var counts map[string]es.GroupCount
	var err error

	counts, err = es.CountDocuments(flagIndex, flagTerm, flagExists, flagNested, flagGroupBy, flagSize, flagTimeout, flagRefresh)
	if err != nil {
		fmt.Printf("Failed to get document counts: %v\n", err)
		os.Exit(1)
	}

	columnDefs := []output.ColumnDefaults{
		{Header: "INDEX", Type: output.Text},
		{Header: strings.ToUpper(flagGroupBy), Type: output.Text},
		{Header: "COUNT", Type: output.Number},
	}

	data := [][]string{}

	for index, groupCount := range counts {
		for group, count := range groupCount {
			rowData := map[string]string{
				"INDEX":                      index,
				strings.ToUpper(flagGroupBy): group,
				"COUNT":                      strconv.Itoa(count),
			}

			row := make([]string, len(columnDefs))
			for i, colDef := range columnDefs {
				row[i] = rowData[colDef.Header]
			}
			data = append(data, row)
		}
	}

	if len(flagSortBy) > 0 {
		sortCols := output.ParseSortColumns(flagSortBy)
		output.PrintTable(columnDefs, data, sortCols)
	} else {
		sortCols := output.ParseSortColumns("INDEX")
		output.PrintTable(columnDefs, data, sortCols)
	}
}

func init() {
	countCmd.Flags().StringVarP(&flagIndex, "index", "i", "", "Filter by specific indices or patterns")
	countCmd.Flags().StringSliceVarP(&flagTerm, "term", "t", []string{}, "Term filters to apply")
	countCmd.Flags().StringSliceVarP(&flagExists, "exists", "e", []string{}, "Exists filters to apply")
	countCmd.Flags().StringArrayVar(&flagNested, "nested", []string{}, "Nested paths")
	countCmd.Flags().StringVarP(&flagGroupBy, "group-by", "g", "", "Field to group the documents by")
	countCmd.Flags().StringVarP(&flagSortBy, "sort-by", "s", "", "Columns to sort by (comma-separated)")
	countCmd.Flags().IntVar(&flagSize, "size", 0, "Set max results per group")
	countCmd.Flags().StringVar(&flagTimeout, "timeout", "", "Set timeout for group by query")
	countCmd.Flags().BoolVar(&flagRefresh, "refresh", false, "Refresh index before counting documents")
}
