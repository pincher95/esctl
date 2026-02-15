package template

import (
	"fmt"
	"strings"

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

		# List templates by name substring
		esctl template list --name logs

		# Output as JSON
		esctl template list -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		templates, err := template.List(ctx)
		if err != nil {
			return err
		}

		// Also fetch legacy templates and merge (legacy templates won't overwrite composable ones)
		legacyTemplates, err := template.ListLegacy(ctx)
		if err == nil {
			for name, tmpl := range legacyTemplates {
				if _, exists := templates[name]; !exists {
					templates[name] = tmpl
				}
			}
		}

		if templateListName != "" {
			filtered := make(template.ListResponse)
			for name, tmpl := range templates {
				if strings.Contains(name, templateListName) {
					filtered[name] = tmpl
				}
			}
			if len(filtered) == 0 {
				return fmt.Errorf("no templates matched: %s", templateListName)
			}
			templates = filtered
		}

		if len(templates) == 0 {
			fmt.Println("No templates found")
			return nil
		}

		columnDefs := []output.ColumnDefaults{
			{Header: "name", Type: output.Text},
			{Header: "index_patterns", Type: output.Text},
			{Header: "priority", Type: output.Number},
			{Header: "version", Type: output.Number},
			{Header: "composed_of", Type: output.Text},
		}

		rows := make([]map[string]any, 0, len(templates))
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

			rows = append(rows, map[string]any{
				"name":           name,
				"index_patterns": tmpl.IndexPatterns,
				"priority":       tmpl.Priority,
				"version":        tmpl.Version,
				"composed_of":    tmpl.ComposedOf,
			})

			data = append(data, []string{
				name,
				patterns,
				fmt.Sprintf("%d", tmpl.Priority),
				fmt.Sprintf("%d", tmpl.Version),
				composedOf,
			})
		}

		if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
			return output.Render(rows)
		}

		return output.PrintTable(columnDefs, data, nil)
	},
}

var templateListName string

func init() {
	listCmd.Flags().StringVar(&templateListName, "name", "", "Filter templates by name or substring of template name")
}
