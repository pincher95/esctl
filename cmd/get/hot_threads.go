package get

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/node"
	"github.com/spf13/cobra"
)

var (
	flagHotThreadsNode     string
	flagHotThreadsCount    int
	flagHotThreadsInterval string
	flagHotThreadsType     string
)

var getHotThreadsCmd = &cobra.Command{
	Use:   "hot-threads",
	Short: "Get hot threads diagnostics",
	Long: utils.Trim(`
		Retrieves information about hot threads on nodes in the cluster. Shows which threads
		are consuming the most CPU time, which can help identify performance issues.

		Thread types:
		  - cpu: Threads with high CPU usage (default)
		  - wait: Threads that are waiting
		  - block: Threads that are blocked
	`),
	Example: utils.TrimAndIndent(`
		# Get hot threads for all nodes
		esctl get hot-threads

		# Get hot threads for a specific node
		esctl get hot-threads --node node-1

		# Get top 5 hot threads with 1 second sampling interval
		esctl get hot-threads --threads 5 --interval 1s

		# Get hot threads by type (cpu, wait, or block)
		esctl get hot-threads --type wait

		# Combine multiple options
		esctl get hot-threads --node node-1 --threads 10 --interval 500ms --type cpu
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleGetHotThreadsLogic(ctx, flagHotThreadsNode, flagHotThreadsCount, flagHotThreadsInterval, flagHotThreadsType)
	},
}

func init() {
	getHotThreadsCmd.Flags().StringVar(&flagHotThreadsNode, "node", "", "Node ID or name to filter by")
	getHotThreadsCmd.Flags().IntVar(&flagHotThreadsCount, "threads", 3, "Number of hot threads to return")
	getHotThreadsCmd.Flags().StringVar(&flagHotThreadsInterval, "interval", "", "Sampling interval (e.g., 500ms, 1s)")
	getHotThreadsCmd.Flags().StringVar(&flagHotThreadsType, "type", "", "Thread type: cpu, wait, or block")
}

func handleGetHotThreadsLogic(ctx context.Context, nodeID string, threads int, interval, threadType string) error {
	// Validate thread type if provided
	if threadType != "" && threadType != "cpu" && threadType != "wait" && threadType != "block" {
		return fmt.Errorf("invalid thread type: %s (must be cpu, wait, or block)", threadType)
	}

	result, err := node.GetHotThreads(ctx, nodeID, threads, interval, threadType)
	if err != nil {
		return fmt.Errorf("failed to get hot threads: %w", err)
	}

	// Hot threads returns plain text, just print it
	fmt.Println(result)
	return nil
}
