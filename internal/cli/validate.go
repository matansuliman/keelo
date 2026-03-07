package cli

import (
	"fmt"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := "project.yaml" // Default path
		if len(args) > 0 {
			configPath = args[0]
		}

		cfg, err := config.LoadProjectConfig(configPath)
		if err != nil {
			return fmt.Errorf("loading project config: %w", err)
		}

		// Assume modules are in "modules" folder by default (can be configurable later)
		loader := modules.NewLoader("modules", ".keelo/cache")
		loadedModules, err := loader.LoadProjectModules(cfg)
		if err != nil {
			return fmt.Errorf("loading modules: %w", err)
		}

		// Validate module inputs
		for i, modNode := range cfg.Modules {
			def, ok := loadedModules[modNode.Name]
			if !ok {
				return fmt.Errorf("module '%s' was not loaded", modNode.Name)
			}

			if err := validator.ValidateModuleInputs(&cfg.Modules[i], def); err != nil {
				return fmt.Errorf("validation error in module '%s': %w", modNode.Name, err)
			}
		}

		fmt.Printf("Successfully validated project config: %s\n", cfg.Project)
		fmt.Printf("Loaded and validated %d modules.\n", len(loadedModules))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
