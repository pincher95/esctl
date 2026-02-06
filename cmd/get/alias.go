package get

import (
	"github.com/pincher95/esctl/cmd/alias"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getAliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Get details of a specific alias",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Get alias details
	esctl get alias --name my-alias
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return alias.HandleAliasGet(cmd.Context(), getAliasName)
	},
}

var getAliasName string

func init() {
	getAliasCmd.Flags().StringVar(&getAliasName, "name", "", "Alias name")
	getAliasCmd.MarkFlagRequired("name")
}
