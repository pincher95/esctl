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
	Use:   "alias",
	Short: "Move an alias from one index to another",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Move alias from old index to new index
	esctl update alias --name logs-current --from=logs-2023 --to=logs-2024
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return alias.HandleAliasMove(cmd.Context(), updateAliasName, updateAliasFrom, updateAliasTo)
	},
}

func init() {
	updateAliasCmd.Flags().StringVar(&updateAliasName, "name", "", "Alias name")
	updateAliasCmd.Flags().StringVar(&updateAliasFrom, "from", "", "Index to move the alias from")
	updateAliasCmd.Flags().StringVar(&updateAliasTo, "to", "", "Index to move the alias to")
	updateAliasCmd.MarkFlagRequired("name")
	updateAliasCmd.MarkFlagRequired("from")
	updateAliasCmd.MarkFlagRequired("to")
}

var updateAliasName string
