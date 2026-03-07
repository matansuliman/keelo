package types

// ModuleDefinition represents the parsed contents of a module.yaml file.
type ModuleDefinition struct {
	Name         string                 `yaml:"name"`
	Version      string                 `yaml:"version"`
	Description  string                 `yaml:"description"`
	Inputs       map[string]ModuleInput `yaml:"inputs"`
	Dependencies []string               `yaml:"dependencies"`

	// Subpath stores the local directory path where the module was found.
	Subpath string `yaml:"-"`
}

// ModuleInput defines the schema for a single configuration value required by a module.
type ModuleInput struct {
	Type        string      `yaml:"type"`
	Default     interface{} `yaml:"default"`
	Required    bool        `yaml:"required"`
	Sensitive   bool        `yaml:"sensitive"`
	Description string      `yaml:"description"`
}
