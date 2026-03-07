package merger

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"keelo/pkg/types"
)

// MergeComposeFragments takes multiple rendered modules and combines them into a single YAML structure.
// It detects conflicts in service names, volumes, and networks.
func MergeComposeFragments(fragments []*types.RenderedModule) ([]byte, error) {
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
