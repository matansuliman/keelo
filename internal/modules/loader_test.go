package modules

import (
	"os"
	"path/filepath"
	"testing"

	"keelo/pkg/types"
)

func TestLoader_LoadModule(t *testing.T) {
	tempDir := t.TempDir()

	// Create a mock module 'postgres'
	pgDir := filepath.Join(tempDir, "postgres")
	if err := os.Mkdir(pgDir, 0755); err != nil {
		t.Fatalf("Failed to create postgres dir: %v", err)
	}

	pgYaml := `
name: postgres
version: 1.0.0
description: PostgreSQL database module
inputs:
  POSTGRES_DB:
    type: string
    default: appdb
    required: true
  PORT:
    type: int
    default: 5432
`
	if err := os.WriteFile(filepath.Join(pgDir, "module.yaml"), []byte(pgYaml), 0644); err != nil {
		t.Fatalf("Failed to write postgres module.yaml: %v", err)
	}

	loader := NewLoader(tempDir, t.TempDir(), false)

	def, err := loader.LoadModule("postgres")
	if err != nil {
		t.Fatalf("Failed to load postgres module: %v", err)
	}

	if def.Name != "postgres" {
		t.Errorf("Expected module name 'postgres', got '%s'", def.Name)
	}
	if def.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", def.Version)
	}

	if input, ok := def.Inputs["POSTGRES_DB"]; !ok {
		t.Errorf("Expected input POSTGRES_DB, but it was missing")
	} else {
		if input.Default != "appdb" {
			t.Errorf("Expected POSTGRES_DB default 'appdb', got '%v'", input.Default)
		}
		if !input.Required {
			t.Errorf("Expected POSTGRES_DB to be required")
		}
	}
}

func TestLoader_LoadProjectModules(t *testing.T) {
	tempDir := t.TempDir()

	// Create mock module 'postgres'
	pgDir := filepath.Join(tempDir, "postgres")
	os.Mkdir(pgDir, 0755)
	pgYaml := `
name: postgres
`
	os.WriteFile(filepath.Join(pgDir, "module.yaml"), []byte(pgYaml), 0644)

	loader := NewLoader(tempDir, t.TempDir(), false)
	cfg := &types.ProjectConfig{
		Project: "test-proj",
		Modules: []types.ModuleNode{
			{Name: "postgres"},
		},
	}

	loaded, err := loader.LoadProjectModules(cfg)
	if err != nil {
		t.Fatalf("Failed to load project modules: %v", err)
	}

	if len(loaded) != 1 {
		t.Errorf("Expected 1 module loaded, got %d", len(loaded))
	}
	if _, ok := loaded["postgres"]; !ok {
		t.Errorf("Expected 'postgres' to be loaded")
	}
}
