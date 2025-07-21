package get

import (
	"context"
	"fmt"
	"sort"
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

var baseHealthColumns = []output.ColumnDefaults{
	{Header: "STATUS", Type: output.Text},
	{Header: "NODE_TOTAL", Type: output.Number},
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
	// build dynamic role columns
	roleCounts := countRoles(ctx, client)
	// get sorted unique role keys
	roles := make([]string, 0, len(roleCounts))
	for r := range roleCounts {
		roles = append(roles, r)
	}
	sort.Strings(roles)

	// Build headers in upper-case but keep original keys for lookup
	roleHeaders := make([]string, len(roles))
	for i, r := range roles {
		roleHeaders[i] = strings.ToUpper(r)
	}

	// build column definitions dynamically
	dynamicCols := make([]output.ColumnDefaults, 0, len(roleHeaders))
	for _, h := range roleHeaders {
		dynamicCols = append(dynamicCols, output.ColumnDefaults{Header: h, Type: output.Number})
	}

	allColumns := append([]output.ColumnDefaults{}, baseHealthColumns...) // copy
	// insert role columns after NODE_TOTAL (index 2)
	allColumns = append(allColumns[:2], append(dynamicCols, allColumns[2:]...)...)

	columnDefs, _ := getColumnDefs(conf, "STATUS", allColumns)

	// build row matching columns order
	row := []string{
		h.Status,
		fmt.Sprintf("%d", h.NodeTotal),
	}
	for _, r := range roles {
		row = append(row, fmt.Sprintf("%d", roleCounts[r]))
	}
	row = append(row,
		fmt.Sprintf("%d", h.Shards),
		fmt.Sprintf("%d", h.Relo),
		fmt.Sprintf("%d", h.Init),
		fmt.Sprintf("%d", h.Unassign),
		fmt.Sprintf("%d", h.PendingTasks),
		h.ActiveShardsPercent,
	)

	output.PrintTable(columnDefs, [][]string{row}, nil)
	return nil
}

func countRoles(ctx context.Context, c cat.Cat) map[string]int {
	m := make(map[string]int)
	nodes, err := c.CatNodes(ctx, "", "", "", "")
	if err != nil {
		return m
	}
	for _, n := range nodes {
		for _, r := range n.RolesList() {
			m[r]++
		}
	}
	return m
}

func init() {
	// attach to get root by referencing variable in get package init order.
}
