package template

import (
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/template"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a specific index template",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		# Get template details
		esctl template get --name logs-template

		# Get as JSON
		esctl template get --name logs-template -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		name := templateName

		tmpl, err := template.Get(ctx, name)
		if err != nil {
			return err
		}

		// Output in JSON/YAML for detailed view
		return output.Render(tmpl)
	},
}

var templateName string

func init() {
	getCmd.Flags().StringVar(&templateName, "name", "", "Template name")
	getCmd.MarkFlagRequired("name")
}
