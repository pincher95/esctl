package describe

import (
	"fmt"
	"os"
	"strings"

	"github.com/pincher95/esctl/constants"
	"github.com/pincher95/esctl/es"
	"github.com/pincher95/esctl/output"
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Short:     "Print detailed information about an entity",
	Args:      cobra.RangeArgs(1, 2),
	ValidArgs: []string{"cluster", "index", "node"},
	Run: func(cmd *cobra.Command, args []string) {
		entity := args[0]
		switch entity {
		case constants.EntityCluster:
			handleDescribeCluster()
		case constants.EntityIndex:
			if len(args) < 2 {
				fmt.Println("Index name is required.")
				cmd.Help()
				os.Exit(1)
			}
			handleDescribeIndex(args[1])
		case constants.EntityNode:
			node := ""
			if len(args) == 2 {
				node = args[1]
			}
			handleDescribeNode(node)
		default:
			fmt.Printf("Unknown entity: %s\n", entity)
			cmd.Help()
			os.Exit(1)
		}
	},
}

func Cmd() *cobra.Command {
	return describeCmd
}

func handleDescribeCluster() {
	cluster, err := es.GetCluster()
	if err != nil {
		fmt.Println("Failed to retrieve cluster information:", err)
		return
	}

	print(cluster)
}

func handleDescribeIndex(index string) {
	shouldGetMappings := flagMappings || !flagSettings
	shouldGetSettings := flagSettings || !flagMappings

	details, err := es.GetIndexDetails(index, shouldGetMappings, shouldGetSettings)
	if err != nil {
		fmt.Println("Failed to retrieve index details:", err)
		return
	}

	print(details)
}

func handleDescribeNode(node string) {
	nodeDetails, err := es.GetNodeDetails(node)
	if err != nil {
		fmt.Println("Failed to retrieve node details:", err)
		return
	}

	print(nodeDetails)
}

func print(data interface{}) {
	switch flagOutput {
	case "json":
		output.PrintJson(data)
	case "yaml":
		output.PrintYaml(data)
	default:
		fmt.Printf("Unknown output type: %s\n", flagOutput)
		os.Exit(1)
	}
}

func init() {
	describeCmd.Use = fmt.Sprintf(`describe [%s] [NAME]`, strings.Join(describeCmd.ValidArgs, "|"))
	describeCmd.Long = fmt.Sprintf("Print detailed information about the specified entity.\nAvailable entities: %s.", strings.Join(describeCmd.ValidArgs, ", "))

	describeCmd.Flags().BoolVar(&flagMappings, "mappings", false, "If set, retrieve and print index mappings")
	describeCmd.Flags().BoolVar(&flagSettings, "settings", false, "If set, retrieve and print index settings")
	describeCmd.Flags().StringVarP(&flagOutput, "output", "o", "json", "Print output as json or yaml")
}
