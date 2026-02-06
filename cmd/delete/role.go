package delete

import (
	"github.com/pincher95/esctl/cmd/security"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var deleteRoleCmd = &cobra.Command{
	Use:   "role",
	Short: "Delete a role",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Delete role
	esctl delete role --name read_only
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return security.HandleRoleDelete(cmd.Context(), deleteRoleName)
	},
}

var deleteRoleName string

func init() {
	deleteRoleCmd.Flags().StringVar(&deleteRoleName, "name", "", "Role name")
	deleteRoleCmd.MarkFlagRequired("name")
}
