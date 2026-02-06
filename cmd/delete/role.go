package delete

import (
	"github.com/pincher95/esctl/cmd/security"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var deleteRoleCmd = &cobra.Command{
	Use:   "role <role>",
	Short: "Delete a role",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Delete role
	esctl delete role read_only
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return security.HandleRoleDelete(cmd.Context(), args[0])
	},
}
