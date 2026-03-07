package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new project configuration interactively",
	Long:  `Creates a new project.yaml file in the current directory, asking for project details and a starting template.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat("project.yaml"); err == nil {
			return fmt.Errorf("project.yaml already exists in this directory")
		}

		// 1. Determine Project Name
		defaultName := "my-project"
		if len(args) > 0 {
			defaultName = args[0]
		} else {
			// Try to use current directory name
			if cwd, err := os.Getwd(); err == nil {
				defaultName = filepath.Base(cwd)
			}
		}

		projectName := defaultName
		selectedTemplate := "Basic Web Service"

		nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
		if !nonInteractive {
			var err error
			projectName, err = PromptString("Project Name", defaultName, ValidateNotEmpty)
			if err != nil {
				return fmt.Errorf("prompt failed: %w", err)
			}

			// 2. Select Template
			templates := []string{
				"Empty Template",
				"Basic Web Service",
				"Full Stack (Web + Postgres)",
			}

			selectedTemplate, err = PromptSelect("Starting Template", templates)
			if err != nil {
				return fmt.Errorf("prompt failed: %w", err)
			}
		}

		// 3. Generate YAML based on selection
		var content string
		switch selectedTemplate {
		case "Empty Template":
			content = fmt.Sprintf(`project: %s
modules: []
`, projectName)
		case "Basic Web Service":
			content = fmt.Sprintf(`project: %s
modules:
  - name: web
    values:
      PORT: "8080"
`, projectName)
		case "Full Stack (Web + Postgres)":
			content = fmt.Sprintf(`project: %s
modules:
  - name: postgres
    values:
      POSTGRES_DB: appdb
      POSTGRES_PASSWORD: secret
  - name: web
    values:
      PORT: "8080"
      DB_HOST: app-postgres
`, projectName)
		}

		if err := os.WriteFile("project.yaml", []byte(content), 0644); err != nil {
			return fmt.Errorf("creating project.yaml: %w", err)
		}

		fmt.Printf("✨ Successfully initialized new project '%s' using template '%s'!\n", projectName, selectedTemplate)
		return nil
	},
}

func init() {
	initCmd.Flags().Bool("non-interactive", false, "Skip interactive prompts and use defaults")
	rootCmd.AddCommand(initCmd)
}
