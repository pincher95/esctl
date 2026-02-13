package component

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/template"
	"github.com/spf13/cobra"
)

var (
	flagDeleteName   string
	flagDeleteDryRun bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a component template",
	Long: utils.Trim(`
		Deletes a component template from the Elasticsearch cluster. This operation cannot be undone.
	`),
	Example: utils.TrimAndIndent(`
		# Delete a component template
		esctl template component delete --name my-component

		# Dry run - show what would be deleted
		esctl template component delete --name my-component --dry-run
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleDeleteComponent(ctx, flagDeleteName, flagDeleteDryRun)
	},
}

func init() {
	deleteCmd.Flags().StringVar(&flagDeleteName, "name", "", "Component template name")
	deleteCmd.Flags().BoolVar(&flagDeleteDryRun, "dry-run", false, "Show what would be deleted without actually deleting")
	_ = deleteCmd.MarkFlagRequired("name")
}

func handleDeleteComponent(ctx context.Context, name string, dryRun bool) error {
	if dryRun {
		fmt.Printf("DRY RUN: Would delete component template: %s\n", name)
		return nil
	}

	if err := template.DeleteComponent(ctx, name); err != nil {
		return err
	}

	fmt.Printf("Successfully deleted component template: %s\n", name)
	return nil
}
