package modules

import (
	"os"
	"path/filepath"
	"testing"
	"time"

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

func TestLoader_LoadRemoteProjectModules(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping remote load test in short mode")
	}

	cacheDir := t.TempDir()

	// Ensure cleanup
	defer func() {
		os.RemoveAll(cacheDir)
	}()

	loader := NewLoader(t.TempDir(), cacheDir, false)

	cfg := &types.ProjectConfig{
		Project: "test-remote-proj",
		Modules: []types.ModuleNode{
			{
				Name:   "remote-postgres",
				Source: "git::https://github.com/matansuliman/keelo-modules//base-postgres",
			},
		},
	}

	// Loading process should fetch the module and cache it
	t.Logf("Attempting to load remote modules into cache dir: %s", cacheDir)

	start := time.Now()
	loaded, err := loader.LoadProjectModules(cfg)
	t.Logf("Finished loading remote modules. Took: %v", time.Since(start))

	if err != nil {
		t.Fatalf("Failed to load remote project modules: %v", err)
	}

	if len(loaded) != 1 {
		t.Fatalf("Expected 1 remote module loaded, got %d", len(loaded))
	}

	remoteMod, ok := loaded["remote-postgres"]
	if !ok {
		t.Fatalf("Expected 'remote-postgres' to be loaded")
	}

	// Verify loaded properties
	if remoteMod.Name != "base-postgres" {
		t.Errorf("Expected base-postgres schema name, got '%s'", remoteMod.Name)
	}

	if input, ok := remoteMod.Inputs["DB_NAME"]; !ok {
		t.Errorf("Expected remote DB_NAME input, but it was missing")
	} else if !input.Required {
		t.Errorf("Expected remote DB_NAME to be required")
	}

	// Verify cache directory was populated
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		t.Fatalf("Failed to read cache dir: %v", err)
	}

	if len(entries) == 0 {
		t.Errorf("Expected cache dir %s to contain downloaded module, but it was empty", cacheDir)
	}

	foundRemote := false
	for _, e := range entries {
		if len(e.Name()) == 12 { // It's a 12-character hash directory
			foundRemote = true

			// Verify module files exist inside the cached path
			// go-getter extracts the subdirectory contents directly into the target dst dir
			modYamlPath := filepath.Join(cacheDir, e.Name(), "module.yaml")
			if _, err := os.Stat(modYamlPath); os.IsNotExist(err) {
				t.Errorf("Cached remote module.yaml not found at expected path: %s", modYamlPath)
			}
			tmplPath := filepath.Join(cacheDir, e.Name(), "compose.yaml.tmpl")
			if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
				t.Errorf("Cached remote template not found at expected path: %s", tmplPath)
			}
			break
		}
	}

	if !foundRemote {
		t.Errorf("Expected to find a hash-named folder in cache dir, contents: %v", entries)
	}
}
