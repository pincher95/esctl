package component

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/template"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all component templates",
	Long: utils.Trim(`
		Lists all component templates in the Elasticsearch cluster. Component templates are
		building blocks that can be composed into index templates.
	`),
	Example: utils.TrimAndIndent(`
		# List all component templates
		esctl template component list

		# List with JSON output
		esctl template component list -o json

		# List with YAML output
		esctl template component list -o yaml
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleListComponents(ctx)
	},
}

func handleListComponents(ctx context.Context) error {
	templates, err := template.ListComponents(ctx)
	if err != nil {
		return err
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(templates)
	}

	columnDefs := []output.ColumnDefaults{
		{Header: "NAME", Type: output.Text},
		{Header: "VERSION", Type: output.Number},
	}

	data := make([][]string, 0, len(templates))
	for name, t := range templates {
		data = append(data, []string{
			name,
			fmt.Sprintf("%d", t.Version),
		})
	}

	return output.PrintTable(columnDefs, data, nil)
}
