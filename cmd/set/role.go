package set

import (
	"github.com/pincher95/esctl/cmd/security"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/spf13/cobra"
)

var setRoleFile string

var setRoleCmd = &cobra.Command{
	Use:   "role",
	Short: "Create or update a role",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Create role from file
	esctl set role --name read_only --file=role.json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return security.HandleRoleCreate(cmd.Context(), setRoleName, setRoleFile)
	},
}

func init() {
	setRoleCmd.Flags().StringVar(&setRoleName, "name", "", "Role name")
	setRoleCmd.Flags().StringVar(&setRoleFile, "file", "", "JSON file containing role definition")
	setRoleCmd.MarkFlagRequired("name")
	setRoleCmd.MarkFlagRequired("file")
}

var setRoleName string
