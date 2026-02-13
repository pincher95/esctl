package get

import (
	"github.com/pincher95/esctl/cmd/security"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getUsersCmd = &cobra.Command{
	Use:   "users [NAME]",
	Short: "Get or list users",
	Args:  cobra.MaximumNArgs(1),
	Example: utils.TrimAndIndent(`
	# List all users
	esctl get users

	# Get a specific user by name (positional argument)
	esctl get users john_doe

	# Get a specific user by name (flag)
	esctl get users --name john_doe

	# List users by name substring
	esctl get users --name admin
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// Determine if we're getting a specific user or listing all
		var userName string
		if len(args) > 0 {
			userName = args[0]
		} else if flagUsersName != "" {
			userName = flagUsersName
		}

		// If a specific user name is provided, get that user
		if userName != "" {
			return security.HandleUserGet(ctx, userName)
		}

		// Otherwise, list all users (with optional filtering)
		return security.HandleUserList(ctx, "")
	},
}

var flagUsersName string

func init() {
	getUsersCmd.Flags().StringVar(&flagUsersName, "name", "", "User name (for getting specific user) or substring (for filtering list)")
}
