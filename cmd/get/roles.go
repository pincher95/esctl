package get

import (
	"github.com/pincher95/esctl/cmd/security"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getRolesCmd = &cobra.Command{
	Use:   "roles [NAME]",
	Short: "Get or list roles",
	Args:  cobra.MaximumNArgs(1),
	Example: utils.TrimAndIndent(`
	# List all roles
	esctl get roles

	# Get a specific role by name (positional argument)
	esctl get roles read_only

	# Get a specific role by name (flag)
	esctl get roles --name read_only

	# List roles by name substring
	esctl get roles --name admin
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// Determine if we're getting a specific role or listing all
		var roleName string
		if len(args) > 0 {
			roleName = args[0]
		} else if flagRolesName != "" {
			roleName = flagRolesName
		}

		// If a specific role name is provided, get that role
		if roleName != "" {
			return security.HandleRoleGet(ctx, roleName)
		}

		// Otherwise, list all roles (with optional filtering)
		return security.HandleRoleList(ctx, "")
	},
}

var flagRolesName string

func init() {
	getRolesCmd.Flags().StringVar(&flagRolesName, "name", "", "Role name (for getting specific role) or substring (for filtering list)")
}
