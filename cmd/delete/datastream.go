package delete

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/datastream"
	"github.com/spf13/cobra"
)

var (
	flagDataStreamName   string
	flagDataStreamDryRun bool
)

var dataStreamCmd = &cobra.Command{
	Use:   "data-stream",
	Short: "Delete a data stream",
	Long: utils.Trim(`
		Deletes a data stream and all of its backing indices. This operation cannot be undone.
		Use with caution as all data in the backing indices will be permanently deleted.
	`),
	Example: utils.TrimAndIndent(`
		# Delete a data stream
		esctl delete data-stream --name logs-app

		# Dry run - show what would be deleted
		esctl delete data-stream --name logs-app --dry-run
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleDeleteDataStream(ctx, flagDataStreamName, flagDataStreamDryRun)
	},
}

func init() {
	dataStreamCmd.Flags().StringVar(&flagDataStreamName, "name", "", "Data stream name")
	dataStreamCmd.Flags().BoolVar(&flagDataStreamDryRun, "dry-run", false, "Show what would be deleted without actually deleting")
	_ = dataStreamCmd.MarkFlagRequired("name")
}

func handleDeleteDataStream(ctx context.Context, name string, dryRun bool) error {
	if dryRun {
		fmt.Printf("DRY RUN: Would delete data stream: %s\n", name)
		fmt.Println("WARNING: This would also delete all backing indices and their data!")
		return nil
	}

	// Get approval from user
	fmt.Printf("WARNING: Deleting data stream '%s' will also delete all its backing indices and data.\n", name)
	approved, err := utils.GetApproval()
	if err != nil {
		return err
	}

	if !approved {
		fmt.Println("Operation cancelled")
		return nil
	}

	if err := datastream.Delete(ctx, name); err != nil {
		return fmt.Errorf("failed to delete data stream: %w", err)
	}

	fmt.Printf("Successfully deleted data stream: %s\n", name)
	return nil
}
