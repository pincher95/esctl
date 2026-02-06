package version

import (
	"fmt"
	"os"
	"runtime"

	"github.com/pincher95/esctl/internal/buildinfo"
	"github.com/pincher95/esctl/output"
	"github.com/pincher95/esctl/shared"
	"github.com/spf13/cobra"
)

type info struct {
	Version    string `json:"version"`
	Commit     string `json:"commit"`
	Date       string `json:"date"`
	BuiltBy    string `json:"built_by"`
	GoVersion  string `json:"go_version"`
	Platform   string `json:"platform"`
	Executable string `json:"executable"`
	ConfigPath string `json:"config_path"`
}

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		exe, err := os.Executable()
		if err != nil {
			exe = "unknown"
		}

		configPath := defaultConfigPath()
		data := info{
			Version:    buildinfo.Version,
			Commit:     buildinfo.Commit,
			Date:       buildinfo.Date,
			BuiltBy:    buildinfo.BuiltBy,
			GoVersion:  runtime.Version(),
			Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
			Executable: exe,
			ConfigPath: configPath,
		}

		switch shared.OutputFormat {
		case "json":
			return output.PrintJson(data)
		case "yaml":
			return output.PrintYaml(data)
		default:
			fmt.Printf("Version:     %s\n", data.Version)
			fmt.Printf("Commit:      %s\n", data.Commit)
			fmt.Printf("Date:        %s\n", data.Date)
			fmt.Printf("Built By:    %s\n", data.BuiltBy)
			fmt.Printf("Go Version:  %s\n", data.GoVersion)
			fmt.Printf("Platform:    %s\n", data.Platform)
			fmt.Printf("Executable:  %s\n", data.Executable)
			fmt.Printf("Config Path: %s\n", data.ConfigPath)
		}

		return nil
	},
}

func defaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "unknown"
	}
	return fmt.Sprintf("%s/.config/esctl.yaml", home)
}
