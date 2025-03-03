package update

import (
	"fmt"
	"os"

	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/cluster"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var updateRerouteCmd = &cobra.Command{
	Use:   "reroute",
	Short: "Changes the allocation of shards in a cluster",
	Long: utils.Trim(`
	The reroute command allows for manual changes to the allocation of individual shards in the cluster..
	`),
	Example: utils.TrimAndIndent(`
	# Reroute the shards in the cluster.
	esctl update reroute

	# Reroute the shards in the cluster with a dry-run.
	esctl update reroute --dry-run

	# Reroute the shards in the cluster with an explanation.
	esctl update reroute --explain

	# Reroute the shards in the cluster with a retry-failed.
	esctl update reroute --retry-failed

	# Reroute the shards in the cluster with a metric.
	esctl update reroute --metric '_all'

	# Reroute the shards in the cluster with a dry-run, explanation, retry-failed, and metric.
	esctl update reroute --dry-run --explain --retry-failed --metric 'none'
	`),
	Run: func(cmd *cobra.Command, args []string) {
		handleRerouteLogic()
	},
}

func init() {
	updateRerouteCmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "Simulates the operation only and returns the resulting state (Default: false)")
	updateRerouteCmd.Flags().BoolVar(&flagExplain, "explain", true, "contains an explanation of why the commands can or cannot be executed (Default: true)")
	updateRerouteCmd.Flags().BoolVar(&flagRetryFailed, "retry-failed", true, "Retries allocation of shards that are blocked due to too many subsequent allocation failures (Default: true)")
	updateRerouteCmd.Flags().StringVar(&flagMertic, "metric", "none", "Limits the information returned to the specified metrics (Default: none)")
}

func handleRerouteLogic() {
	reroute, err := cluster.ClusterReroute(nil, &flagMertic, flagDryRun, flagExplain, flagRetryFailed)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve reroute:", err)
		os.Exit(1)
	}

	output.PrintJson(reroute)
}
