package analyze

import (
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/analysis"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagIndex    string
	flagAnalyzer string
	flagField    string
	flagText     string
)

var Cmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze text with the specified analyzer or field",
	Example: utils.TrimAndIndent(`
    # Analyze text with standard analyzer
    esctl analyze --text="The quick brown fox"

    # Analyze using a field's analyzer
    esctl analyze --index=my-index --field=title --text="Elasticsearch"
    `),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		req := analysis.AnalyzeRequest{
			Analyzer: flagAnalyzer,
			Field:    flagField,
			Text:     []string{flagText},
		}
		resp, err := analysis.Analyze(ctx, flagIndex, req)
		if err != nil {
			return err
		}
		return output.Render(resp)
	},
}

func init() {
	Cmd.Flags().StringVar(&flagIndex, "index", "", "Index to use for analysis (optional)")
	Cmd.Flags().StringVar(&flagAnalyzer, "analyzer", "", "Analyzer to use")
	Cmd.Flags().StringVar(&flagField, "field", "", "Field to derive analyzer from")
	Cmd.Flags().StringVar(&flagText, "text", "", "Text to analyze")
	Cmd.MarkFlagRequired("text")
}
