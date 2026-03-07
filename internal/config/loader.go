package config

import (
	"fmt"
	"os"

	"keelo/pkg/types"

	"gopkg.in/yaml.v3"
)

// LoadProjectConfig loads a project configuration from the given file path.
func LoadProjectConfig(path string) (*types.ProjectConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("project config file not found: %s", path)
		}
		return nil, fmt.Errorf("failed to read project config (%s): %w", path, err)
	}

	var cfg types.ProjectConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid YAML in project config (%s): %w", path, err)
	}

	if cfg.Project == "" {
		return nil, fmt.Errorf("project config missing 'project' name: %s", path)
	}

	return &cfg, nil
}
