package update

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/cluster"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var (
	flagRerouteIndex          string
	flagRerouteShard          int
	flagRerouteNode           string
	flagRerouteFromNode       string
	flagRerouteToNode         string
	flagRerouteAcceptDataLoss bool
	flagRerouteAllowPrimary   bool
	flagRerouteCmdDryRun      bool
	flagRerouteCmdExplain     bool
)

var rerouteAllocateStalePrimaryCmd = &cobra.Command{
	Use:   "allocate-stale-primary",
	Short: "Allocate a stale primary shard copy to a node (may lose recent writes)",
	Long: utils.Trim(`
Force-allocate a primary shard using a stale copy on the given node. Use this only when the
in-sync copy is permanently lost (for example its node will not return) and you have located a
surviving copy with 'esctl get shard-stores'. Recent writes not present on the stale copy are
lost, so --accept-data-loss must be set explicitly.`),
	Example: utils.TrimAndIndent(`
# Locate a surviving copy, then force-allocate it as primary.
esctl get shard-stores --status red
esctl update reroute allocate-stale-primary --index my-index --shard 0 --node es-data-1 --accept-data-loss

# Preview the effect without applying it.
esctl update reroute allocate-stale-primary --index my-index --shard 0 --node es-data-1 --accept-data-loss --dry-run
`),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !flagRerouteAcceptDataLoss {
			return fmt.Errorf("--accept-data-loss must be set: allocating a stale primary can lose recent writes")
		}
		return runRerouteCommand(cmd.Context(), map[string]any{
			"allocate_stale_primary": map[string]any{
				"index":            flagRerouteIndex,
				"shard":            flagRerouteShard,
				"node":             flagRerouteNode,
				"accept_data_loss": true,
			},
		})
	},
}

var rerouteAllocateEmptyPrimaryCmd = &cobra.Command{
	Use:   "allocate-empty-primary",
	Short: "Allocate an empty primary shard to a node (DISCARDS all data for the shard)",
	Long: utils.Trim(`
Force-allocate a brand-new EMPTY primary shard on the given node. This permanently discards all
data for the shard and should only be used as a last resort when every copy is unrecoverable and
you accept losing the shard's data to restore the index to a writable (green) state.
--accept-data-loss must be set explicitly.`),
	Example: utils.TrimAndIndent(`
# Last resort: no recoverable copy exists; create an empty primary to unblock the index.
esctl update reroute allocate-empty-primary --index my-index --shard 0 --node es-data-1 --accept-data-loss
`),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !flagRerouteAcceptDataLoss {
			return fmt.Errorf("--accept-data-loss must be set: allocating an empty primary discards all data for the shard")
		}
		return runRerouteCommand(cmd.Context(), map[string]any{
			"allocate_empty_primary": map[string]any{
				"index":            flagRerouteIndex,
				"shard":            flagRerouteShard,
				"node":             flagRerouteNode,
				"accept_data_loss": true,
			},
		})
	},
}

var rerouteMoveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move a started shard from one node to another",
	Long: utils.Trim(`
Move an already-started shard from one node to another, for example to drain a node before
maintenance or to relieve a hot node.`),
	Example: utils.TrimAndIndent(`
esctl update reroute move --index my-index --shard 0 --from-node es-data-0 --to-node es-data-2
`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRerouteCommand(cmd.Context(), map[string]any{
			"move": map[string]any{
				"index":     flagRerouteIndex,
				"shard":     flagRerouteShard,
				"from_node": flagRerouteFromNode,
				"to_node":   flagRerouteToNode,
			},
		})
	},
}

var rerouteCancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel allocation or recovery of a shard on a node",
	Long: utils.Trim(`
Cancel the allocation (or ongoing recovery) of a shard on a node. Use --allow-primary to cancel a
primary shard's allocation, which discards its data on that node.`),
	Example: utils.TrimAndIndent(`
# Cancel a stuck replica recovery.
esctl update reroute cancel --index my-index --shard 0 --node es-data-2

# Cancel a primary allocation (discards the copy on that node).
esctl update reroute cancel --index my-index --shard 0 --node es-data-2 --allow-primary
`),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRerouteCommand(cmd.Context(), map[string]any{
			"cancel": map[string]any{
				"index":         flagRerouteIndex,
				"shard":         flagRerouteShard,
				"node":          flagRerouteNode,
				"allow_primary": flagRerouteAllowPrimary,
			},
		})
	},
}

func runRerouteCommand(ctx context.Context, command map[string]any) error {
	result, err := cluster.ClusterRerouteCommands(ctx, []map[string]any{command}, flagRerouteCmdDryRun, flagRerouteCmdExplain)
	if err != nil {
		return fmt.Errorf("failed to execute reroute command: %w", err)
	}
	return output.Render(result)
}

// addRerouteCommonFlags registers flags shared by every reroute subcommand.
func addRerouteCommonFlags(c *cobra.Command) {
	c.Flags().StringVar(&flagRerouteIndex, "index", "", "Index the shard belongs to (required)")
	c.Flags().IntVar(&flagRerouteShard, "shard", 0, "Shard number (required)")
	c.Flags().BoolVar(&flagRerouteCmdDryRun, "dry-run", false, "Simulate the command and return the resulting state without applying it")
	c.Flags().BoolVar(&flagRerouteCmdExplain, "explain", true, "Include an explanation of why the command can or cannot be executed")
	_ = c.MarkFlagRequired("index")
	_ = c.MarkFlagRequired("shard")
}

func init() {
	for _, c := range []*cobra.Command{
		rerouteAllocateStalePrimaryCmd,
		rerouteAllocateEmptyPrimaryCmd,
		rerouteMoveCmd,
		rerouteCancelCmd,
	} {
		addRerouteCommonFlags(c)
	}

	// Node-targeted commands.
	rerouteAllocateStalePrimaryCmd.Flags().StringVar(&flagRerouteNode, "node", "", "Node to allocate the shard on (required)")
	rerouteAllocateStalePrimaryCmd.Flags().BoolVar(&flagRerouteAcceptDataLoss, "accept-data-loss", false, "Acknowledge that recent writes may be lost (required)")
	_ = rerouteAllocateStalePrimaryCmd.MarkFlagRequired("node")

	rerouteAllocateEmptyPrimaryCmd.Flags().StringVar(&flagRerouteNode, "node", "", "Node to allocate the shard on (required)")
	rerouteAllocateEmptyPrimaryCmd.Flags().BoolVar(&flagRerouteAcceptDataLoss, "accept-data-loss", false, "Acknowledge that all data for the shard will be discarded (required)")
	_ = rerouteAllocateEmptyPrimaryCmd.MarkFlagRequired("node")

	rerouteCancelCmd.Flags().StringVar(&flagRerouteNode, "node", "", "Node the shard copy is on (required)")
	rerouteCancelCmd.Flags().BoolVar(&flagRerouteAllowPrimary, "allow-primary", false, "Allow cancelling a primary shard's allocation")
	_ = rerouteCancelCmd.MarkFlagRequired("node")

	rerouteMoveCmd.Flags().StringVar(&flagRerouteFromNode, "from-node", "", "Node to move the shard from (required)")
	rerouteMoveCmd.Flags().StringVar(&flagRerouteToNode, "to-node", "", "Node to move the shard to (required)")
	_ = rerouteMoveCmd.MarkFlagRequired("from-node")
	_ = rerouteMoveCmd.MarkFlagRequired("to-node")

	updateRerouteCmd.AddCommand(rerouteAllocateStalePrimaryCmd)
	updateRerouteCmd.AddCommand(rerouteAllocateEmptyPrimaryCmd)
	updateRerouteCmd.AddCommand(rerouteMoveCmd)
	updateRerouteCmd.AddCommand(rerouteCancelCmd)
}
