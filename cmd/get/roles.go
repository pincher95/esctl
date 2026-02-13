package get

import (
	"github.com/pincher95/esctl/cmd/security"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getRolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "Get or list roles",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# List all roles
	esctl get roles

	# Get a specific role by name
	esctl get roles --name read_only

	# List roles by name substring
	esctl get roles --name admin
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// If a specific role name is provided, get that role
		if flagRolesName != "" {
			return security.HandleRoleGet(ctx, flagRolesName)
		}

		// Otherwise, list all roles (with optional filtering)
		return security.HandleRoleList(ctx, "")
	},
}

var flagRolesName string

func init() {
	getRolesCmd.Flags().StringVar(&flagRolesName, "name", "", "Role name to retrieve or substring for filtering")
}
