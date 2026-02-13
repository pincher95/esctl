package set

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/script"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	flagSetScriptID     string
	flagSetScriptLang   string
	flagSetScriptSource string
	flagSetScriptFile   string
)

var setScriptCmd = &cobra.Command{
	Use:   "script",
	Short: "Create or update a stored script",
	Long: utils.Trim(`
		Creates or updates a stored script in Elasticsearch. Scripts can be written in
		Painless, Mustache, Expression, or other supported languages.

		You can provide the script inline using --source or from a file using --file.
		The file should be JSON or YAML format with the script definition.
	`),
	Example: utils.TrimAndIndent(`
		# Create a script with inline source
		esctl set script --id my-script --lang painless --source "Math.log(_score * 2) + params.multiplier"

		# Create a script from a JSON file
		esctl set script --id my-script --file script.json

		# Create a Mustache template script
		esctl set script --id my-template --lang mustache --source '{"query":{"match":{"{{field}}":"{{value}}"}}}'

		# Example JSON file format:
		# {
		#   "script": {
		#     "lang": "painless",
		#     "source": "Math.log(_score * 2) + params.multiplier"
		#   }
		# }
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleSetScriptLogic(ctx, flagSetScriptID, flagSetScriptLang, flagSetScriptSource, flagSetScriptFile)
	},
}

func init() {
	setScriptCmd.Flags().StringVar(&flagSetScriptID, "id", "", "Script ID")
	setScriptCmd.Flags().StringVar(&flagSetScriptLang, "lang", "", "Script language (painless, mustache, expression)")
	setScriptCmd.Flags().StringVar(&flagSetScriptSource, "source", "", "Script source code (inline)")
	setScriptCmd.Flags().StringVar(&flagSetScriptFile, "file", "", "Path to JSON/YAML file containing script definition")
	_ = setScriptCmd.MarkFlagRequired("id")
}

func handleSetScriptLogic(ctx context.Context, id, lang, source, filePath string) error {
	var scriptData script.Script

	if filePath != "" {
		// Load from file
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		// Try JSON first
		var fileContent struct {
			Script script.Script `json:"script" yaml:"script"`
		}
		if err := json.Unmarshal(data, &fileContent); err != nil {
			// Try YAML
			if yamlErr := yaml.Unmarshal(data, &fileContent); yamlErr != nil {
				return fmt.Errorf("failed to parse file as JSON or YAML: json=%v, yaml=%v", err, yamlErr)
			}
		}
		scriptData = fileContent.Script
	} else {
		// Use inline flags
		if lang == "" {
			return fmt.Errorf("--lang is required when not using --file")
		}
		if source == "" {
			return fmt.Errorf("--source is required when not using --file")
		}
		scriptData = script.Script{
			Lang:   lang,
			Source: source,
		}
	}

	if err := script.Put(ctx, id, scriptData); err != nil {
		return fmt.Errorf("failed to create/update script: %w", err)
	}

	fmt.Printf("Successfully created/updated script: %s\n", id)
	return nil
}
