package ilm

import (
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/ilm"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <policy>",
	Short: "Get details of a specific ILM policy",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
		# Get ILM policy details
		esctl ilm get hot_delete_policy

		# Get as JSON
		esctl ilm get hot_delete_policy -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		name := args[0]

		policy, err := ilm.Get(ctx, name)
		if err != nil {
			return err
		}

		return output.Render(policy)
	},
}
