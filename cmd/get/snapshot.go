package get

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/snapshot"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getSnapshotCmd = &cobra.Command{
	Use:   "snapshot [repository] [snapshot]",
	Short: "Get snapshot details or list snapshots",
	Args:  cobra.RangeArgs(0, 2),
	Example: utils.TrimAndIndent(`
	# Get a specific snapshot
	esctl get snapshot my-repo my-snapshot

	# List snapshots in a repository
	esctl get snapshot my-repo

	# List snapshots using flags
	esctl get snapshot --repository my-repo --status SUCCESS
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		switch len(args) {
		case 2:
			if flagRepository != "" && flagRepository != args[0] {
				return fmt.Errorf("repository specified twice: %s and %s", flagRepository, args[0])
			}
			return snapshot.HandleSnapshotGet(ctx, args[0], args[1])
		case 1:
			if flagRepository != "" && flagRepository != args[0] {
				return fmt.Errorf("repository specified twice: %s and %s", flagRepository, args[0])
			}
			return handleSnapshotList(ctx, args[0])
		default:
			if flagRepository == "" {
				return fmt.Errorf("repository is required to list snapshots")
			}
			return handleSnapshotList(ctx, flagRepository)
		}
	},
}

func init() {
	getSnapshotCmd.Flags().StringVarP(&flagRepository, "repository", "p", "", "Name of the snapshot repository")
	getSnapshotCmd.Flags().StringVar(&flagStatus, "status", "", "Filter snapshots by status")
	getSnapshotCmd.Flags().StringVar(&flagFilter, "name", "", "Filter snapshots by name or substring of snapshot name e.g. 'snapshot-1', 'snapshot', 'data'")
}

var snapshotsColumns = []output.ColumnDefaults{
	{Header: "ID", Type: output.Text},
	{Header: "STATUS", Type: output.Text},
	{Header: "START-TIME", Type: output.Text},
	{Header: "START-DATE", Type: output.Date},
	{Header: "END-DATE", Type: output.Date},
	{Header: "DURATION", Type: output.Text},
	{Header: "INDICES", Type: output.Number},
	{Header: "SUCCESSFUL-SHARDS", Type: output.Number},
	{Header: "FAILD-SHARDS", Type: output.Number},
	{Header: "TOTAL-SHARDS", Type: output.Number},
}

func handleSnapshotList(ctx context.Context, repository string) error {
	snapshotsClient := cat.NewCat()
	conf, err := config.ParseConfigFile()
	if err != nil {
		return err
	}

	snapshots, err := snapshotsClient.CatSnapshots(ctx, "", repository, flagFilter)
	if err != nil {
		return fmt.Errorf("failed to retrieve indices: %w", err)
	}

	columnDefs, err := getColumnDefs(*conf, "end_epoch", snapshotsColumns)
	if err != nil {
		return fmt.Errorf("failed to get column definitions: %w", err)
	}

	data := [][]string{}

	for _, snapshot := range snapshots {
		if includeSnapshotByStatus(snapshot) {
			rowData := map[string]string{
				"ID":                snapshot.ID,
				"STATUS":            snapshot.Status,
				"START-TIME":        snapshot.StartTime,
				"START-DATE":        time.Unix(int64(snapshot.StartEpoch), 0).Format(time.RFC822Z),
				"END-DATE":          time.Unix(int64(snapshot.EndEpoch), 0).Format(time.RFC822Z),
				"DURATION":          snapshot.Duration,
				"INDICES":           strconv.Itoa(utils.SafeInt(&snapshot.Indices)),
				"SUCCESSFUL-SHARDS": strconv.Itoa(utils.SafeInt(&snapshot.SuccessfulShards)),
				"FAILD-SHARDS":      strconv.Itoa(utils.SafeInt(&snapshot.FailedShards)),
				"TOTAL-SHARDS":      strconv.Itoa(utils.SafeInt(&snapshot.TotalShards)),
			}

			row := make([]string, len(columnDefs))
			for i, colDef := range columnDefs {
				row[i] = rowData[colDef.Header]
			}
			data = append(data, row)
		}
	}

	if len(flagSortBy) > 0 {
		sortCols := output.ParseSortColumns(flagSortBy)
		return output.PrintTable(columnDefs, data, sortCols)
	}
	sortCols := output.ParseSortColumns("ID")
	return output.PrintTable(columnDefs, data, sortCols)
}

func includeSnapshotByStatus(snapshot cat.CatSnapshotResponse) bool {
	return snapshot.Status == strings.ToUpper(flagStatus) || flagStatus == ""
}
