package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper to run root command and return output
func executeCommand(args ...string) (string, error) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String(), err
}

func TestCLI_Validate_E2E(t *testing.T) {
	// Point to our testdata fixtures
	cwd, _ := os.Getwd()

	// We need to change the modules directory loaded by the validate command.
	// We'll run the cmd from the testdata directory.
	testdataDir := filepath.Join(cwd, "..", "..", "testdata")
	os.Chdir(testdataDir)
	defer os.Chdir(cwd)

	output, err := executeCommand("validate", "project.yaml")
	if err != nil {
		t.Fatalf("Validate command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(output, "Successfully validated project config: e2e-test-project") {
		t.Errorf("Unexpected validate output: %s", output)
	}
}

func TestCLI_Render_E2E(t *testing.T) {
	cwd, _ := os.Getwd()
	testdataDir := filepath.Join(cwd, "..", "..", "testdata")
	outPath := filepath.Join(testdataDir, "docker-compose.test-render.yaml")
	defer os.Remove(outPath) // cleanup

	os.Chdir(testdataDir)
	defer os.Chdir(cwd)

	output, err := executeCommand("render", "--config", "project.yaml", "--output", "docker-compose.test-render.yaml")
	if err != nil {
		t.Fatalf("Render command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(output, "Successfully rendered project 'e2e-test-project'") {
		t.Errorf("Unexpected render output: %s", output)
	}

	// Verify the produced file
	data, err := os.ReadFile("docker-compose.test-render.yaml")
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "image: my-app:v1.2.3") {
		t.Errorf("Rendered template missing inputs: %s", content)
	}
	if !strings.Contains(content, "replicas: 2") {
		t.Errorf("Rendered template missing replicas: %s", content)
	}
}

func TestCLI_Init_ListModules(t *testing.T) {
	tempDir := t.TempDir()
	cwd, _ := os.Getwd()

	os.Chdir(tempDir)
	defer os.Chdir(cwd)

	// Test init
	output, err := executeCommand("init", "test-init")
	if err != nil {
		t.Fatalf("Init command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(output, "Successfully initialized new project 'test-init'") {
		t.Errorf("Unexpected init output: %s", output)
	}

	if _, err := os.Stat("project.yaml"); os.IsNotExist(err) {
		t.Fatalf("project.yaml was not created")
	}

	// Make a mock module to test list
	os.MkdirAll("modules/mock-mod", 0755)
	os.WriteFile("modules/mock-mod/module.yaml", []byte(`
name: mock-mod
description: A mock module
version: 1.0.0
`), 0644)

	// Test list-modules
	output, err = executeCommand("list-modules")
	if err != nil {
		t.Fatalf("List-modules command failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(output, "mock-mod") || !strings.Contains(output, "A mock module") {
		t.Errorf("Unexpected list-modules output: %s", output)
	}
}
