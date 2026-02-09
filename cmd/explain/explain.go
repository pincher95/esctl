package explain

import (
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/analysis"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagIndex string
	flagDocID string
	flagQuery string
)

var Cmd = &cobra.Command{
	Use:   "explain",
	Short: "Explain why a document matches or doesn't match a query",
	Example: utils.TrimAndIndent(`
    esctl explain --index=my-index --id=1 --query='{"match":{"title":"fox"}}'
    `),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		q, err := utils.ParseJSONMap(flagQuery, "invalid query JSON")
		if err != nil {
			return err
		}
		req := analysis.ExplainRequest{Query: q}
		resp, err := analysis.Explain(ctx, flagIndex, flagDocID, req)
		if err != nil {
			return err
		}
		return output.Render(resp)
	},
}

func init() {
	Cmd.Flags().StringVar(&flagIndex, "index", "", "Index name")
	Cmd.Flags().StringVar(&flagDocID, "id", "", "Document ID")
	Cmd.Flags().StringVar(&flagQuery, "query", "", "Query JSON")
	Cmd.MarkFlagRequired("index")
	Cmd.MarkFlagRequired("id")
	Cmd.MarkFlagRequired("query")
}
