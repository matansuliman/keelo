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

		forceRefresh, _ := cmd.Flags().GetBool("force-refresh")
		loader := modules.NewLoader("modules", modules.DefaultCacheDir(), forceRefresh)
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

				// Calculate secure directory checksum
				checksum, err := modules.HashDirectory(def.Subpath)
				if err != nil {
					return fmt.Errorf("failed to hash module '%s': %w", modNode.Name, err)
				}

				lock.Modules = append(lock.Modules, types.LockedModule{
					Name:     modNode.Name,
					Source:   modNode.Source,
					Resolved: resolved,
					Checksum: checksum,
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
	getCmd.Flags().Bool("force-refresh", false, "Force re-download of remote modules, bypassing cache")
	rootCmd.AddCommand(getCmd)
}
