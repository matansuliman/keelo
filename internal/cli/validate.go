package cli

import (
	"fmt"
	"os"

	"keelo/internal/config"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [config-file]",
	Short: "Validate a project configuration file",
	Long:  `Validates that the provided project configuration file exists and has a valid schema.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configPath := "project.yaml" // Default path
		if len(args) > 0 {
			configPath = args[0]
		}

		cfg, err := config.LoadProjectConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully validated project config: %s\n", cfg.Project)
		fmt.Printf("Loaded %d modules.\n", len(cfg.Modules))
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
