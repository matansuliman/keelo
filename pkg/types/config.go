package types

// ProjectConfig represents the root structure of a project.yaml file.
type ProjectConfig struct {
	Project string         `yaml:"project"`
	Mixins  *ProjectMixins `yaml:"mixins,omitempty"`
	Modules []ModuleNode   `yaml:"modules"`
}

// ProjectMixins defines global properties to be injected into all services.
type ProjectMixins struct {
	Labels      map[string]string `yaml:"labels,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
}

// ModuleNode represents a single module declaration within the project config.
type ModuleNode struct {
	Name   string                 `yaml:"name"`
	Source string                 `yaml:"source,omitempty"`
	Values map[string]interface{} `yaml:"values"`
}
