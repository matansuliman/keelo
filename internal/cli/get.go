package cli

import (
	"fmt"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, _ := cmd.Flags().GetString("config")

		cfg, err := config.LoadProjectConfig(configPath)
		if err != nil {
			return fmt.Errorf("loading project config: %w", err)
		}

		fmt.Printf("Fetching remote modules for project '%s'...\n", cfg.Project)

		loader := modules.NewLoader("modules", ".keelo/cache", false)
		loadedModules, err := loader.LoadProjectModules(cfg)
		if err != nil {
			return fmt.Errorf("fetching modules: %w", err)
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
				return fmt.Errorf("saving keelo.lock: %w", err)
			}
			fmt.Println("Generated keelo.lock")
		}

		fmt.Println("Successfully fetched all modules.")
		return nil
	},
}

func init() {
	getCmd.Flags().StringP("config", "c", "project.yaml", "Path to project configuration file")
	rootCmd.AddCommand(getCmd)
}
