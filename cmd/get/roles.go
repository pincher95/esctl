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

	# List roles by name substring
	esctl get roles --name admin
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return security.HandleRoleList(cmd.Context(), getRolesName)
	},
}

var getRolesName string

func init() {
	getRolesCmd.Flags().StringVar(&getRolesName, "name", "", "Filter roles by name or substring of role name")
}
