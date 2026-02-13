package component

import (
	"context"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/template"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagGetName string
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a specific component template",
	Long: utils.Trim(`
		Retrieves the details of a specific component template by name.
	`),
	Example: utils.TrimAndIndent(`
		# Get a component template
		esctl template component get --name my-component

		# Get with JSON output
		esctl template component get --name my-component -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleGetComponent(ctx, flagGetName)
	},
}

func init() {
	getCmd.Flags().StringVar(&flagGetName, "name", "", "Component template name")
	_ = getCmd.MarkFlagRequired("name")
}

func handleGetComponent(ctx context.Context, name string) error {
	tmpl, err := template.GetComponent(ctx, name)
	if err != nil {
		return err
	}

	return output.Render(tmpl)
}
