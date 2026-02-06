package delete

import (
	"github.com/pincher95/esctl/cmd/security"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var deleteUserCmd = &cobra.Command{
	Use:   "user <username>",
	Short: "Delete a user",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Delete user
	esctl delete user john_doe
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return security.HandleUserDelete(cmd.Context(), args[0])
	},
}
