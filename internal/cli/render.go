package cli

import (
	"fmt"
	"os"

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
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadProjectConfig(renderConfigPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading project config: %v\n", err)
			os.Exit(1)
		}

		// Load Modules
		loader := modules.NewLoader("modules", ".keelo/cache")
		loadedModules, err := loader.LoadProjectModules(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading modules: %v\n", err)
			os.Exit(1)
		}

		// Validation and Rendering
		var fragments []*types.RenderedModule
		for i, modNode := range cfg.Modules {
			def, ok := loadedModules[modNode.Name]
			if !ok {
				fmt.Fprintf(os.Stderr, "Error: Module '%s' was not loaded\n", modNode.Name)
				os.Exit(1)
			}

			if err := validator.ValidateModuleInputs(&cfg.Modules[i], def); err != nil {
				fmt.Fprintf(os.Stderr, "Validation error in module '%s': %v\n", modNode.Name, err)
				os.Exit(1)
			}

			rendered, err := renderer.RenderModuleTemplate(cfg.Project, &cfg.Modules[i], def)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Template rendering error for module '%s': %v\n", modNode.Name, err)
				os.Exit(1)
			}
			fragments = append(fragments, rendered)
		}

		// Merging
		mergedOutput, err := merger.MergeComposeFragments(fragments, cfg.Mixins)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Merge error: %v\n", err)
			os.Exit(1)
		}

		// Writing Output
		writer := compose.NewOutputWriter(renderOutput)
		if err := writer.Write(mergedOutput); err != nil {
			fmt.Fprintf(os.Stderr, "Output error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully rendered project '%s' into %s\n", cfg.Project, renderOutput)
	},
}

func init() {
	renderCmd.Flags().StringVarP(&renderOutput, "output", "o", compose.DefaultOutputFileName(), "Path to write the generated Compose file")
	renderCmd.Flags().StringVarP(&renderConfigPath, "config", "c", "project.yaml", "Path to the project configuration file")
	rootCmd.AddCommand(renderCmd)
}
