package get

import (
	"context"
	"fmt"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/node"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var flagThreadPoolsNode string

var getThreadPoolsCmd = &cobra.Command{
	Use:   "thread-pools",
	Short: "Get thread pool statistics",
	Long: utils.Trim(`
		Retrieves thread pool statistics for Elasticsearch nodes. Shows the number of threads,
		queue size, active threads, rejected tasks, and completed tasks for each thread pool.
	`),
	Example: utils.TrimAndIndent(`
		# Get thread pools for all nodes
		esctl get thread-pools

		# Get thread pools for a specific node
		esctl get thread-pools --node node-1

		# Get thread pools in JSON format
		esctl get thread-pools -o json

		# Watch thread pools in real-time
		esctl get thread-pools --watch
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return handleGetThreadPoolsLogic(ctx, flagThreadPoolsNode)
	},
}

func init() {
	getThreadPoolsCmd.Flags().StringVar(&flagThreadPoolsNode, "node", "", "Node ID or name to filter by")
}

func handleGetThreadPoolsLogic(ctx context.Context, nodeID string) error {
	result, err := node.GetThreadPools(ctx, nodeID)
	if err != nil {
		return fmt.Errorf("failed to get thread pools: %w", err)
	}

	// Convert to a flat structure for table output
	type threadPoolRow struct {
		Node      string `json:"node" yaml:"node"`
		Pool      string `json:"pool" yaml:"pool"`
		Threads   int    `json:"threads" yaml:"threads"`
		Queue     int    `json:"queue" yaml:"queue"`
		Active    int    `json:"active" yaml:"active"`
		Rejected  int    `json:"rejected" yaml:"rejected"`
		Largest   int    `json:"largest" yaml:"largest"`
		Completed int    `json:"completed" yaml:"completed"`
	}

	var rows []threadPoolRow
	for nodeID, nodeData := range result.Nodes {
		nodeName := nodeData.Name
		if nodeName == "" {
			nodeName = nodeID
		}

		for poolName, poolStats := range nodeData.ThreadPools {
			rows = append(rows, threadPoolRow{
				Node:      nodeName,
				Pool:      poolName,
				Threads:   poolStats.Threads,
				Queue:     poolStats.Queue,
				Active:    poolStats.Active,
				Rejected:  poolStats.Rejected,
				Largest:   poolStats.Largest,
				Completed: poolStats.Completed,
			})
		}
	}

	if len(rows) == 0 {
		fmt.Println("No thread pool data found")
		return nil
	}

	return output.Render(rows)
}
