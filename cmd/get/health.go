package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getHealthCmd = &cobra.Command{
	Use:   "health",
	Short: "Get cluster health in cat format",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cat.NewCat()
		conf := config.ParseConfigFile()
		ctx := cmd.Context()

		if !flagRefresh {
			return handleHealth(ctx, client, *conf)
		}
		return utils.WatchLoop(flagRefreshInterval, func() error {
			return handleHealth(ctx, client, *conf)
		})
	},
}

var healthColumns = []output.ColumnDefaults{
	{Header: "STATUS", Type: output.Text},
	{Header: "NODE_TOTAL", Type: output.Number},
	{Header: "DATA", Type: output.Number},
	{Header: "MASTER", Type: output.Number},
	{Header: "SHARDS", Type: output.Number},
	{Header: "RELOCATING", Type: output.Number},
	{Header: "INIT", Type: output.Number},
	{Header: "UNASSIGN", Type: output.Number},
	{Header: "PENDING_TASKS", Type: output.Number},
	{Header: "ACTIVE_PERC", Type: output.Text},
}

func handleHealth(ctx context.Context, client cat.Cat, conf config.Config) error {
	h, err := client.CatHealth(ctx)
	if err != nil {
		return err
	}
	columnDefs, _ := getColumnDefs(conf, "STATUS", healthColumns)
	row := []string{
		h.Status,
		fmt.Sprintf("%d", h.NodeTotal),
		fmt.Sprintf("%d", countRole(ctx, client)),
		fmt.Sprintf("%d", countMaster(ctx, client)),
		fmt.Sprintf("%d", h.Shards),
		fmt.Sprintf("%d", h.Relo),
		fmt.Sprintf("%d", h.Init),
		fmt.Sprintf("%d", h.Unassign),
		fmt.Sprintf("%d", h.PendingTasks),
		h.ActiveShardsPercent,
	}
	output.PrintTable(columnDefs, [][]string{row}, nil)
	return nil
}

// countRole counts data-role nodes
func countRole(ctx context.Context, c cat.Cat) int {
	nodes, err := c.CatNodes(ctx, "", "", "", "")
	if err != nil {
		return 0
	}
	cnt := 0
	for _, n := range nodes {
		if strings.Contains(n.Roles, "d") {
			cnt++
		}
	}
	return cnt
}

func countMaster(ctx context.Context, c cat.Cat) int {
	nodes, err := c.CatNodes(ctx, "", "", "", "")
	if err != nil {
		return 0
	}
	cnt := 0
	for _, n := range nodes {
		if strings.Contains(n.Roles, "m") || n.Master == "*" {
			cnt++
		}
	}
	return cnt
}

func init() {
	// attach to get root by referencing variable in get package init order.
}
