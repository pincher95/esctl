package set

import (
	"fmt"
	"os"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/template"
	"github.com/pincher95/esctl/internal/validation"
	"github.com/spf13/cobra"
)

var (
	setTemplateName     string
	setTemplateFile     string
	setTemplatePatterns []string
	setTemplatePriority int

	setComponentTemplateName string
	setComponentTemplateFile string
)

var setTemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "Create or update an index template",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		# From a file
		esctl set template --name logs-template --file template.json

		# From flags
		esctl set template --name logs-template --patterns "logs-*" --patterns "app-*" --priority 100
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		if err := validation.ValidateTemplateName(setTemplateName); err != nil {
			return fmt.Errorf("invalid template name: %w", err)
		}

		var tmpl template.Template
		switch {
		case setTemplateFile != "":
			data, err := os.ReadFile(setTemplateFile)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			if err := utils.UnmarshalJSON(data, &tmpl, "failed to parse template JSON"); err != nil {
				return err
			}
		case len(setTemplatePatterns) > 0:
			for _, pattern := range setTemplatePatterns {
				if err := validation.ValidateIndexPattern(pattern); err != nil {
					return fmt.Errorf("invalid index pattern %q: %w", pattern, err)
				}
			}
			if err := validation.ValidatePriority(setTemplatePriority); err != nil {
				return fmt.Errorf("invalid priority: %w", err)
			}
			tmpl.IndexPatterns = setTemplatePatterns
			tmpl.Priority = setTemplatePriority
		default:
			return fmt.Errorf("either --file or --patterns must be specified")
		}

		if err := template.Put(ctx, setTemplateName, tmpl); err != nil {
			return err
		}
		fmt.Printf("Template '%s' created/updated successfully\n", setTemplateName)
		return nil
	},
}

var setComponentTemplateCmd = &cobra.Command{
	Use:   "component-template",
	Short: "Create or update a component template",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
		esctl set component-template --name settings-common --file component.json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		data, err := os.ReadFile(setComponentTemplateFile)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		var ct template.ComponentTemplate
		if err := utils.UnmarshalJSON(data, &ct, "failed to parse component template JSON"); err != nil {
			return err
		}

		if err := template.PutComponent(ctx, setComponentTemplateName, ct); err != nil {
			return err
		}
		fmt.Printf("Component template '%s' created/updated successfully\n", setComponentTemplateName)
		return nil
	},
}

func init() {
	setTemplateCmd.Flags().StringVar(&setTemplateName, "name", "", "Template name (required)")
	setTemplateCmd.Flags().StringVar(&setTemplateFile, "file", "", "JSON file containing the template definition")
	setTemplateCmd.Flags().StringSliceVar(&setTemplatePatterns, "patterns", nil, "Index patterns (alternative to --file)")
	setTemplateCmd.Flags().IntVar(&setTemplatePriority, "priority", 0, "Template priority (used with --patterns)")
	_ = setTemplateCmd.MarkFlagRequired("name")

	setComponentTemplateCmd.Flags().StringVar(&setComponentTemplateName, "name", "", "Component template name (required)")
	setComponentTemplateCmd.Flags().StringVar(&setComponentTemplateFile, "file", "", "JSON file containing the component template (required)")
	_ = setComponentTemplateCmd.MarkFlagRequired("name")
	_ = setComponentTemplateCmd.MarkFlagRequired("file")
}
