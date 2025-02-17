package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/count"
	"github.com/pincher95/esctl/cmd/describe"
	"github.com/pincher95/esctl/cmd/get"
	"github.com/pincher95/esctl/cmd/query"
	"github.com/pincher95/esctl/constants"
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "esctl",
	Short: "esctl is CLI for Elasticsearch",
	Long:  `esctl is a read-only CLI for Elasticsearch that allows users to manage and monitor their Elasticsearch clusters.`,
}

func Execute(ctx context.Context) error {
	return RootCmd.ExecuteContext(ctx)
}

func init() {
	cobra.OnInitialize(initialize)

	initProtocolFlag()
	initHostFlag()
	initPortFlag()
	initUsernameFlag()
	initPasswordFlag()

	RootCmd.PersistentFlags().StringVar(&shared.Context, "context", "", "Override context")
	RootCmd.PersistentFlags().BoolVar(&shared.Debug, "debug", false, "Enable debug mode")

	RootCmd.AddCommand(config.Cmd())
	RootCmd.AddCommand(count.Cmd())
	RootCmd.AddCommand(describe.Cmd())
	RootCmd.AddCommand(get.Cmd())
	RootCmd.AddCommand(query.Cmd())
}

func initialize() {
	if shared.ElasticsearchHost == "" {
		conf := config.ParseConfigFile()
		readContextFromConfig(*conf)
	}
}

func readContextFromConfig(conf config.Config) {
	if len(conf.Contexts) == 0 {
		fmt.Println("Error: No contexts defined in the configuration.")
		os.Exit(1)
	}

	var context string

	if shared.Context != "" {
		context = shared.Context
	} else if conf.CurrentContext != "" {
		context = conf.CurrentContext
	} else {
		context = conf.Contexts[0].Name
	}

	clusterFound := false
	for _, cluster := range conf.Contexts {
		if cluster.Name == context {
			shared.ElasticsearchProtocol = cluster.Protocol
			if shared.ElasticsearchProtocol == "" {
				shared.ElasticsearchProtocol = constants.DefaultElasticsearchProtocol
			}
			shared.ElasticsearchPort = cluster.Port
			if shared.ElasticsearchPort == 0 {
				shared.ElasticsearchPort = constants.DefaultElasticsearchPort
			}
			shared.ElasticsearchUsername = cluster.Username
			shared.ElasticsearchPassword = cluster.Password
			shared.ElasticsearchHost = cluster.Host
			if shared.ElasticsearchHost == "" {
				fmt.Println("Error: 'host' field is not specified in the configuration for the current cluster.")
				os.Exit(1)
			}
			clusterFound = true
			break
		}
	}

	if !clusterFound {
		fmt.Printf("Error: No cluster found with the name '%s' in the configuration.\n", conf.CurrentContext)
		os.Exit(1)
	}
}

func initProtocolFlag() {
	defaultProtocol := constants.DefaultElasticsearchProtocol
	defaultProtocolEnv := os.Getenv(constants.ElasticsearchProtocolEnvVar)
	if defaultProtocolEnv != "" {
		defaultProtocol = defaultProtocolEnv
	}
	RootCmd.PersistentFlags().StringVar(&shared.ElasticsearchProtocol, "protocol", defaultProtocol, "Elasticsearch protocol")
}

func initHostFlag() {
	defaultHost := os.Getenv(constants.ElasticsearchHostEnvVar)
	RootCmd.PersistentFlags().StringVar(&shared.ElasticsearchHost, "host", defaultHost, "Elasticsearch host")
}

func initPortFlag() {
	defaultPort := constants.DefaultElasticsearchPort
	defaultPortStr := os.Getenv(constants.ElasticsearchPortEnvVar)
	if defaultPortStr != "" {
		parsedPort, err := strconv.Atoi(defaultPortStr)
		if err != nil || parsedPort <= 0 {
			fmt.Printf("Invalid value for %s environment variable: %s\n", constants.ElasticsearchPortEnvVar, defaultPortStr)
			os.Exit(1)
		}
		defaultPort = parsedPort
	}
	RootCmd.PersistentFlags().IntVar(&shared.ElasticsearchPort, "port", defaultPort, "Elasticsearch port")
}

func initUsernameFlag() {
	defaultUsername := os.Getenv(constants.ElasticsearchUsernameEnvVar)
	RootCmd.PersistentFlags().StringVar(&shared.ElasticsearchUsername, "username", defaultUsername, "Elasticsearch username")
}

func initPasswordFlag() {
	defaultPassword := os.Getenv(constants.ElasticsearchPasswordEnvVar)
	RootCmd.PersistentFlags().StringVar(&shared.ElasticsearchPassword, "password", defaultPassword, "Elasticsearch password")
}
