package delete

import (
	"github.com/pincher95/esctl/cmd/security"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var deleteUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Delete a user",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Delete user
	esctl delete user --name john_doe
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return security.HandleUserDelete(cmd.Context(), deleteUserName)
	},
}

var deleteUserName string

func init() {
	deleteUserCmd.Flags().StringVar(&deleteUserName, "name", "", "User name")
	deleteUserCmd.MarkFlagRequired("name")
}
