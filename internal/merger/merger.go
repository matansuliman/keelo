package merger

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"keelo/pkg/types"
)

// MergeComposeFragments takes multiple rendered modules and combines them into a single YAML structure.
// It detects conflicts in service names, volumes, and networks, and injects global mixins.
func MergeComposeFragments(fragments []*types.RenderedModule, mixins *types.ProjectMixins) ([]byte, error) {
	rootMap := &yaml.Node{Kind: yaml.MappingNode}

	for _, frag := range fragments {
		var docNode yaml.Node
		if err := yaml.Unmarshal(frag.YAML, &docNode); err != nil {
			return nil, fmt.Errorf("failed to parse YAML for module '%s': %w", frag.ModuleName, err)
		}

		if len(docNode.Content) == 0 || docNode.Content[0].Kind != yaml.MappingNode {
			continue // Skip empty or invalid
		}

		fragMap := docNode.Content[0]
		if err := mergeMappingNodes(rootMap, fragMap, frag.ModuleName, ""); err != nil {
			return nil, err
		}
	}

	// Phase 13: Inject Global Mixins
	if mixins != nil {
		applyMixinsToAST(rootMap, mixins)
	}

	// Because we are building the AST without an enclosing DocumentNode initially,
	// yaml.Marshal will automatically wrap our MappingNode into a DocumentNode.
	output, err := yaml.Marshal(rootMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal merged compose file: %w", err)
	}

	return output, nil
}

func mergeMappingNodes(dst *yaml.Node, src *yaml.Node, moduleName string, path string) error {
	for i := 0; i < len(src.Content); i += 2 {
		srcKeyNode := src.Content[i]
		srcValNode := src.Content[i+1]
		srcKey := srcKeyNode.Value

		var dstValNode *yaml.Node
		dstIdx := -1
		for j := 0; j < len(dst.Content); j += 2 {
			if dst.Content[j].Value == srcKey {
				dstValNode = dst.Content[j+1]
				dstIdx = j
				break
			}
		}

		if dstIdx == -1 {
			// Append both key and value nodes
			dst.Content = append(dst.Content, srcKeyNode, srcValNode)
		} else {
			// Key exists in dst
			// Check for explicit conflicts before recursing into children
			if path == "services" {
				return fmt.Errorf("conflict detected: service '%s' is defined multiple times (seen in module '%s')", srcKey, moduleName)
			}
			if path == "volumes" {
				if !isEmptyNode(dstValNode) || !isEmptyNode(srcValNode) {
					return fmt.Errorf("conflict detected: volume '%s' is defined multiple times with options (seen in module '%s')", srcKey, moduleName)
				}
			}
			if path == "networks" {
				if !isEmptyNode(dstValNode) || !isEmptyNode(srcValNode) {
					return fmt.Errorf("conflict detected: network '%s' is defined multiple times with options (seen in module '%s')", srcKey, moduleName)
				}
			}

			if dstValNode.Kind == yaml.MappingNode && srcValNode.Kind == yaml.MappingNode {
				newPath := srcKey
				if path != "" {
					newPath = path + "." + srcKey
				}
				if err := mergeMappingNodes(dstValNode, srcValNode, moduleName, newPath); err != nil {
					return err
				}
			} else {
				// For primitive values simply overwrite
				dst.Content[dstIdx+1] = srcValNode
			}
		}
	}
	return nil
}

func isEmptyNode(n *yaml.Node) bool {
	if n == nil {
		return true
	}
	if n.Tag == "!!null" || n.Value == "~" || n.Value == "null" {
		return true
	}
	if n.Kind == yaml.MappingNode && len(n.Content) == 0 {
		return true
	}
	return false
}

func applyMixinsToAST(root *yaml.Node, mixins *types.ProjectMixins) {
	var servicesNode *yaml.Node
	for i := 0; i < len(root.Content); i += 2 {
		if root.Content[i].Value == "services" {
			servicesNode = root.Content[i+1]
			break
		}
	}

	if servicesNode == nil || servicesNode.Kind != yaml.MappingNode {
		return
	}

	for i := 1; i < len(servicesNode.Content); i += 2 {
		svc := servicesNode.Content[i]
		if svc.Kind != yaml.MappingNode {
			continue
		}

		if len(mixins.Environment) > 0 {
			injectMap(svc, "environment", mixins.Environment)
		}
		if len(mixins.Labels) > 0 {
			injectMap(svc, "labels", mixins.Labels)
		}
	}
}

func injectMap(svcNode *yaml.Node, targetKey string, data map[string]string) {
	var targetMap *yaml.Node
	for i := 0; i < len(svcNode.Content); i += 2 {
		if svcNode.Content[i].Value == targetKey {
			targetMap = svcNode.Content[i+1]
			break
		}
	}

	if targetMap == nil {
		targetMap = &yaml.Node{Kind: yaml.MappingNode}
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: targetKey}
		svcNode.Content = append(svcNode.Content, keyNode, targetMap)
	}

	if targetMap.Kind != yaml.MappingNode {
		return
	}

	for k, v := range data {
		replaced := false
		for i := 0; i < len(targetMap.Content); i += 2 {
			if targetMap.Content[i].Value == k {
				targetMap.Content[i+1] = &yaml.Node{Kind: yaml.ScalarNode, Value: v}
				replaced = true
				break
			}
		}
		if !replaced {
			kNode := &yaml.Node{Kind: yaml.ScalarNode, Value: k}
			vNode := &yaml.Node{Kind: yaml.ScalarNode, Value: v}
			targetMap.Content = append(targetMap.Content, kNode, vNode)
		}
	}
}
