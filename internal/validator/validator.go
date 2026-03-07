package validator

import (
	"fmt"
	"reflect"

	"keelo/pkg/types"
)

// ValidateModuleInputs checks if the user-provided values in the project config
// satisfy the rules defined in the module's module.yaml.
func ValidateModuleInputs(moduleNode *types.ModuleNode, moduleDef *types.ModuleDefinition) error {
	providedValues := moduleNode.Values
	if providedValues == nil {
		providedValues = make(map[string]interface{})
		moduleNode.Values = providedValues // ensure it's not nil for later
	}

	for inputName, inputDef := range moduleDef.Inputs {
		val, exists := providedValues[inputName]

		if !exists {
			if inputDef.Required && inputDef.Default == nil {
				return fmt.Errorf("module '%s' is missing required input: %s", moduleDef.Name, inputName)
			}
			// Apply default if it exists and no value provided
			if inputDef.Default != nil {
				providedValues[inputName] = inputDef.Default
			}
			continue
		}

		// Type validation could happen here.
		// For MVP, we'll do best-effort loose type checking.
		if err := validateType(inputName, val, inputDef.Type); err != nil {
			return fmt.Errorf("module '%s': %w", moduleDef.Name, err)
		}
	}

	return nil
}

// validateType performs a very basic check against the expected type string natively from YAML.
func validateType(key string, val interface{}, expectedType string) error {
	if expectedType == "" {
		return nil // no type specified, anything goes
	}

	var valid bool
	switch expectedType {
	case "string":
		_, valid = val.(string)
	case "int":
		_, valid = val.(int)
		if !valid {
			// YAML parser might read it as float64 depending on format
			_, validFloat := val.(float64)
			valid = validFloat
		}
	case "bool":
		_, valid = val.(bool)
	default:
		// Unsupported type check in MVP, let it pass
		valid = true
	}

	if !valid {
		return fmt.Errorf("invalid type for '%s': expected %s, got %v", key, expectedType, reflect.TypeOf(val))
	}

	return nil
}
