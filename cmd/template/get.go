package template

import (
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/template"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <name>",
	Short: "Get details of a specific index template",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
		# Get template details
		esctl template get logs-template

		# Get as JSON
		esctl template get logs-template -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		name := args[0]

		tmpl, err := template.Get(ctx, name)
		if err != nil {
			return err
		}

		// Output in JSON/YAML for detailed view
		return output.Render(tmpl)
	},
}
