package set

import (
	"github.com/pincher95/esctl/cmd/alias"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	setAliasIndices string
	setAliasFilter  string
	setAliasRouting string
)

var setAliasCmd = &cobra.Command{
	Use:   "alias <alias> --indices=<index1,index2,...>",
	Short: "Add an alias to one or more indices",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Add alias to a single index
	esctl set alias my-alias --indices=my-index

	# Add alias to multiple indices
	esctl set alias logs-current --indices="logs-2023,logs-2024"
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return alias.HandleAliasAdd(cmd.Context(), args[0], setAliasIndices, setAliasFilter, setAliasRouting)
	},
}

func init() {
	setAliasCmd.Flags().StringVar(&setAliasIndices, "indices", "", "Comma-separated list of indices to add the alias to")
	setAliasCmd.Flags().StringVar(&setAliasFilter, "filter", "", "Filter to apply to the alias (JSON)")
	setAliasCmd.Flags().StringVar(&setAliasRouting, "routing", "", "Routing value for the alias")
	setAliasCmd.MarkFlagRequired("indices")
}
