package update

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/searchtemplate"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	flagRenderTemplateID     string
	flagRenderTemplateParams string
	flagRenderTemplateFile   string
)

var updateSearchTemplateCmd = &cobra.Command{
	Use:   "search-template",
	Short: "Render a search template with parameters",
	Long: utils.Trim(`
		Renders a stored search template with the provided parameters to see the final query.
		This is useful for testing and debugging search templates before using them.
	`),
}

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render a search template with parameters",
	Long: utils.Trim(`
		Renders a stored search template with the provided parameters to see the final query.
		Parameters can be provided as a JSON string or from a file.
	`),
	Example: utils.TrimAndIndent(`
		# Render a template with inline JSON parameters
		esctl update search-template render --id my-template --params '{"field":"title","value":"search"}'

		# Render a template with parameters from a file
		esctl update search-template render --id my-template --file params.json

		# Render and view output in JSON format
		esctl update search-template render --id my-template --params '{"query":"test"}' -o json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleRenderSearchTemplateLogic(ctx, flagRenderTemplateID, flagRenderTemplateParams, flagRenderTemplateFile)
	},
}

func init() {
	renderCmd.Flags().StringVar(&flagRenderTemplateID, "id", "", "Search template ID")
	renderCmd.Flags().StringVar(&flagRenderTemplateParams, "params", "", "Parameters as JSON string")
	renderCmd.Flags().StringVar(&flagRenderTemplateFile, "file", "", "Path to JSON/YAML file containing parameters")
	_ = renderCmd.MarkFlagRequired("id")

	updateSearchTemplateCmd.AddCommand(renderCmd)
}

func handleRenderSearchTemplateLogic(ctx context.Context, id, paramsJSON, filePath string) error {
	var params map[string]any

	if filePath != "" {
		// Load params from file
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		// Try JSON first
		if err := json.Unmarshal(data, &params); err != nil {
			// Try YAML
			if yamlErr := yaml.Unmarshal(data, &params); yamlErr != nil {
				return fmt.Errorf("failed to parse file as JSON or YAML: json=%v, yaml=%v", err, yamlErr)
			}
		}
	} else if paramsJSON != "" {
		// Parse inline JSON params
		if err := json.Unmarshal([]byte(paramsJSON), &params); err != nil {
			return fmt.Errorf("failed to parse params JSON: %w", err)
		}
	}

	result, err := searchtemplate.Render(ctx, id, params)
	if err != nil {
		return fmt.Errorf("failed to render search template: %w", err)
	}

	return output.Render(result.TemplateOutput)
}
