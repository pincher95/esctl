package update

import (
	"github.com/pincher95/esctl/cmd/alias"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	updateAliasFrom string
	updateAliasTo   string
)

var updateAliasCmd = &cobra.Command{
	Use:   "alias <alias> --from=<from-index> --to=<to-index>",
	Short: "Move an alias from one index to another",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Move alias from old index to new index
	esctl update alias logs-current --from=logs-2023 --to=logs-2024
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return alias.HandleAliasMove(cmd.Context(), args[0], updateAliasFrom, updateAliasTo)
	},
}

func init() {
	updateAliasCmd.Flags().StringVar(&updateAliasFrom, "from", "", "Index to move the alias from")
	updateAliasCmd.Flags().StringVar(&updateAliasTo, "to", "", "Index to move the alias to")
	updateAliasCmd.MarkFlagRequired("from")
	updateAliasCmd.MarkFlagRequired("to")
}
