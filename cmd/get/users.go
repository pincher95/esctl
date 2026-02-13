package get

import (
	"github.com/pincher95/esctl/cmd/security"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Get or list users",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# List all users
	esctl get users

	# Get a specific user by name
	esctl get users --name john_doe

	# List users by name substring
	esctl get users --name admin
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// If a specific user name is provided, get that user
		if flagUsersName != "" {
			return security.HandleUserGet(ctx, flagUsersName)
		}

		// Otherwise, list all users (with optional filtering)
		return security.HandleUserList(ctx, "")
	},
}

var flagUsersName string

func init() {
	getUsersCmd.Flags().StringVar(&flagUsersName, "name", "", "User name to retrieve or substring for filtering")
}
