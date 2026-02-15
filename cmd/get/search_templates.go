package get

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/searchtemplate"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

var flagSearchTemplatesID string

var getSearchTemplatesCmd = &cobra.Command{
	Use:   "search-templates",
	Short: "List all stored search templates or get a specific search template",
	Long: utils.Trim(`
		Lists all stored search templates in Elasticsearch, or gets a specific search template by ID.
		Search templates allow you to define template queries with parameters that can be reused across searches.
	`),
	Example: utils.TrimAndIndent(`
		# List all search templates
		esctl get search-templates

		# Get a specific search template by ID
		esctl get search-templates --id my-template

		# List search templates in JSON format
		esctl get search-templates -o json

		# Get specific search template in YAML format
		esctl get search-templates --id my-template -o yaml
	`),
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// If ID provided, get specific template; otherwise list all
		if flagSearchTemplatesID != "" {
			return handleGetSearchTemplateLogic(ctx, flagSearchTemplatesID)
		}
		return handleGetSearchTemplatesLogic(ctx)
	},
}

func init() {
	getSearchTemplatesCmd.Flags().StringVar(&flagSearchTemplatesID, "id", "", "Search template ID")
}

func handleGetSearchTemplatesLogic(ctx context.Context) error {
	templates, err := searchtemplate.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list search templates: %w", err)
	}

	if len(templates) == 0 {
		fmt.Println("No stored search templates found")
		return nil
	}

	// Convert map to slice for consistent output
	type templateItem struct {
		ID     string `json:"id" yaml:"id"`
		Lang   string `json:"lang" yaml:"lang"`
		Source string `json:"source" yaml:"source"`
	}

	var items []templateItem
	for id, tmpl := range templates {
		lang, _ := tmpl.Template["lang"].(string)
		source, _ := tmpl.Template["source"].(string)
		items = append(items, templateItem{
			ID:     id,
			Lang:   lang,
			Source: source,
		})
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(items)
	}

	columnDefs := []output.ColumnDefaults{
		{Header: "ID", Type: output.Text},
		{Header: "LANG", Type: output.Text},
		{Header: "SOURCE", Type: output.Text},
	}

	data := make([][]string, 0, len(items))
	for _, item := range items {
		data = append(data, []string{item.ID, item.Lang, item.Source})
	}

	return output.PrintTable(columnDefs, data, nil)
}

func handleGetSearchTemplateLogic(ctx context.Context, id string) error {
	templateResp, err := searchtemplate.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get search template: %w", err)
	}

	// Create a displayable structure
	type templateDisplay struct {
		ID       string         `json:"id" yaml:"id"`
		Template map[string]any `json:"template" yaml:"template"`
	}

	display := templateDisplay{
		ID:       templateResp.ID,
		Template: templateResp.Template,
	}

	if shared.OutputFormat == "json" || shared.OutputFormat == "yaml" {
		return output.Render(display)
	}

	columnDefs := []output.ColumnDefaults{
		{Header: "setting", Type: output.Text},
		{Header: "value", Type: output.Text},
	}

	data := [][]string{
		{"id", display.ID},
	}
	for k, v := range display.Template {
		valJSON, _ := json.Marshal(v)
		data = append(data, []string{k, string(valJSON)})
	}

	return output.PrintTable(columnDefs, data, nil)
}
