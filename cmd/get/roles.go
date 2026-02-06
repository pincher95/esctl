package get

import (
	"github.com/pincher95/esctl/cmd/security"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getRolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "List roles",
	Example: utils.TrimAndIndent(`
	# List all roles
	esctl get roles
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return security.HandleRoleList(cmd.Context())
	},
}
