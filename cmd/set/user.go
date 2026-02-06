package set

import (
	"github.com/pincher95/esctl/cmd/security"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var setUserFile string

var setUserCmd = &cobra.Command{
	Use:   "user <username> --file=user.json",
	Short: "Create or update a user",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Create user from file
	esctl set user john_doe --file=user.json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return security.HandleUserCreate(cmd.Context(), args[0], setUserFile)
	},
}

func init() {
	setUserCmd.Flags().StringVar(&setUserFile, "file", "", "JSON file containing user definition")
	setUserCmd.MarkFlagRequired("file")
}
