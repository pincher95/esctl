package template

import (
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/template"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all index templates",
	Example: utils.TrimAndIndent(`
		# List all templates
		esctl template list

		# Output as JSON
		esctl template list -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		templates, err := template.List(ctx)
		if err != nil {
			return err
		}

		if len(templates) == 0 {
			fmt.Println("No templates found")
			return nil
		}

		if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
			return output.Render(templates)
		}

		// Prepare table data
		columnDefs := []output.ColumnDefaults{
			{Header: "NAME", Type: output.Text},
			{Header: "INDEX-PATTERNS", Type: output.Text},
			{Header: "PRIORITY", Type: output.Number},
			{Header: "VERSION", Type: output.Number},
			{Header: "COMPOSED-OF", Type: output.Text},
		}

		var data [][]string
		for name, tmpl := range templates {
			patterns := ""
			if len(tmpl.IndexPatterns) > 0 {
				patterns = tmpl.IndexPatterns[0]
				if len(tmpl.IndexPatterns) > 1 {
					patterns += fmt.Sprintf(" (+%d more)", len(tmpl.IndexPatterns)-1)
				}
			}

			composedOf := ""
			if len(tmpl.ComposedOf) > 0 {
				composedOf = tmpl.ComposedOf[0]
				if len(tmpl.ComposedOf) > 1 {
					composedOf += fmt.Sprintf(" (+%d more)", len(tmpl.ComposedOf)-1)
				}
			}

			data = append(data, []string{
				name,
				patterns,
				fmt.Sprintf("%d", tmpl.Priority),
				fmt.Sprintf("%d", tmpl.Version),
				composedOf,
			})
		}

		return output.PrintTable(columnDefs, data, nil)
	},
}
