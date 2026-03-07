package merger

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"keelo/pkg/types"
)

// MergeComposeFragments takes multiple rendered modules and combines them into a single YAML structure.
// It detects conflicts in service names, volumes, and networks, and injects global mixins.
func MergeComposeFragments(fragments []*types.RenderedModule, mixins *types.ProjectMixins) ([]byte, error) {
	finalCompose := make(map[string]interface{})
	finalServices := make(map[string]interface{})
	finalVolumes := make(map[string]interface{})
	finalNetworks := make(map[string]interface{})

	for _, frag := range fragments {
		var parsed map[string]interface{}
		if err := yaml.Unmarshal(frag.YAML, &parsed); err != nil {
			return nil, fmt.Errorf("failed to parse YAML for module '%s': %w", frag.ModuleName, err)
		}

		// Merge Services
		if services, ok := parsed["services"].(map[string]interface{}); ok {
			for name, def := range services {
				if _, exists := finalServices[name]; exists {
					return nil, fmt.Errorf("conflict detected: service '%s' is defined multiple times (seen in module '%s')", name, frag.ModuleName)
				}
				finalServices[name] = def
			}
		}

		// Merge Volumes
		if volumes, ok := parsed["volumes"].(map[string]interface{}); ok {
			for name, def := range volumes {
				if existingOpts, exists := finalVolumes[name]; exists {
					// We only check for exact collision of names. In a real system we might merge volume options.
					// For MVP, allow declaring the same volume name across modules if they don't have conflicting configurations.
					// Actually, simplest is to just allow it or error. The specs say "detect name conflicts". Let's error for services and warn/merge for volumes, or just error for all to be safe.
					// Let's error for strictness in MVP unless it's identical nil.
					if existingOpts != nil || def != nil {
						return nil, fmt.Errorf("conflict detected: volume '%s' is defined multiple times with options (seen in module '%s')", name, frag.ModuleName)
					}
				}
				finalVolumes[name] = def
			}
		}

		// Merge Networks
		if networks, ok := parsed["networks"].(map[string]interface{}); ok {
			for name, def := range networks {
				if existingOpts, exists := finalNetworks[name]; exists {
					if existingOpts != nil || def != nil {
						return nil, fmt.Errorf("conflict detected: network '%s' is defined multiple times with options (seen in module '%s')", name, frag.ModuleName)
					}
				}
				finalNetworks[name] = def
			}
		}

		// (Optional) Merge other top-level keys if needed, but MVP focuses on services, volumes, networks.
	}

	// Phase 13: Inject Global Mixins
	if mixins != nil {
		for svcName, svcDef := range finalServices {
			svcMap, ok := svcDef.(map[string]interface{})
			if !ok {
				continue // Skip if a service isn't a proper map
			}

			// Apply Environment Mixins
			if len(mixins.Environment) > 0 {
				existingEnv, exists := svcMap["environment"].(map[string]interface{})
				if !exists {
					existingEnv = make(map[string]interface{})
					svcMap["environment"] = existingEnv
				}
				for k, v := range mixins.Environment {
					// We only inject if it doesn't already exist natively in the module,
					// or we can allow mixins to override. Overriding makes mixins more powerful.
					existingEnv[k] = v
				}
			}

			// Apply Labels Mixins
			if len(mixins.Labels) > 0 {
				existingLabels, exists := svcMap["labels"].(map[string]interface{})
				if !exists {
					existingLabels = make(map[string]interface{})
					svcMap["labels"] = existingLabels
				}
				for k, v := range mixins.Labels {
					existingLabels[k] = v
				}
			}
			finalServices[svcName] = svcMap
		}
	}

	if len(finalServices) > 0 {
		finalCompose["services"] = finalServices
	}
	if len(finalVolumes) > 0 {
		finalCompose["volumes"] = finalVolumes
	}
	if len(finalNetworks) > 0 {
		finalCompose["networks"] = finalNetworks
	}

	// Deterministic ordering: we rely on yaml.Marshal which sorts map keys for map[string]interface{}.
	// According to gopkg.in/yaml.v3, map keys are serialized in alphabetical order.
	// This satisfies the "preserve deterministic ordering" requirement for the output.
	output, err := yaml.Marshal(finalCompose)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal merged compose file: %w", err)
	}

	return output, nil
}
