package set

import (
	"github.com/pincher95/esctl/cmd/security"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var setUserFile string

var setUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Create or update a user",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Create user from file
	esctl set user --name john_doe --file=user.json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return security.HandleUserCreate(cmd.Context(), setUserName, setUserFile)
	},
}

func init() {
	setUserCmd.Flags().StringVar(&setUserName, "name", "", "User name")
	setUserCmd.Flags().StringVar(&setUserFile, "file", "", "JSON file containing user definition")
	setUserCmd.MarkFlagRequired("name")
	setUserCmd.MarkFlagRequired("file")
}

var setUserName string
