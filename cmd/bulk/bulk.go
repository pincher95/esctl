package bulk

import (
	"context"
	"fmt"
	"os"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/bulk"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

var bulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Execute bulk operations",
	Long:  "Execute bulk index, update, and delete operations from NDJSON files or stdin",
}

var (
	fromFileIndex   string
	fromFileRefresh string
	fromFileTimeout string
)

var fromFileCmd = &cobra.Command{
	Use:   "from-file <file>",
	Short: "Execute bulk operations from a file or stdin",
	Args:  cobra.MaximumNArgs(1),
	Example: utils.TrimAndIndent(`
	# Execute bulk operations from a file
	esctl bulk from-file operations.ndjson

	# Execute bulk operations from stdin
	cat operations.ndjson | esctl bulk from-file

	# Execute with specific index and refresh
	esctl bulk from-file operations.ndjson --index=my-index --refresh=wait_for
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		var filename string
		if len(args) > 0 {
			filename = args[0]
		}

		return handleBulkFromFile(ctx, filename)
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate [output-file]",
	Short: "Generate a sample bulk NDJSON template",
	Args:  cobra.MaximumNArgs(1),
	Example: utils.TrimAndIndent(`
	# Generate template to stdout
	esctl bulk generate

	# Generate template to file
	esctl bulk generate template.ndjson
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		var filename string
		if len(args) > 0 {
			filename = args[0]
		}

		return handleBulkGenerate(filename)
	},
}

func init() {
	fromFileCmd.Flags().StringVar(&fromFileIndex, "index", "", "Default index for operations that don't specify one")
	fromFileCmd.Flags().StringVar(&fromFileRefresh, "refresh", "", "Refresh policy (true, false, wait_for)")
	fromFileCmd.Flags().StringVar(&fromFileTimeout, "timeout", "", "Timeout for the bulk operation")

	bulkCmd.AddCommand(fromFileCmd)
	bulkCmd.AddCommand(generateCmd)
}

func Cmd() *cobra.Command {
	return bulkCmd
}

func handleBulkFromFile(ctx context.Context, filename string) error {
	var file *os.File
	var err error

	if filename == "" || filename == "-" {
		// Read from stdin
		file = os.Stdin
		fmt.Println("Reading bulk operations from stdin...")
	} else {
		// Read from file
		file, err = os.Open(filename)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()
		fmt.Printf("Reading bulk operations from %s...\n", filename)
	}

	result, err := bulk.ExecuteBulk(ctx, file, fromFileIndex, fromFileRefresh, fromFileTimeout)
	if err != nil {
		return fmt.Errorf("bulk operation failed: %w", err)
	}

	// Print summary only for table output to avoid mixing with JSON/YAML.
	if shared.OutputFormat == "table" {
		fmt.Printf("Bulk operation completed in %dms\n", result.Took)
	}

	successCount := 0
	errorCount := 0
	for _, item := range result.Items {
		for _, op := range item {
			if op.Error != nil {
				errorCount++
			} else {
				successCount++
			}
		}
	}

	if shared.OutputFormat == "table" {
		fmt.Printf("Operations: %d successful, %d errors\n", successCount, errorCount)
	}

	if result.Errors && shared.OutputFormat == "table" {
		fmt.Println("Errors occurred during bulk operation:")
		for i, item := range result.Items {
			for opType, op := range item {
				if op.Error != nil {
					fmt.Printf("  Item %d (%s): %s - %s\n", i+1, opType, op.Error.Type, op.Error.Reason)
				}
			}
		}
	}

	return output.Render(result)
}

func handleBulkGenerate(filename string) error {
	template := bulk.GenerateBulkTemplate()

	if filename == "" {
		// Output to stdout
		fmt.Print(template)
	} else {
		// Output to file
		err := os.WriteFile(filename, []byte(template), 0644)
		if err != nil {
			return fmt.Errorf("failed to write template file: %w", err)
		}
		fmt.Printf("Bulk template written to %s\n", filename)
	}

	return nil
}
