package modules

import (
	"testing"
)

func TestResolveSource(t *testing.T) {
	tests := []struct {
		name     string
		modName  string
		source   string
		expected string
	}{
		{
			name:     "No source, standard name",
			modName:  "postgres",
			source:   "",
			expected: "",
		},
		{
			name:     "No source, shorthand github name",
			modName:  "matansuliman/postgres",
			source:   "",
			expected: "github.com/matansuliman/postgres",
		},
		{
			name:     "Explicit github URL",
			modName:  "postgres",
			source:   "github.com/matansuliman/keelo-modules//postgres",
			expected: "github.com/matansuliman/keelo-modules//postgres",
		},
		{
			name:     "Explicit git prefixed URL",
			modName:  "postgres",
			source:   "git::https://github.com/matansuliman/keelo-modules//postgres",
			expected: "git::https://github.com/matansuliman/keelo-modules//postgres",
		},
		{
			name:     "Shorthand github source",
			modName:  "db",
			source:   "matansuliman/keelo-modules//postgres",
			expected: "github.com/matansuliman/keelo-modules//postgres",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := resolveSource(tc.modName, tc.source)
			if actual != tc.expected {
				t.Errorf("Expected resolveSource('%s', '%s') to be '%s', but got '%s'", tc.modName, tc.source, tc.expected, actual)
			}
		})
	}
}
