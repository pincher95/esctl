package component

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

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(tmpl)
	}

	columnDefs := []output.ColumnDefaults{
		{Header: "SETTING", Type: output.Text},
		{Header: "VALUE", Type: output.Text},
	}

	data := [][]string{
		{"name", name},
		{"version", fmt.Sprintf("%d", tmpl.Version)},
	}

	// Flatten and display settings
	if len(tmpl.Template.Settings) > 0 {
		flat := utils.FlattenSettingsMap(tmpl.Template.Settings)
		keys := make([]string, 0, len(flat))
		for k := range flat {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			data = append(data, []string{"settings." + k, fmt.Sprintf("%v", flat[k])})
		}
	}

	// Display aliases
	if len(tmpl.Template.Aliases) > 0 {
		aliases := make([]string, 0, len(tmpl.Template.Aliases))
		for a := range tmpl.Template.Aliases {
			aliases = append(aliases, a)
		}
		sort.Strings(aliases)
		data = append(data, []string{"aliases", strings.Join(aliases, ", ")})
	}

	return output.PrintTable(columnDefs, data, nil)
}
