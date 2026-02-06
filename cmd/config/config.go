package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Modify and view the configuration",
}

// Flags for the add-context command
var (
	contextName     string
	contextHost     string
	contextPort     int
	contextProtocol string
	contextUsername string
	contextPassword string
)

var addContextCmd = &cobra.Command{
	Use:   "add-context",
	Short: "Add a new context to the configuration",
	Long:  `Add a new named context with connection details (e.g., host, port, username, password) to esctl.yml`,
	RunE:  runAddContext,
}

var updateContextCmd = &cobra.Command{
	Use:   "update-context",
	Short: "Update an existing context",
	Long:  `Update an existing context with new connection details (e.g., host, port, username, password) in esctl.yml`,
	RunE:  runUpdateContext,
}

var deleteContextCmd = &cobra.Command{
	Use:   "delete-context",
	Short: `Delete a context`,
	Long:  `Delete an existing context from the configuration`,
	RunE:  runDeleteContext,
}

var useContextCmd = &cobra.Command{
	Use:   "use-context",
	Short: "Set the current context",
	Long:  `Set the current context to connect to. This command updates the 'current-context' field in the configuration file.`,
	RunE:  runUseContext,
}

var getContextsCmd = &cobra.Command{
	Use:   "get-contexts",
	Short: "List the contexts defined in the esctl.yml file",
	RunE:  runGetContexts,
}

var currentContextCmd = &cobra.Command{
	Use:   "current-context",
	Short: "Display the current context",
	RunE:  runCurrentContext,
}

func init() {
	configCmd.AddCommand(useContextCmd)
	configCmd.AddCommand(getContextsCmd)
	configCmd.AddCommand(currentContextCmd)
	configCmd.AddCommand(addContextCmd)
	configCmd.AddCommand(updateContextCmd)
	configCmd.AddCommand(deleteContextCmd)

	// Define the flags for `use-context`
	useContextCmd.Flags().StringVarP(&contextName, "name", "c", "", "Name of the context to use (required)")
	// Mark name as required
	_ = useContextCmd.MarkFlagRequired("name")

	// Define the flags for `add-context`
	addContextCmd.Flags().StringVarP(&contextName, "name", "n", "", "Name of the new context (required)")
	addContextCmd.Flags().StringVarP(&contextHost, "host", "", "", "Elasticsearch host, e.g. example.com")
	addContextCmd.Flags().IntVar(&contextPort, "port", 9200, "Elasticsearch port (default: 9200)")
	addContextCmd.Flags().StringVar(&contextProtocol, "protocol", "http", "Protocol, e.g. http or https")
	addContextCmd.Flags().StringVarP(&contextUsername, "username", "u", "", "Username for Elasticsearch (if needed)")
	addContextCmd.Flags().StringVarP(&contextPassword, "password", "p", "", "Password for Elasticsearch (if needed)")
	// Mark name as required
	_ = addContextCmd.MarkFlagRequired("name")
	_ = addContextCmd.MarkFlagRequired("host")

	// Define the flags for `update-context`
	updateContextCmd.Flags().StringVarP(&contextName, "name", "n", "", "Name of the context to update (required)")
	updateContextCmd.Flags().StringVarP(&contextHost, "host", "", "", "Elasticsearch host")
	updateContextCmd.Flags().IntVar(&contextPort, "port", 9200, "Elasticsearch port")
	updateContextCmd.Flags().StringVar(&contextProtocol, "protocol", "http", "Protocol, e.g. http or https")
	updateContextCmd.Flags().StringVarP(&contextUsername, "username", "u", "", "Username for Elasticsearch")
	updateContextCmd.Flags().StringVarP(&contextPassword, "password", "p", "", "Password for Elasticsearch")
	// Mark name as required
	_ = updateContextCmd.MarkFlagRequired("name")

	// Define the flags for `delete-context`
	deleteContextCmd.Flags().StringVarP(&contextName, "name", "n", "", "Name of the context to delete (required)")
	// Mark name as required
	_ = deleteContextCmd.MarkFlagRequired("name")
}

func Cmd() *cobra.Command {
	return configCmd
}

