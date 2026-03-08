package modules

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"keelo/internal/config"
	"keelo/pkg/types"

	"gopkg.in/yaml.v3"
)

// Loader acts as the module registry and loader.
type Loader struct {
	modulesDir   string
	cacheDir     string
	downloader   *Downloader
	forceRefresh bool
}

// DefaultCacheDir returns the global cache directory for Keelo modules.
func DefaultCacheDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return ".keelo/cache" // Fallback to local if UserCacheDir fails
	}
	return filepath.Join(dir, "keelo", "modules")
}

// NewLoader creates a new module loader.
func NewLoader(modulesDir, cacheDir string, forceRefresh bool) *Loader {
	return &Loader{
		modulesDir:   modulesDir,
		cacheDir:     cacheDir,
		downloader:   NewDownloader(cacheDir),
		forceRefresh: forceRefresh,
	}
}

// HashDirectory calculates a robust SHA256 checksum of all files within a directory block.
func HashDirectory(dirPath string) (string, error) {
	hash := sha256.New()

	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(dirPath, path)
			// Skip files that shouldn't affect the module's core integrity (e.g., .git inner contents)
			// But for cached modules, we might just hash everything except .git
			if filepath.Base(relPath) == ".git" {
				return filepath.SkipDir
			}
			files = append(files, relPath)
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	sort.Strings(files)

	for _, relPath := range files {
		fullPath := filepath.Join(dirPath, relPath)
		f, err := os.Open(fullPath)
		if err != nil {
			return "", err
		}
		// Write the filepath relative to the module root to the hash (prevents rename exploits)
		hash.Write([]byte(relPath))

		// Write the file contents into the hash
		if _, err := io.Copy(hash, f); err != nil {
			f.Close()
			return "", err
		}
		f.Close()
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// LoadModule reads and parses the module.yaml file for a given module name from the local modules directory.
func (l *Loader) LoadModule(name string) (*types.ModuleDefinition, error) {
	modulePath := filepath.Join(l.modulesDir, name)
	return l.loadDefinitionFromPath(name, modulePath)
}

// loadDefinitionFromPath reads and parses the module.yaml file from a specific path.
func (l *Loader) loadDefinitionFromPath(name, modulePath string) (*types.ModuleDefinition, error) {
	yamlPath := filepath.Join(modulePath, "module.yaml")

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("module '%s' not found at %s", name, yamlPath)
		}
		return nil, fmt.Errorf("failed to read module.yaml for '%s': %w", name, err)
	}

	var def types.ModuleDefinition
	if err := yaml.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("invalid YAML in module definition for '%s': %w", name, err)
	}

	if def.Name == "" {
		def.Name = name // fallback
	}

	def.Subpath = modulePath
	return &def, nil
}

// resolveSource expands a shorthand name or source into a canonical go-getter URL.
func resolveSource(name, source string) string {
	if source == "" {
		if strings.Contains(name, "/") {
			return "github.com/" + name
		}
		return ""
	}

	// Expand shorthand sources like "matansuliman/postgres"
	if !strings.Contains(source, "://") && !strings.HasPrefix(source, "git::") && !strings.HasPrefix(source, "github.com/") {
		if strings.Contains(source, "/") {
			return "github.com/" + source
		}
	}
	return source
}

// LoadProjectModules loads all module definitions referenced by a project configuration, pulling remote ones if needed.
func (l *Loader) LoadProjectModules(cfg *types.ProjectConfig) (map[string]*types.ModuleDefinition, error) {
	loaded := make(map[string]*types.ModuleDefinition)
	lock, _ := config.LoadLockFile("keelo.lock")

	for _, modNode := range cfg.Modules {
		if err := l.resolveAndLoad(modNode.Name, modNode.Source, loaded, lock); err != nil {
			return nil, err
		}
	}

	return loaded, nil
}

// resolveAndLoad is a recursive helper to load a module and its dependencies.
func (l *Loader) resolveAndLoad(name string, source string, loaded map[string]*types.ModuleDefinition, lock *types.LockFile) error {
	if _, exists := loaded[name]; exists {
		return nil // already loaded
	}

	actualSource := resolveSource(name, source)
	var modulePath string

	if actualSource != "" {
		// Remote module logic
		var err error
		modulePath, err = l.downloader.Download(actualSource, l.forceRefresh)
		if err != nil {
			return fmt.Errorf("failed to load remote module '%s': %w", name, err)
		}

		if lock != nil && !l.forceRefresh {
			if err := verifyChecksum(modulePath, name, lock); err != nil {
				return err
			}
		}
	} else {
		// Local module logic
		modulePath = filepath.Join(l.modulesDir, name)
	}

	def, err := l.loadDefinitionFromPath(name, modulePath)
	if err != nil {
		return err
	}
	loaded[name] = def

	// Load dependencies recursively (they can now be remote shorthands)
	for _, dep := range def.Dependencies {
		if err := l.resolveAndLoad(dep, "", loaded, lock); err != nil {
			return fmt.Errorf("failed to load dependency '%s' for module '%s': %w", dep, name, err)
		}
	}

	return nil
}

// verifyChecksum checks if a downloaded module matches the checksum in the lock file.
func verifyChecksum(modulePath, name string, lock *types.LockFile) error {
	for _, lockedMod := range lock.Modules {
		if lockedMod.Name == name && lockedMod.Checksum != "" {
			actualHash, hashErr := HashDirectory(modulePath)
			if hashErr != nil {
				return fmt.Errorf("failed to hash cached module for validation: %w", hashErr)
			}
			if actualHash != lockedMod.Checksum {
				return fmt.Errorf("TAMPERING DETECTED: cached module '%s' checksum does not match keelo.lock. Run with --force-refresh to overwrite cache", name)
			}
			break
		}
	}
	return nil
}
