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
}

// NewLoader creates a new module loader.
func NewLoader(modulesDir string) *Loader {
	return &Loader{
		modulesDir: modulesDir,
	}
}

// LoadModule reads and parses the module.yaml file for a given module name.
func (l *Loader) LoadModule(name string) (*types.ModuleDefinition, error) {
	modulePath := filepath.Join(l.modulesDir, name)
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

// LoadProjectModules loads all module definitions referenced by a project configuration.
func (l *Loader) LoadProjectModules(cfg *types.ProjectConfig) (map[string]*types.ModuleDefinition, error) {
	loaded := make(map[string]*types.ModuleDefinition)

	for _, modNode := range cfg.Modules {
		// Avoid loading the same module twice
		if _, exists := loaded[modNode.Name]; exists {
			continue
		}

		def, err := l.LoadModule(modNode.Name)
		if err != nil {
			return nil, err
		}
		loaded[modNode.Name] = def

		// In later phases we can recursively load Dependencies.
		// For MVP, if it's explicitly in the project config it's enough,
		// but let's pre-load dependencies to be safe.
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
