package template

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/template"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
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

		if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
			return output.Render(tmpl)
		}

		// Build table with template settings as rows
		columnDefs := []output.ColumnDefaults{
			{Header: "setting", Type: output.Text},
			{Header: "value", Type: output.Text},
		}

		var data [][]string

		// Template metadata
		data = append(data, []string{"name", tmpl.Name})
		if len(tmpl.IndexPatterns) > 0 {
			data = append(data, []string{"index_patterns", strings.Join(tmpl.IndexPatterns, ", ")})
		}
		data = append(data, []string{"priority", fmt.Sprintf("%d", tmpl.Priority)})
		data = append(data, []string{"version", fmt.Sprintf("%d", tmpl.Version)})
		if len(tmpl.ComposedOf) > 0 {
			data = append(data, []string{"composed_of", strings.Join(tmpl.ComposedOf, ", ")})
		}

		// Flattened settings
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
	},
}

var templateName string

func init() {
	getCmd.Flags().StringVar(&templateName, "name", "", "Template name")
	getCmd.MarkFlagRequired("name")
}
