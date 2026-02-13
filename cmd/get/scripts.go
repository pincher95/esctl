package get

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/script"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var flagScriptsID string

var getScriptsCmd = &cobra.Command{
	Use:   "scripts",
	Short: "List all stored scripts or get a specific script",
	Long: utils.Trim(`
		Lists all stored scripts in Elasticsearch, or gets a specific script by ID.
		Stored scripts can be used to centralize frequently used script logic in
		Painless, Mustache, or other supported languages.
	`),
	Example: utils.TrimAndIndent(`
		# List all stored scripts
		esctl get scripts

		# Get a specific script by ID
		esctl get scripts --id my-script

		# List scripts in JSON format
		esctl get scripts -o json

		# Get specific script in YAML format
		esctl get scripts --id my-script -o yaml
	`),
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// Determine script ID from flag
		scriptID := flagScriptsID

		// If ID provided, get specific script; otherwise list all
		if scriptID != "" {
			return handleGetScriptLogic(ctx, scriptID)
		}
		return handleGetScriptsLogic(ctx)
	},
}

func init() {
	getScriptsCmd.Flags().StringVar(&flagScriptsID, "id", "", "Script ID to retrieve")
}

func handleGetScriptsLogic(ctx context.Context) error {
	scripts, err := script.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list scripts: %w", err)
	}

	if len(scripts) == 0 {
		fmt.Println("No stored scripts found")
		return nil
	}

	// Convert map to slice for consistent output
	type scriptItem struct {
		ID     string `json:"id" yaml:"id"`
		Lang   string `json:"lang" yaml:"lang"`
		Source string `json:"source" yaml:"source"`
	}

	var items []scriptItem
	for id, sr := range scripts {
		items = append(items, scriptItem{
			ID:     id,
			Lang:   sr.Script.Lang,
			Source: sr.Script.Source,
		})
	}

	return output.Render(items)
}

func handleGetScriptLogic(ctx context.Context, id string) error {
	scriptResp, err := script.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get script: %w", err)
	}

	// Create a displayable structure
	type scriptDisplay struct {
		ID     string         `json:"id" yaml:"id"`
		Lang   string         `json:"lang" yaml:"lang"`
		Source string         `json:"source" yaml:"source"`
		Params map[string]any `json:"params,omitempty" yaml:"params,omitempty"`
	}

	display := scriptDisplay{
		ID:     scriptResp.ID,
		Lang:   scriptResp.Script.Lang,
		Source: scriptResp.Script.Source,
		Params: scriptResp.Script.Params,
	}

	return output.Render(display)
}
