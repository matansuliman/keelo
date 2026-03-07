package types

// ProjectConfig represents the root structure of a project.yaml file.
type ProjectConfig struct {
	Project string       `yaml:"project"`
	Modules []ModuleNode `yaml:"modules"`
}

// ModuleNode represents a single module declaration within the project config.
type ModuleNode struct {
	Name   string                 `yaml:"name"`
	Source string                 `yaml:"source,omitempty"`
	Values map[string]interface{} `yaml:"values"`
}
