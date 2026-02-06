package get

import (
	"github.com/pincher95/esctl/cmd/security"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "List users",
	Example: utils.TrimAndIndent(`
	# List all users
	esctl get users

	# List users by name substring
	esctl get users --name admin
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return security.HandleUserList(cmd.Context(), getUsersName)
	},
}

var getUsersName string

func init() {
	getUsersCmd.Flags().StringVar(&getUsersName, "name", "", "Filter users by name or substring of user name")
}
