package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new project configuration",
	Long:  `Creates a new project.yaml file in the current directory with basic boilerplate structure.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := "my-project"
		if len(args) > 0 {
			projectName = args[0]
		}

		if _, err := os.Stat("project.yaml"); err == nil {
			return fmt.Errorf("project.yaml already exists in this directory")
		}

		boilerplate := fmt.Sprintf(`project: %s
modules:
  # Example module declaration
  # - name: postgres
  #   values:
  #     POSTGRES_DB: mydb
  #     POSTGRES_PASSWORD: secret
`, projectName)

		if err := os.WriteFile("project.yaml", []byte(boilerplate), 0644); err != nil {
			return fmt.Errorf("creating project.yaml: %w", err)
		}

		fmt.Printf("Successfully initialized new project '%s' in project.yaml\n", projectName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
