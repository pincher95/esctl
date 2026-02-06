package get

import (
	"github.com/pincher95/esctl/cmd/security"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getRoleCmd = &cobra.Command{
	Use:   "role",
	Short: "Get a role",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Get a specific role
	esctl get role --name read_only
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return security.HandleRoleGet(cmd.Context(), getRoleName)
	},
}

var getRoleName string

func init() {
	getRoleCmd.Flags().StringVar(&getRoleName, "name", "", "Role name")
	getRoleCmd.MarkFlagRequired("name")
}
