package get

import (
	"github.com/pincher95/esctl/cmd/security"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var getUserCmd = &cobra.Command{
	Use:   "user <username>",
	Short: "Get a user",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Get a specific user
	esctl get user john_doe
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return security.HandleUserGet(cmd.Context(), args[0])
	},
}
