package get

import (
	"github.com/pincher95/esctl/cmd/security"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getRoleCmd = &cobra.Command{
	Use:   "role <role>",
	Short: "Get a role",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Get a specific role
	esctl get role read_only
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return security.HandleRoleGet(cmd.Context(), args[0])
	},
}
