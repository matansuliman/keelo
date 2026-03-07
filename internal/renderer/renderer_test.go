package renderer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"keelo/pkg/types"
)

func TestRenderModuleTemplate(t *testing.T) {
	tempDir := t.TempDir()

	// Mock module definition
	modDef := &types.ModuleDefinition{
		Name:    "postgres",
		Subpath: tempDir,
	}

	// Mock template file
	tmplContent := `
services:
  postgres:
    image: postgres:15
    container_name: {{ .ProjectName }}-postgres
    environment:
      POSTGRES_DB: {{ .Values.POSTGRES_DB }}
`
	tmplPath := filepath.Join(tempDir, "compose.yaml.tmpl")
	if err := os.WriteFile(tmplPath, []byte(tmplContent), 0644); err != nil {
		t.Fatalf("Failed to write mock template: %v", err)
	}

	// Mock project node
	modNode := &types.ModuleNode{
		Name: "postgres",
		Values: map[string]interface{}{
			"POSTGRES_DB": "testdb",
		},
	}

	// Test Render
	rendered, err := RenderModuleTemplate("myproject", modNode, modDef)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	if rendered.ModuleName != "postgres" {
		t.Errorf("Expected ModuleName 'postgres', got '%s'", rendered.ModuleName)
	}

	output := string(rendered.YAML)

	if !strings.Contains(output, "container_name: myproject-postgres") {
		t.Errorf("Template did not correctly render ProjectName. Output: \n%s", output)
	}

	if !strings.Contains(output, "POSTGRES_DB: testdb") {
		t.Errorf("Template did not correctly render POSTGRES_DB. Output: \n%s", output)
	}
}

func TestRenderModuleTemplate_MissingTemplate(t *testing.T) {
	modDef := &types.ModuleDefinition{
		Name:    "missing-tmpl",
		Subpath: t.TempDir(), // Empty dir
	}
	modNode := &types.ModuleNode{Name: "missing-tmpl"}

	_, err := RenderModuleTemplate("testproj", modNode, modDef)
	if err == nil {
		t.Fatal("Expected error for missing template, got nil")
	}
	if !strings.Contains(err.Error(), "module template not found") {
		t.Errorf("Expected 'module template not found' error, got: %v", err)
	}
}
