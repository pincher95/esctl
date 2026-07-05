package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/pincher95/esctl/cmd/analyze"
	"github.com/pincher95/esctl/cmd/bulk"
	"github.com/pincher95/esctl/cmd/config"
	"github.com/pincher95/esctl/cmd/count"
	"github.com/pincher95/esctl/cmd/delete"
	"github.com/pincher95/esctl/cmd/describe"
	"github.com/pincher95/esctl/cmd/explain"
	"github.com/pincher95/esctl/cmd/get"
	"github.com/pincher95/esctl/cmd/profile"
	"github.com/pincher95/esctl/cmd/query"
	setcmd "github.com/pincher95/esctl/cmd/set"
	"github.com/pincher95/esctl/cmd/update"
	"github.com/pincher95/esctl/cmd/version"
	"github.com/pincher95/esctl/constants"
	"github.com/pincher95/esctl/internal/client"
	"github.com/pincher95/esctl/internal/logger"
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "esctl",
	Short: "esctl is CLI for Elasticsearch",
	Long:  `esctl is CLI for Elasticsearch that allows users to manage and monitor their Elasticsearch clusters.`,
	// main() is the single error reporter: it prints runtime errors to stderr and
	// exits quietly on user interruption. Silencing cobra avoids a duplicate
	// "Error: ..." line and a noisy usage dump on every runtime failure.
	SilenceErrors: true,
	SilenceUsage:  true,
}

type cancelContextKey struct{}

var portEnvErr error

func Execute(ctx context.Context) error {
	return RootCmd.ExecuteContext(ctx)
}

func init() {
	initProtocolFlag()
	initHostFlag()
	initPortFlag()
	initUsernameFlag()
	initPasswordFlag()

	RootCmd.PersistentFlags().StringVar(&shared.Context, "context", "", "Override context")
	RootCmd.PersistentFlags().StringVar(&shared.ElasticsearchAPIKey, "api-key", os.Getenv(constants.ElasticsearchAPIKeyEnvVar), "Elasticsearch API key (Authorization: ApiKey ...); takes precedence over username/password")
	RootCmd.PersistentFlags().StringVar(&shared.CACertPath, "ca-cert", os.Getenv(constants.ElasticsearchCACertEnvVar), "Path to a PEM CA certificate bundle to verify the server's TLS certificate")
	RootCmd.PersistentFlags().BoolVar(&shared.TLSInsecure, "insecure", false, "Skip TLS certificate verification (INSECURE; testing only — prefer --ca-cert)")
	RootCmd.PersistentFlags().BoolVar(&shared.Debug, "debug", false, "Enable debug mode")
	RootCmd.PersistentFlags().StringVarP(&shared.OutputFormat, "output", "o", "table", "Output format: table|json|yaml")
	RootCmd.PersistentFlags().DurationVar(&shared.TimeoutDuration, "timeout", 0, "Global timeout for command execution (e.g. 30s, 2m)")

	// Root level context timeout handling and initialization
	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if portEnvErr != nil {
			return portEnvErr
		}
		if shared.TimeoutDuration > 0 {
			ctx, cancel := context.WithTimeout(cmd.Context(), shared.TimeoutDuration)
			ctx = context.WithValue(ctx, cancelContextKey{}, cancel)
			cmd.SetContext(ctx)
		}
		return initialize()
	}
	RootCmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
		if cancel, ok := cmd.Context().Value(cancelContextKey{}).(context.CancelFunc); ok {
			cancel()
		}
	}

	RootCmd.AddCommand(bulk.Cmd())
	RootCmd.AddCommand(config.Cmd())
	RootCmd.AddCommand(count.Cmd())
	RootCmd.AddCommand(delete.Cmd())
	RootCmd.AddCommand(describe.Cmd())
	RootCmd.AddCommand(get.Cmd())
	RootCmd.AddCommand(query.Cmd())
	RootCmd.AddCommand(update.Cmd())
	RootCmd.AddCommand(setcmd.Cmd())
	RootCmd.AddCommand(analyze.Cmd)
	RootCmd.AddCommand(explain.Cmd)
	RootCmd.AddCommand(profile.Cmd)
	RootCmd.AddCommand(version.Cmd)
}

func initialize() error {
	// Initialize logger first
	logger.Init(shared.Debug)

	logger.Debug("initializing esctl",
		"debug", shared.Debug,
		"output", shared.OutputFormat,
		"timeout", shared.TimeoutDuration)

	if shared.ElasticsearchHost == "" {
		conf, err := config.ParseConfigFile()
		if err != nil {
			return err
		}
		if err := readContextFromConfig(conf); err != nil {
			return err
		}
	}

	if err := initClient(); err != nil {
		return err
	}

	logger.Debug("esctl initialized",
		"host", shared.ElasticsearchHost,
		"port", shared.ElasticsearchPort,
		"protocol", shared.ElasticsearchProtocol)
	return nil
}

func readContextFromConfig(conf *config.Config) error {
	if len(conf.Contexts) == 0 {
		return fmt.Errorf("no contexts defined in the configuration")
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
				return fmt.Errorf("'host' field is not specified in the configuration for the current cluster")
			}
			clusterFound = true
			break
		}
	}

	if !clusterFound {
		return fmt.Errorf("no cluster found with the name '%s' in the configuration", conf.CurrentContext)
	}
	return nil
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
			portEnvErr = fmt.Errorf("invalid value for %s environment variable: %s", constants.ElasticsearchPortEnvVar, defaultPortStr)
		} else {
			defaultPort = parsedPort
		}
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

func initClient() error {
	baseURL := fmt.Sprintf("%s://%s:%d", shared.ElasticsearchProtocol, shared.ElasticsearchHost, shared.ElasticsearchPort)

	cfg := &client.Config{
		BaseURL:     baseURL,
		Debug:       shared.Debug,
		Username:    shared.ElasticsearchUsername,
		Password:    shared.ElasticsearchPassword,
		APIKey:      shared.ElasticsearchAPIKey,
		CACertPath:  shared.CACertPath,
		TLSInsecure: shared.TLSInsecure,
		Timeout:     shared.TimeoutDuration,
	}

	cli, err := client.NewClient(cfg)
	if err != nil {
		return err
	}
	shared.SetClient(cli)
	return nil
}
