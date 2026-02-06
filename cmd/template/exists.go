package template

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/template"
	"github.com/spf13/cobra"
)

var existsCmd = &cobra.Command{
	Use:   "exists",
	Short: "Check if an index template exists",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		# Check if template exists
		esctl template exists --name logs-template
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		name := templateExistsName

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

var templateExistsName string

func init() {
	existsCmd.Flags().StringVar(&templateExistsName, "name", "", "Template name")
	existsCmd.MarkFlagRequired("name")
}
