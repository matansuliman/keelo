package renderer

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"keelo/pkg/types"
)

// TemplateContext is the data structure passed to the Go template engine.
type TemplateContext struct {
	ProjectName string
	Values      map[string]interface{}
}

// RenderModuleTemplate loads a module's compose.yaml.tmpl file and renders it.
func RenderModuleTemplate(projectName string, moduleNode *types.ModuleNode, moduleDef *types.ModuleDefinition) (*types.RenderedModule, error) {
	tmplPath := filepath.Join(moduleDef.Subpath, "compose.yaml.tmpl")

	tmplContent, err := os.ReadFile(tmplPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("module template not found at %s", tmplPath)
		}
		return nil, fmt.Errorf("failed to read template for module '%s': %w", moduleDef.Name, err)
	}

	tmpl, err := template.New(moduleDef.Name).Parse(string(tmplContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template for module '%s': %w", moduleDef.Name, err)
	}

	ctx := TemplateContext{
		ProjectName: projectName,
		Values:      moduleNode.Values,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil {
		return nil, fmt.Errorf("failed to execute template for module '%s': %w", moduleDef.Name, err)
	}

	return &types.RenderedModule{
		ModuleName: moduleDef.Name,
		YAML:       buf.Bytes(),
	}, nil
}
