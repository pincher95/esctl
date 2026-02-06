package profile

import (
	"encoding/json"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/analysis"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagIndex string
	flagQuery string
)

var Cmd = &cobra.Command{
	Use:   "profile",
	Short: "Run a query with profiling enabled",
	Example: utils.TrimAndIndent(`
    esctl profile --index=logs --query='{"match_all":{}}'
    `),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		var q map[string]any
		if err := json.Unmarshal([]byte(flagQuery), &q); err != nil {
			return fmt.Errorf("invalid query JSON: %w", err)
		}
		resp, err := analysis.ProfileSearch(ctx, flagIndex, q)
		if err != nil {
			return err
		}
		return output.Render(resp)
	},
}

func init() {
	Cmd.Flags().StringVar(&flagIndex, "index", "", "Index name")
	Cmd.Flags().StringVar(&flagQuery, "query", "", "Query JSON")
	Cmd.MarkFlagRequired("index")
	Cmd.MarkFlagRequired("query")
}
