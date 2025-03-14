package get

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/utils"
	"github.com/pincher95/esctl/es/cat"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var getSnapshotsCmd = &cobra.Command{
	Use:   "snapshots",
	Short: "Get Elasticsearch snapshots",
	Long: utils.Trim(`
	Get Elasticsearch snapshots. You can filter the results using the index flag.
	`),
	Example: utils.TrimAndIndent(`
	# Retrieve all snapshots.
	esctl get snapshots

	# Retrieve snapshots for a specific index.
	esctl get snapshots --repository my_index
	`),
	Run: func(cmd *cobra.Command, args []string) {
		snapshotsClient := cat.NewCat()
		config := config.ParseConfigFile()

		// If --watch is NOT set, just run once
		if !flagRefresh {
			handleSnapshotsLogic(snapshotsClient, *config)
			return
		}

		// If --watch is set, run in a loop
		for {
			clearScreen() // optional, to mimic "watch" clearing
			handleSnapshotsLogic(snapshotsClient, *config)
			time.Sleep(flagRefreshInterval)
		}
	},
}

func init() {
	getSnapshotsCmd.PersistentFlags().StringVarP(&flagRepository, "repository", "p", "", "Name of the snapshot repository")
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

func handleSnapshotsLogic(snapshotsClient cat.Cat, conf config.Config) {
	snapshots, err := snapshotsClient.CatSnapshots(nil, &flagRepository)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to retrieve indices:", err)
		os.Exit(1)
	}

	columnDefs, err := getColumnDefs(conf, "end_epoch", snapshotsColumns)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get column definitions:", err)
		os.Exit(1)
	}

	data := [][]string{}

	for _, snapshot := range *snapshots {
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

	if len(flagSortBy) > 0 {
		sortCols := output.ParseSortColumns(flagSortBy)
		output.PrintTable(columnDefs, data, sortCols)
	} else {
		sortCols := output.ParseSortColumns("ID")
		output.PrintTable(columnDefs, data, sortCols)
	}
}
