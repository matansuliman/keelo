package modules

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

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

// LoadProjectModules loads all module definitions referenced by a project configuration, pulling remote ones if needed.
func (l *Loader) LoadProjectModules(cfg *types.ProjectConfig) (map[string]*types.ModuleDefinition, error) {
	loaded := make(map[string]*types.ModuleDefinition)

	for _, modNode := range cfg.Modules {
		// Avoid loading the same module twice
		if _, exists := loaded[modNode.Name]; exists {
			continue
		}

		var modulePath string
		if modNode.Source != "" {
			// Remote module: download first
			var err error
			modulePath, err = l.downloader.Download(modNode.Source, l.forceRefresh)
			if err != nil {
				return nil, fmt.Errorf("failed to load remote module '%s': %w", modNode.Name, err)
			}
		} else {
			// Local module: use default path
			modulePath = filepath.Join(l.modulesDir, modNode.Name)
		}

		def, err := l.loadDefinitionFromPath(modNode.Name, modulePath)
		if err != nil {
			return nil, err
		}
		loaded[modNode.Name] = def

		// Load dependencies (local only for now, as dependencies don't have 'source' in their module.yaml yet)
		// Future improvement: support remote dependencies in module.yaml
		for _, dep := range def.Dependencies {
			if _, exists := loaded[dep]; !exists {
				depDef, err := l.LoadModule(dep)
				if err != nil {
					return nil, fmt.Errorf("failed to load dependency '%s' for module '%s': %w", dep, def.Name, err)
				}
				loaded[dep] = depDef
			}
		}
	}

	return loaded, nil
}
