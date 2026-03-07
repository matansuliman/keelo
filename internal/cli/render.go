package cli

import (
	"fmt"

	"keelo/internal/compose"
	"keelo/internal/config"
	"keelo/internal/merger"
	"keelo/internal/modules"
	"keelo/internal/renderer"
	"keelo/internal/validator"
	"keelo/pkg/types"

	"github.com/spf13/cobra"
)

var renderOutput string
var renderConfigPath string

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render the project configuration into a Docker Compose file",
	Long:  `Loads the project configuration, resolves modules, validates inputs, and renders a final docker-compose.generated.yaml file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadProjectConfig(renderConfigPath)
		if err != nil {
			return fmt.Errorf("loading project config: %w", err)
		}

		// Load Modules
		loader := modules.NewLoader("modules", ".keelo/cache")
		loadedModules, err := loader.LoadProjectModules(cfg)
		if err != nil {
			return fmt.Errorf("loading modules: %w", err)
		}

		// Validation and Rendering
		var fragments []*types.RenderedModule
		for i, modNode := range cfg.Modules {
			def, ok := loadedModules[modNode.Name]
			if !ok {
				return fmt.Errorf("module '%s' was not loaded", modNode.Name)
			}

			if err := validator.ValidateModuleInputs(&cfg.Modules[i], def); err != nil {
				return fmt.Errorf("validation error in module '%s': %w", modNode.Name, err)
			}

			rendered, err := renderer.RenderModuleTemplate(cfg.Project, &cfg.Modules[i], def)
			if err != nil {
				return fmt.Errorf("template rendering error for module '%s': %w", modNode.Name, err)
			}
			fragments = append(fragments, rendered)
		}

		// Merging
		mergedOutput, err := merger.MergeComposeFragments(fragments, cfg.Mixins)
		if err != nil {
			return fmt.Errorf("merge error: %w", err)
		}

		// Writing Output
		writer := compose.NewOutputWriter(renderOutput)
		if err := writer.Write(mergedOutput); err != nil {
			return fmt.Errorf("output error: %w", err)
		}

		fmt.Printf("Successfully rendered project '%s' into %s\n", cfg.Project, renderOutput)
		return nil
	},
}

func init() {
	renderCmd.Flags().StringVarP(&renderOutput, "output", "o", compose.DefaultOutputFileName(), "Path to write the generated Compose file")
	renderCmd.Flags().StringVarP(&renderConfigPath, "config", "c", "project.yaml", "Path to the project configuration file")
	rootCmd.AddCommand(renderCmd)
}
