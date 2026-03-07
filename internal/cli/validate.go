package cli

import (
	"fmt"
	"os"

	"keelo/internal/config"
	"keelo/internal/modules"
	"keelo/internal/validator"

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
			fmt.Fprintf(os.Stderr, "Error loading project config: %v\n", err)
			os.Exit(1)
		}

		// Assume modules are in "modules" folder by default (can be configurable later)
		loader := modules.NewLoader("modules")
		loadedModules, err := loader.LoadProjectModules(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading modules: %v\n", err)
			os.Exit(1)
		}

		// Validate module inputs
		for i, modNode := range cfg.Modules {
			def, ok := loadedModules[modNode.Name]
			if !ok {
				fmt.Fprintf(os.Stderr, "Error: Module '%s' was not loaded\n", modNode.Name)
				os.Exit(1)
			}

			if err := validator.ValidateModuleInputs(&cfg.Modules[i], def); err != nil {
				fmt.Fprintf(os.Stderr, "Validation error: %v\n", err)
				os.Exit(1)
			}
		}

		fmt.Printf("Successfully validated project config: %s\n", cfg.Project)
		fmt.Printf("Loaded and validated %d modules.\n", len(loadedModules))
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
