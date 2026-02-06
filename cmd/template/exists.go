package template

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/template"
	"github.com/spf13/cobra"
)

var existsCmd = &cobra.Command{
	Use:   "exists <name>",
	Short: "Check if an index template exists",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
		# Check if template exists
		esctl template exists logs-template
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		name := args[0]

		exists, err := template.Exists(ctx, name)
		if err != nil {
			return err
		}

		if exists {
			fmt.Printf("Template '%s' exists\n", name)
		} else {
			fmt.Printf("Template '%s' does not exist\n", name)
		}

		return nil
	},
}
