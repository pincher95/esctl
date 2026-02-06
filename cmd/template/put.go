package template

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/template"
	"github.com/pincher95/esctl/internal/logger"
	"github.com/pincher95/esctl/internal/validation"
	"github.com/spf13/cobra"
)

var (
	flagFile     string
	flagPatterns []string
	flagPriority int
)

var putCmd = &cobra.Command{
	Use:   "put <name>",
	Short: "Create or update an index template",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
		# Create template from file
		esctl template put logs-template --file template.json

		# Create simple template with patterns
		esctl template put logs-template --patterns "logs-*" --patterns "app-logs-*"
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		name := args[0]

		// Validate template name
		if err := validation.ValidateTemplateName(name); err != nil {
			return fmt.Errorf("invalid template name: %w", err)
		}

		logger.Debug("creating/updating template", "name", name)

		var tmpl template.Template

		if flagFile != "" {
			// Read from file
			logger.Debug("reading template from file", "file", flagFile)
			data, err := os.ReadFile(flagFile)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}

			if err := json.Unmarshal(data, &tmpl); err != nil {
				return fmt.Errorf("failed to parse template JSON: %w", err)
			}
		} else if len(flagPatterns) > 0 {
			// Create from flags
			logger.Debug("creating template from flags",
				"patterns", flagPatterns,
				"priority", flagPriority)

			// Validate patterns
			for _, pattern := range flagPatterns {
				if err := validation.ValidateIndexPattern(pattern); err != nil {
					return fmt.Errorf("invalid index pattern %q: %w", pattern, err)
				}
			}

			// Validate priority
			if err := validation.ValidatePriority(flagPriority); err != nil {
				return fmt.Errorf("invalid priority: %w", err)
			}

			tmpl.IndexPatterns = flagPatterns
			tmpl.Priority = flagPriority
		} else {
			return fmt.Errorf("either --file or --patterns must be specified")
		}

		if err := template.Put(ctx, name, tmpl); err != nil {
			logger.Error("failed to put template", "name", name, "error", err)
			return err
		}

		logger.Info("template created/updated", "name", name)
		fmt.Printf("Template '%s' created/updated successfully\n", name)
		return nil
	},
}

func init() {
	putCmd.Flags().StringVar(&flagFile, "file", "", "JSON file containing template definition")
	putCmd.Flags().StringSliceVar(&flagPatterns, "patterns", []string{}, "Index patterns")
	putCmd.Flags().IntVar(&flagPriority, "priority", 0, "Template priority")
}
