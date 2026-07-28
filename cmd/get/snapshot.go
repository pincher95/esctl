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
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

var getSnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Get snapshot details or list snapshots",
	Args:  cobra.NoArgs,
	Example: utils.TrimAndIndent(`
	# Get a specific snapshot
	esctl get snapshot --repository my-repo --name my-snapshot

	# List snapshots in a repository
	esctl get snapshot --repository my-repo

	# List snapshots using filters
	esctl get snapshot --repository my-repo --status SUCCESS

	# List snapshots by name substring
	esctl get snapshot --repository my-repo --name snap
	`),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if flagRepository == "" {
			return fmt.Errorf("repository is required")
		}
		return runWithWatch(ctx, func() error {
			if flagSnapshotName != "" {
				if isTableOutput() {
					return handleSnapshotList(ctx, flagRepository)
				}
				return snapshot.HandleSnapshotGet(ctx, flagRepository, flagSnapshotName)
			}
			return handleSnapshotList(ctx, flagRepository)
		})
	},
}

func init() {
	getSnapshotCmd.Flags().StringVarP(&flagRepository, "repository", "p", "", "Name of the snapshot repository")
	getSnapshotCmd.Flags().StringVar(&flagStatus, "status", "", "Filter snapshots by status")
	getSnapshotCmd.Flags().StringVar(&flagSnapshotName, "name", "", "Snapshot name")
	getSnapshotCmd.MarkFlagRequired("repository")
}

var (
	flagSnapshotName string
)

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

	filter := ""
	if flagSnapshotName != "" {
		filter = flagSnapshotName
	}
	snapshots, err := snapshotsClient.CatSnapshots(ctx, "", repository, filter)
	if err != nil {
		return fmt.Errorf("failed to retrieve snapshots in repository %q: %w", repository, err)
	}

	columnDefs, err := getColumnDefs(*conf, "end_epoch", snapshotsColumns)
	if err != nil {
		return fmt.Errorf("failed to get column definitions: %w", err)
	}

	data := make([][]string, 0, len(snapshots))

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

func isTableOutput() bool {
	switch strings.ToLower(shared.OutputFormat) {
	case "", "table":
		return true
	default:
		return false
	}
}
