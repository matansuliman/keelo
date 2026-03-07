package modules

import (
	"fmt"
	"os"
	"path/filepath"

	"keelo/pkg/types"

	"gopkg.in/yaml.v3"
)

// Loader acts as the module registry and loader.
type Loader struct {
	modulesDir string
	cacheDir   string
	downloader *Downloader
}

// NewLoader creates a new module loader.
func NewLoader(modulesDir, cacheDir string) *Loader {
	return &Loader{
		modulesDir: modulesDir,
		cacheDir:   cacheDir,
		downloader: NewDownloader(cacheDir),
	}
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
			modulePath, err = l.downloader.Download(modNode.Source)
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
