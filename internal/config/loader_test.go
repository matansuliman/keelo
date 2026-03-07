package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadProjectConfig(t *testing.T) {
	tempDir := t.TempDir()

	validConfigPath := filepath.Join(tempDir, "valid.yaml")
	validYAML := `
project: myapp
modules:
  - name: postgres
    values:
      POSTGRES_DB: mydb
`
	if err := os.WriteFile(validConfigPath, []byte(validYAML), 0644); err != nil {
		t.Fatalf("Failed to write valid config file: %v", err)
	}

	invalidConfigPath := filepath.Join(tempDir, "invalid.yaml")
	invalidYAML := `
project:
  - bad: format
`
	if err := os.WriteFile(invalidConfigPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("Failed to write invalid config file: %v", err)
	}

	missingProjectConfigPath := filepath.Join(tempDir, "missing_project.yaml")
	missingProjectYAML := `
modules:
  - name: postgres
`
	if err := os.WriteFile(missingProjectConfigPath, []byte(missingProjectYAML), 0644); err != nil {
		t.Fatalf("Failed to write missing project config file: %v", err)
	}

	tests := []struct {
		name        string
		path        string
		expectError bool
		check       func(*testing.T, error)
	}{
		{
			name:        "Valid Config",
			path:        validConfigPath,
			expectError: false,
		},
		{
			name:        "File Not Found",
			path:        filepath.Join(tempDir, "nonexistent.yaml"),
			expectError: true,
			check: func(t *testing.T, err error) {
				if err == nil || !contains(err.Error(), "project config file not found") {
					t.Errorf("Expected 'file not found' error, got: %v", err)
				}
			},
		},
		{
			name:        "Invalid YAML Structure",
			path:        invalidConfigPath,
			expectError: true,
			check: func(t *testing.T, err error) {
				if err == nil || !contains(err.Error(), "invalid YAML") {
					t.Errorf("Expected 'invalid YAML' error, got: %v", err)
				}
			},
		},
		{
			name:        "Missing Project Name",
			path:        missingProjectConfigPath,
			expectError: true,
			check: func(t *testing.T, err error) {
				if err == nil || !contains(err.Error(), "missing 'project' name") {
					t.Errorf("Expected 'missing project name' error, got: %v", err)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := LoadProjectConfig(tc.path)
			if tc.expectError && err == nil {
				t.Fatalf("Expected error for %s, got nil", tc.name)
			}
			if !tc.expectError && err != nil {
				t.Fatalf("Expected no error for %s, got %v", tc.name, err)
			}
			if tc.check != nil {
				tc.check(t, err)
			}
		})
	}
}

func contains(s, substr string) bool {
	// Simple string contains helper to avoid importing strings just for this
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
