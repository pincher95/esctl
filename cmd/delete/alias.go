package delete

import (
	"github.com/pincher95/esctl/cmd/alias"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var deleteAliasIndices string

var deleteAliasCmd = &cobra.Command{
	Use:   "alias <alias> --indices=<index1,index2,...>",
	Short: "Remove an alias from one or more indices",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Remove alias from a single index
	esctl delete alias my-alias --indices=my-index
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return alias.HandleAliasRemove(cmd.Context(), args[0], deleteAliasIndices)
	},
}

func init() {
	deleteAliasCmd.Flags().StringVar(&deleteAliasIndices, "indices", "", "Comma-separated list of indices to remove the alias from")
	deleteAliasCmd.MarkFlagRequired("indices")
}
