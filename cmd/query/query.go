package query

import (
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagId     []string
	flagTerm   []string
	flagNested []string
	flagSort   []string
	flagFrom   int
	flagSize   int
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query Elasticsearch",
	Long:  `This command allows you to query Elasticsearch.`,
	Example: utils.TrimAndIndent(`
esctl query articles
esctl query articles --id 61
esctl query articles --term "price:10" --size 1
esctl query articles --sort "price:desc" --from 10 --size 10`),
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		index := args[0]

		response, err := es.SearchDocuments(ctx, index, flagId, flagTerm, flagFrom, flagSize, flagNested, flagSort)
		if err != nil {
			return err
		}
		output.PrintJson(response["hits"])
		return nil
	},
}

func Cmd() *cobra.Command {
	return queryCmd
}

func init() {
	queryCmd.Flags().StringArrayVar(&flagId, "id", []string{}, "Document IDs to fetch")
	queryCmd.Flags().StringArrayVarP(&flagTerm, "term", "t", []string{}, "Term filter(s)")
	queryCmd.Flags().StringArrayVar(&flagNested, "nested", []string{}, "Nested path(s)")
	queryCmd.Flags().StringArrayVarP(&flagSort, "sort", "s", []string{}, "Sort definition(s)")
	queryCmd.Flags().IntVar(&flagFrom, "from", 0, "Starting document offset")
	queryCmd.Flags().IntVar(&flagSize, "size", 1, "Number of hits to return")
}
