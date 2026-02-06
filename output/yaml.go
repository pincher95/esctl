package output

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func PrintYaml(data any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// Convert panics from yaml.Marshal into a regular error.
			err = fmt.Errorf("failed to marshal data to YAML: %v", r)
		}
	}()

	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to YAML: %w", err)
	}

	fmt.Println(string(yamlData))
	return nil
}