// runAddContext is executed when someone calls `esctl config add-context --host=... --name=...`
func runAddContext(cmd *cobra.Command, args []string) error {
	// 1. Parse existing config
	config, err := ParseConfigFile()
	if err != nil {
		return err
	}

	// 2. Validate required flags
	if contextName == "" {
		return fmt.Errorf("--name is required")
	}

	contextExists := false
	for _, context := range config.Contexts {
		if context.Name == contextName {
			contextExists = true
			break
		}
	}

	if contextExists {
		return fmt.Errorf("context already exists with the name '%s' in the configuration", contextName)
	}

	// 3. Create a new Context from the flags
	newCtx := Context{
		Name:     contextName,
		Host:     contextHost,
		Port:     contextPort,
		Protocol: contextProtocol,
		Username: contextUsername,
		Password: contextPassword,
	}

	// 4. Append to existing contexts
	config.Contexts = append(config.Contexts, newCtx)

	// 5. Update Viper in-memory
	viper.Set("contexts", config.Contexts)

	// 6. Write the updated config to file
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("error writing updated configuration: %w", err)
	}

	// 7. Print success or new context
	fmt.Printf("Context %q added successfully.\n", contextName)
	return nil
}

func runUpdateContext(cmd *cobra.Command, args []string) error {
	config, err := ParseConfigFile()
	if err != nil {
		return err
	}

	if contextName == "" {
		return fmt.Errorf("--name is required")
	}

	contextExists := false
	for i, context := range config.Contexts {
		if context.Name == contextName {
			contextExists = true

			// Update the context with the new values
			if contextHost != "" {
				config.Contexts[i].Host = contextHost
			}
			if contextPort != 0 {
				config.Contexts[i].Port = contextPort
			}
			if contextProtocol != "" {
				config.Contexts[i].Protocol = contextProtocol
			}
			if contextUsername != "" {
				config.Contexts[i].Username = contextUsername
			}
			if contextPassword != "" {
				config.Contexts[i].Password = contextPassword
			}

			break
		}
	}

	if !contextExists {
		return fmt.Errorf("no context found with the name '%s' in the configuration", contextName)
	}

	// Update Viper in-memory
	viper.Set("contexts", config.Contexts)

	// Write the updated config to file
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("error writing updated configuration: %w", err)
	}

	fmt.Printf("Context %q updated successfully.\n", contextName)
	return nil
}

func runDeleteContext(cmd *cobra.Command, args []string) error {
	config, err := ParseConfigFile()
	if err != nil {
		return err
	}

	if contextName == "" {
		return fmt.Errorf("--name is required")
	}

	contextExists := false
	for i, context := range config.Contexts {
		if context.Name == contextName {
			contextExists = true
			config.Contexts = append(config.Contexts[:i], config.Contexts[i+1:]...)
			break
		}
	}

	if !contextExists {
		return fmt.Errorf("no context found with the name '%s' in the configuration", contextName)
	}

	viper.Set("contexts", config.Contexts)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("error writing updated configuration: %w", err)
	}

	fmt.Printf("Context %q deleted successfully.\n", contextName)
	return nil
}

func runUseContext(cmd *cobra.Command, args []string) error {
	config, err := ParseConfigFile()
	if err != nil {
		return err
	}

	if contextName == "" {
		return fmt.Errorf("--name is required")
	}

	contextExists := false
	for _, context := range config.Contexts {
		if context.Name == contextName {
			contextExists = true
			break
		}
	}

	if !contextExists {
		return fmt.Errorf("no context found with the name '%s' in the configuration", contextName)
	}

	viper.Set("current-context", contextName)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("error writing updated configuration: %w", err)
	}
	return nil
}

func runGetContexts(cmd *cobra.Command, args []string) error {
	config, err := ParseConfigFile()
	if err != nil {
		return err
	}
	for _, context := range config.Contexts {
		contextName := context.Name
		if contextName == config.CurrentContext {
			contextName += "(*)"
		}
		fmt.Printf("- name: %s\n", contextName)
		fmt.Printf("  host: %s\n", context.Host)
		if context.Protocol != "" {
			fmt.Printf("  protocol: %s\n", context.Protocol)
		}
		if context.Port != 0 {
			fmt.Printf("  port: %d\n", context.Port)
		}
		if context.Username != "" {
			fmt.Printf("  username: %s\n", context.Username)
		}
		if context.Password != "" {
			fmt.Printf("  password: %s\n", context.Password)
		}
	}
	return nil
}

func runCurrentContext(cmd *cobra.Command, args []string) error {
	config, err := ParseConfigFile()
	if err != nil {
		return err
	}
	fmt.Println(config.CurrentContext)
	return nil
}

type Context struct {
	Name     string `mapstructure:"name"`
	Protocol string `mapstructure:"protocol"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Entity struct {
	Columns []string `mapstructure:"columns"`
}

type Config struct {
	CurrentContext string            `mapstructure:"current-context"`
	Contexts       []Context         `mapstructure:"contexts"`
	Entities       map[string]Entity `mapstructure:"entities"`
}

func ParseConfigFile() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting user's home directory: %w", err)
	}

	viper.AddConfigPath(filepath.Join(home, ".config"))
	viper.SetConfigName("esctl")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config into struct: %w", err)
	}

	return &config, nil
}
