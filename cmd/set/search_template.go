package set

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/searchtemplate"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	flagSetSearchTemplateID   string
	flagSetSearchTemplateFile string
)

var setSearchTemplateCmd = &cobra.Command{
	Use:   "search-template",
	Short: "Create or update a stored search template",
	Long: utils.Trim(`
		Creates or updates a stored search template in Elasticsearch. Search templates
		define parameterized search queries that can be reused.

		The template should be provided in a JSON or YAML file containing the query structure
		with parameter placeholders using Mustache syntax ({{parameter}}).
	`),
	Example: utils.TrimAndIndent(`
		# Create a search template from a JSON file
		esctl set search-template --id my-template --file template.json

		# Create a search template from a YAML file
		esctl set search-template --id my-template --file template.yaml

		# Example JSON file format:
		# {
		#   "query": {
		#     "match": {
		#       "{{my_field}}": "{{my_value}}"
		#     }
		#   },
		#   "from": "{{from}}",
		#   "size": "{{size}}"
		# }
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleSetSearchTemplateLogic(ctx, flagSetSearchTemplateID, flagSetSearchTemplateFile)
	},
}

func init() {
	setSearchTemplateCmd.Flags().StringVar(&flagSetSearchTemplateID, "id", "", "Search template ID")
	setSearchTemplateCmd.Flags().StringVar(&flagSetSearchTemplateFile, "file", "", "Path to JSON/YAML file containing template definition")
	_ = setSearchTemplateCmd.MarkFlagRequired("id")
	_ = setSearchTemplateCmd.MarkFlagRequired("file")
}

func handleSetSearchTemplateLogic(ctx context.Context, id, filePath string) error {
	// Load template from file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var templateDef map[string]any

	// Try JSON first
	if err := json.Unmarshal(data, &templateDef); err != nil {
		// Try YAML
		if yamlErr := yaml.Unmarshal(data, &templateDef); yamlErr != nil {
			return fmt.Errorf("failed to parse file as JSON or YAML: json=%v, yaml=%v", err, yamlErr)
		}
	}

	template := searchtemplate.SearchTemplate{
		Template: templateDef,
	}

	if err := searchtemplate.Put(ctx, id, template); err != nil {
		return fmt.Errorf("failed to create/update search template: %w", err)
	}

	fmt.Printf("Successfully created/updated search template: %s\n", id)
	return nil
}
