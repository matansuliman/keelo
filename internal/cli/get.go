package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"keelo/internal/config"
	"keelo/internal/modules"
	"keelo/pkg/types"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Download remote modules and update the cache",
	Long:  `Fetches all remote modules specified in the project configuration and stores them in the local cache.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")

		cfg, err := config.LoadProjectConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading project config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Fetching remote modules for project '%s'...\n", cfg.Project)

		loader := modules.NewLoader("modules", ".keelo/cache")
		loadedModules, err := loader.LoadProjectModules(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching modules: %v\n", err)
			os.Exit(1)
		}

		// Generate lock file
		lock := &types.LockFile{}
		for _, modNode := range cfg.Modules {
			if modNode.Source != "" {
				def, ok := loadedModules[modNode.Name]
				if !ok {
					continue
				}
				// The hash/resolved ID is the last part of the path in the cache
				resolved := filepath.Base(def.Subpath)
				lock.Modules = append(lock.Modules, types.LockedModule{
					Name:     modNode.Name,
					Source:   modNode.Source,
					Resolved: resolved,
				})
			}
		}

		if len(lock.Modules) > 0 {
			if err := config.SaveLockFile("keelo.lock", lock); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving keelo.lock: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Generated keelo.lock")
		}

		fmt.Println("Successfully fetched all modules.")
	},
}

func init() {
	getCmd.Flags().StringP("config", "c", "project.yaml", "Path to project configuration file")
	rootCmd.AddCommand(getCmd)
}
