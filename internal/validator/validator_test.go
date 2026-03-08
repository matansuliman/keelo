package validator

import (
	"testing"

	"keelo/pkg/types"
)

func TestValidateModuleInputs(t *testing.T) {
	moduleDef := &types.ModuleDefinition{
		Name: "test-module",
		Inputs: map[string]types.ModuleInput{
			"REQ_STR": {
				Type:     "string",
				Required: true,
			},
			"OPT_DEF": {
				Type:     "int",
				Required: false,
				Default:  8080,
			},
			"UNSUPPORTED": {
				Type:     "list",
				Required: false,
			},
			"BOOL_FLAG": {
				Type:     "bool",
				Required: false,
			},
		},
	}

	tests := []struct {
		name        string
		inputVals   map[string]interface{}
		expectError bool
		check       func(*testing.T, error, map[string]interface{})
	}{
		{
			name: "Valid required and optional default applied",
			inputVals: map[string]interface{}{
				"REQ_STR": "my-value",
			},
			expectError: false,
			check: func(t *testing.T, err error, vals map[string]interface{}) {
				if vals["OPT_DEF"] != 8080 {
					t.Errorf("Expected default value 8080 to be applied, got %v", vals["OPT_DEF"])
				}
			},
		},
		{
			name:        "Missing required value",
			inputVals:   map[string]interface{}{},
			expectError: true,
			check: func(t *testing.T, err error, vals map[string]interface{}) {
				if err == nil || !contains(err.Error(), "missing required input") {
					t.Errorf("Expected 'missing required input' error, got: %v", err)
				}
			},
		},
		{
			name: "Invalid type",
			inputVals: map[string]interface{}{
				"REQ_STR": 123, // Expected string
			},
			expectError: true,
			check: func(t *testing.T, err error, vals map[string]interface{}) {
				if err == nil || !contains(err.Error(), "invalid type") {
					t.Errorf("Expected 'invalid type' error, got: %v", err)
				}
			},
		},
		{
			name: "Float resolving to int",
			inputVals: map[string]interface{}{
				"REQ_STR": "my-value",
				"OPT_DEF": float64(42.0), // YAML unmarshals sometimes read bare ints as float64
			},
			expectError: false,
		},
		{
			name: "Unsupported type passes through",
			inputVals: map[string]interface{}{
				"REQ_STR":     "my-value",
				"UNSUPPORTED": []string{"hello", "world"},
			},
			expectError: false,
		},
		{
			name: "Valid bool flag",
			inputVals: map[string]interface{}{
				"REQ_STR":   "my-value",
				"BOOL_FLAG": true,
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := &types.ModuleNode{
				Name:   "test-module",
				Values: tc.inputVals,
			}

			err := ValidateModuleInputs(node, moduleDef)
			if tc.expectError && err == nil {
				t.Fatalf("Expected error for %s, got nil", tc.name)
			}
			if !tc.expectError && err != nil {
				t.Fatalf("Expected no error for %s, got %v", tc.name, err)
			}
			if tc.check != nil {
				tc.check(t, err, node.Values)
			}
		})
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
