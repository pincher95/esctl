package alias

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/index"
	"github.com/spf13/cobra"
)

var (
	addIndices string
	addFilter  string
	addRouting string
)

var addCmd = &cobra.Command{
	Use:   "add <alias> --indices=<index1,index2,...>",
	Short: "Add an alias to one or more indices",
	Args:  cobra.ExactArgs(1),
	Example: utils.TrimAndIndent(`
	# Add alias to a single index
	esctl alias add my-alias --indices=my-index

	# Add alias to multiple indices
	esctl alias add logs-current --indices="logs-2023,logs-2024"

	# Add alias with filter
	esctl alias add active-logs --indices=logs --filter='{"term":{"status":"active"}}'

	# Add alias with routing
	esctl alias add user-data --indices=users --routing=user_id
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleAliasAdd(ctx, args[0])
	},
}

func init() {
	addCmd.Flags().StringVar(&addIndices, "indices", "", "Comma-separated list of indices to add the alias to")
	addCmd.Flags().StringVar(&addFilter, "filter", "", "Filter to apply to the alias (JSON)")
	addCmd.Flags().StringVar(&addRouting, "routing", "", "Routing value for the alias")
	addCmd.MarkFlagRequired("indices")
}

func handleAliasAdd(ctx context.Context, alias string) error {
	if addIndices == "" {
		return fmt.Errorf("indices must be specified")
	}

	indices := strings.Split(addIndices, ",")
	for i, idx := range indices {
		indices[i] = strings.TrimSpace(idx)
	}

	var filter map[string]interface{}
	if addFilter != "" {
		if err := json.Unmarshal([]byte(addFilter), &filter); err != nil {
			return fmt.Errorf("invalid filter JSON: %w", err)
		}
	}

	err := index.AddAlias(ctx, indices, alias, filter, addRouting)
	if err != nil {
		return fmt.Errorf("failed to add alias: %w", err)
	}

	fmt.Printf("Successfully added alias '%s' to indices: %s\n", alias, strings.Join(indices, ", "))
	return nil
}
