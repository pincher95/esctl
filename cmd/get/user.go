package get

import (
	"github.com/pincher95/esctl/cmd/security"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Get a user",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Get a specific user
	esctl get user --name john_doe
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return security.HandleUserGet(cmd.Context(), getUserName)
	},
}

var getUserName string

func init() {
	getUserCmd.Flags().StringVar(&getUserName, "name", "", "User name")
	getUserCmd.MarkFlagRequired("name")
}
