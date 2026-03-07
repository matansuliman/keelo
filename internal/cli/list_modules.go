package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"keelo/internal/modules"

	"github.com/spf13/cobra"
)

var listModulesCmd = &cobra.Command{
	Use:   "list-modules",
	Short: "List all available local modules",
	Long:  `Scans the modules/ directory and lists the name, version, and description of all local modules found.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		modulesDir := "modules"

		if _, err := os.Stat(modulesDir); os.IsNotExist(err) {
			return fmt.Errorf("%s directory not found. Please run this from a directory containing modules", modulesDir)
		}

		entries, err := os.ReadDir(modulesDir)
		if err != nil {
			return fmt.Errorf("reading %s directory: %w", modulesDir, err)
		}

		loader := modules.NewLoader(modulesDir, ".keelo/cache", false)
		foundAny := false

		// Use tabwriter for nice alignment
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
		fmt.Fprintln(w, "NAME\tVERSION\tDESCRIPTION")

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			// verify module.yaml exists
			if _, err := os.Stat(filepath.Join(modulesDir, entry.Name(), "module.yaml")); os.IsNotExist(err) {
				continue
			}

			def, err := loader.LoadModule(entry.Name())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to load module %s: %v\n", entry.Name(), err)
				continue
			}

			foundAny = true
			version := def.Version
			if version == "" {
				version = "unk"
			}

			desc := def.Description
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}

			fmt.Fprintf(w, "%s\t%s\t%s\n", def.Name, version, desc)
		}

		if !foundAny {
			fmt.Println("No modules found.")
			return nil
		}

		w.Flush()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listModulesCmd)
}
