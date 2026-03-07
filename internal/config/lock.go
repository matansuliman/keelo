package config

import (
	"fmt"
	"os"

	"keelo/pkg/types"

	"gopkg.in/yaml.v3"
)

// SaveLockFile writes the LockFile to a yaml file.
func SaveLockFile(path string, lock *types.LockFile) error {
	data, err := yaml.Marshal(lock)
	if err != nil {
		return fmt.Errorf("failed to marshal lock file: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// LoadLockFile reads the LockFile from a yaml file.
func LoadLockFile(path string) (*types.LockFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Not an error if file doesn't exist yet
		}
		return nil, fmt.Errorf("failed to read lock file: %w", err)
	}

	var lock types.LockFile
	if err := yaml.Unmarshal(data, &lock); err != nil {
		return nil, fmt.Errorf("invalid YAML in lock file: %w", err)
	}

	return &lock, nil
}
