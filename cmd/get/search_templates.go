package get

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/searchtemplate"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var flagSearchTemplatesID string

var getSearchTemplatesCmd = &cobra.Command{
	Use:   "search-templates [ID]",
	Short: "List all stored search templates or get a specific search template",
	Long: utils.Trim(`
		Lists all stored search templates in Elasticsearch, or gets a specific search template by ID.
		Search templates allow you to define template queries with parameters that can be reused across searches.
	`),
	Example: utils.TrimAndIndent(`
		# List all search templates
		esctl get search-templates

		# Get a specific search template by ID (positional argument)
		esctl get search-templates my-template

		# Get a specific search template by ID (flag)
		esctl get search-templates --id my-template

		# List search templates in JSON format
		esctl get search-templates -o json

		# Get specific search template in YAML format
		esctl get search-templates my-template -o yaml
	`),
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// Determine template ID from positional arg or flag
		var id string
		if len(args) > 0 {
			id = args[0]
		} else if flagSearchTemplatesID != "" {
			id = flagSearchTemplatesID
		}

		// If ID provided, get specific template; otherwise list all
		if id != "" {
			return handleGetSearchTemplateLogic(ctx, id)
		}
		return handleGetSearchTemplatesLogic(ctx)
	},
}

func init() {
	getSearchTemplatesCmd.Flags().StringVar(&flagSearchTemplatesID, "id", "", "Search template ID (alternative to positional argument)")
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
		ID       string         `json:"id" yaml:"id"`
		Template map[string]any `json:"template" yaml:"template"`
	}

	var items []templateItem
	for id, tmpl := range templates {
		items = append(items, templateItem{
			ID:       id,
			Template: tmpl.Template,
		})
	}

	return output.Render(items)
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

	return output.Render(display)
}
