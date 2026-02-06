package set

import (
	"github.com/pincher95/esctl/cmd/security"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var setRoleFile string

var setRoleCmd = &cobra.Command{
	Use:   "role <role> --file=role.json",
	Short: "Create or update a role",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Create role from file
	esctl set role read_only --file=role.json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return security.HandleRoleCreate(cmd.Context(), args[0], setRoleFile)
	},
}

func init() {
	setRoleCmd.Flags().StringVar(&setRoleFile, "file", "", "JSON file containing role definition")
	setRoleCmd.MarkFlagRequired("file")
}
