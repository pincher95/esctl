package get

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/template"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

var (
	flagTemplateName          string
	flagComponentTemplateName string
)

var getTemplatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "List index templates or get a specific template",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		# List all index templates (composable + legacy)
		esctl get templates

		# Get a specific template
		esctl get templates --name logs-template -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if flagTemplateName != "" {
			return handleGetTemplate(ctx, flagTemplateName)
		}
		return handleListTemplates(ctx)
	},
}

var getComponentTemplatesCmd = &cobra.Command{
	Use:   "component-templates",
	Short: "List component templates or get a specific one",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		# List all component templates
		esctl get component-templates

		# Get a specific component template
		esctl get component-templates --name settings-common -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if flagComponentTemplateName != "" {
			ct, err := template.GetComponent(ctx, flagComponentTemplateName)
			if err != nil {
				return err
			}
			return output.Render(ct)
		}
		return handleListComponentTemplates(ctx)
	},
}

func init() {
	getTemplatesCmd.Flags().StringVar(&flagTemplateName, "name", "", "Template name (omit to list all)")
	getComponentTemplatesCmd.Flags().StringVar(&flagComponentTemplateName, "name", "", "Component template name (omit to list all)")
	getCmd.AddCommand(getTemplatesCmd)
	getCmd.AddCommand(getComponentTemplatesCmd)
}

func handleListTemplates(ctx context.Context) error {
	templates, err := template.List(ctx)
	if err != nil {
		return err
	}
	// Merge legacy templates without overwriting composable ones.
	if legacy, err := template.ListLegacy(ctx); err == nil {
		for name, tmpl := range legacy {
			if _, exists := templates[name]; !exists {
				templates[name] = tmpl
			}
		}
	}

	if len(templates) == 0 {
		fmt.Println("No templates found")
		return nil
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(templates)
	}

	columnDefs := []output.ColumnDefaults{
		{Header: "NAME", Type: output.Text},
		{Header: "INDEX-PATTERNS", Type: output.Text},
		{Header: "PRIORITY", Type: output.Number},
		{Header: "VERSION", Type: output.Number},
		{Header: "COMPOSED-OF", Type: output.Text},
	}

	data := make([][]string, 0, len(templates))
	for name, tmpl := range templates {
		data = append(data, []string{
			name,
			summarizeList(tmpl.IndexPatterns),
			fmt.Sprintf("%d", tmpl.Priority),
			fmt.Sprintf("%d", tmpl.Version),
			summarizeList(tmpl.ComposedOf),
		})
	}

	sortBy := flagSortBy
	if sortBy == "" {
		sortBy = "NAME"
	}
	return output.PrintTable(columnDefs, data, output.ParseSortColumns(sortBy))
}

func handleGetTemplate(ctx context.Context, name string) error {
	tmpl, err := template.Get(ctx, name)
	if err != nil {
		return err
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(tmpl)
	}

	columnDefs := []output.ColumnDefaults{
		{Header: "setting", Type: output.Text},
		{Header: "value", Type: output.Text},
	}
	data := [][]string{{"name", tmpl.Name}}
	if len(tmpl.IndexPatterns) > 0 {
		data = append(data, []string{"index_patterns", strings.Join(tmpl.IndexPatterns, ", ")})
	}
	data = append(data, []string{"priority", fmt.Sprintf("%d", tmpl.Priority)})
	data = append(data, []string{"version", fmt.Sprintf("%d", tmpl.Version)})
	if len(tmpl.ComposedOf) > 0 {
		data = append(data, []string{"composed_of", strings.Join(tmpl.ComposedOf, ", ")})
	}

	flat := utils.FlattenSettingsMap(tmpl.Template.Settings)
	keys := make([]string, 0, len(flat))
	for k := range flat {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		data = append(data, []string{k, utils.FormatSettingValue(flat[k])})
	}

	return output.PrintTable(columnDefs, data, nil)
}

func handleListComponentTemplates(ctx context.Context) error {
	templates, err := template.ListComponents(ctx)
	if err != nil {
		return err
	}
	if len(templates) == 0 {
		fmt.Println("No component templates found")
		return nil
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(templates)
	}

	columnDefs := []output.ColumnDefaults{
		{Header: "NAME", Type: output.Text},
		{Header: "VERSION", Type: output.Number},
	}
	data := make([][]string, 0, len(templates))
	for name, tmpl := range templates {
		data = append(data, []string{name, fmt.Sprintf("%d", tmpl.Version)})
	}

	sortBy := flagSortBy
	if sortBy == "" {
		sortBy = "NAME"
	}
	return output.PrintTable(columnDefs, data, output.ParseSortColumns(sortBy))
}

// summarizeList renders the first element with a "(+N more)" suffix.
func summarizeList(items []string) string {
	if len(items) == 0 {
		return ""
	}
	if len(items) == 1 {
		return items[0]
	}
	return fmt.Sprintf("%s (+%d more)", items[0], len(items)-1)
}
