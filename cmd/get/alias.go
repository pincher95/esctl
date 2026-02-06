package get

import (
	"github.com/pincher95/esctl/cmd/alias"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getAliasCmd = &cobra.Command{
	Use:   "alias <alias>",
	Short: "Get details of a specific alias",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Get alias details
	esctl get alias my-alias
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return alias.HandleAliasGet(cmd.Context(), args[0])
	},
}
