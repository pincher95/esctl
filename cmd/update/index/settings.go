package index

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/index"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagBody         string
	flagSettings     string
	flagFlatSettings bool
	flagIndex        string
)

var SettingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Updates the settings of an index",
	Long: utils.Trim(`
	Updates the settings of an index. You can update the number of replicas, the number of shards, and other settings.
	`),
	Example: utils.TrimAndIndent(`
	# Update the number of replicas for an index.
	esctl update index --index my_index --replicas 2
	`),
	Run: func(cmd *cobra.Command, args []string) {
		idxClient := index.NewIndex()

		ctx := cmd.Context()

		// Split the user input into key and value
		kv := strings.SplitN(flagBody, "=", 2)
		if len(kv) != 2 {
			log.Fatalf("Invalid --setting format. Must be key=value, got: %s", flagBody)
		}
		key, value := kv[0], kv[1]

		// Our top-level container
		// You might keep it at just `root := make(map[string]any)`
		// but ES typically expects the "index" block at the top for these settings.
		body := make(map[string]any, 0)
		if flagSettings == "" {
			applySetting(&body, &key, &value)
		}

		handleIndexLogic(ctx, idxClient, body)
	},
}

func init() {
	// SettingsCmd.Flags().BoolVar(&flagFlatSettings, "no-flat-settings", false, "If set, print settings in a none flat format (Default is false)")
	// updateIndexCmd.Flags().IntVar(&flagReplicas, "replicas", 0, "Number of replicas for the index")
	// updateIndexCmd.Flags().IntVar(&flagShards, "shards", 0, "Number of shards for the index")
	SettingsCmd.Flags().StringVar(&flagBody, "settings", "", "Index settings")

	_ = SettingsCmd.MarkFlagRequired("settings")
}

func handleIndexLogic(ctx context.Context, idx index.Index, body map[string]any) {
	// fmt.Fprintf(os.Stderr, "This operation will update the settings of the index: %s.", flagIndex)
	approved, err := utils.GetApproval()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	if approved {
		settings, err := idx.UpdateIndexSettings(ctx, flagIndex, body, flagFlatSettings)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to update index:", err)
			os.Exit(1)
		}

		output.PrintJson(&settings)
	}
}

// applySetting creates nested maps in 'root' based on dot-separated key segments
// and sets the final segment to 'value'.
func applySetting(root *map[string]any, key *string, value *string) {
	parts := strings.Split(*key, ".")
	currentMap := *root

	for i, part := range parts {
		// If we're at the last part, set the value
		if i == len(parts)-1 {
			currentMap[part] = *value
			return
		}

		// Otherwise, descend into the map (creating if needed)
		next, exists := currentMap[part]
		if !exists {
			// Create a new map for this level
			nested := map[string]any{}
			currentMap[part] = nested
			currentMap = nested
		} else {
			// If it already exists, it should be a map
			nm, ok := next.(map[string]any)
			if !ok {
				// This means there's a conflict in structure (unlikely for well-formed settings)
				// but we handle it gracefully by overwriting with a fresh map
				nm = map[string]any{}
				currentMap[part] = nm
			}
			currentMap = nm
		}
	}
}
