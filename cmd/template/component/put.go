package component

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/template"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	flagPutName string
	flagPutFile string
)

var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Create or update a component template",
	Long: utils.Trim(`
		Creates or updates a component template. The template definition must be provided
		via a JSON or YAML file.
	`),
	Example: utils.TrimAndIndent(`
		# Create a component template from JSON file
		esctl template component put --name my-component --file template.json

		# Create a component template from YAML file
		esctl template component put --name my-component --file template.yaml

		# Update an existing component template
		esctl template component put --name my-component --file updated-template.json
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handlePutComponent(ctx, flagPutName, flagPutFile)
	},
}

func init() {
	putCmd.Flags().StringVar(&flagPutName, "name", "", "Component template name")
	putCmd.Flags().StringVar(&flagPutFile, "file", "", "Path to template definition file (JSON or YAML)")
	_ = putCmd.MarkFlagRequired("name")
	_ = putCmd.MarkFlagRequired("file")
}

func handlePutComponent(ctx context.Context, name, filePath string) error {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var tmpl template.ComponentTemplate

	// Try JSON first
	if err := json.Unmarshal(data, &tmpl); err != nil {
		// Try YAML
		if err := yaml.Unmarshal(data, &tmpl); err != nil {
			return fmt.Errorf("failed to parse template file as JSON or YAML: %w", err)
		}
	}

	// Create/update the template
	if err := template.PutComponent(ctx, name, tmpl); err != nil {
		return err
	}

	fmt.Printf("Successfully created/updated component template: %s\n", name)
	return nil
}
