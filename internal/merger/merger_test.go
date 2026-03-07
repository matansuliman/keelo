package merger

import (
	"strings"
	"testing"

	"keelo/pkg/types"
)

func TestMergeComposeFragments(t *testing.T) {
	frag1 := &types.RenderedModule{
		ModuleName: "postgres",
		YAML: []byte(`
services:
  app-postgres:
    image: postgres:15
volumes:
  pgdata:
`),
	}

	frag2 := &types.RenderedModule{
		ModuleName: "redis",
		YAML: []byte(`
services:
  app-redis:
    image: redis:7
volumes:
  redisdata:
`),
	}

	merged, err := MergeComposeFragments([]*types.RenderedModule{frag1, frag2}, nil)
	if err != nil {
		t.Fatalf("Expected successful merge, got error: %v", err)
	}

	output := string(merged)
	if !strings.Contains(output, "app-postgres:") {
		t.Errorf("Missing app-postgres service in merged output")
	}
	if !strings.Contains(output, "app-redis:") {
		t.Errorf("Missing app-redis service in merged output")
	}
	if !strings.Contains(output, "pgdata:") {
		t.Errorf("Missing pgdata volume in merged output")
	}
	if !strings.Contains(output, "redisdata:") {
		t.Errorf("Missing redisdata volume in merged output")
	}
}

func TestMergeComposeFragments_ServiceConflict(t *testing.T) {
	frag1 := &types.RenderedModule{
		ModuleName: "api1",
		YAML: []byte(`
services:
  api:
    image: api:v1
`),
	}

	frag2 := &types.RenderedModule{
		ModuleName: "api2",
		YAML: []byte(`
services:
  api:
    image: api:v2
`),
	}

	_, err := MergeComposeFragments([]*types.RenderedModule{frag1, frag2}, nil)
	if err == nil {
		t.Fatalf("Expected error for service name conflict, got nil")
	}

	if !strings.Contains(err.Error(), "conflict detected: service 'api' is defined multiple times") {
		t.Errorf("Expected conflict error message, got: %v", err)
	}
}

func TestMergeComposeFragments_VolumeConflict(t *testing.T) {
	frag1 := &types.RenderedModule{
		ModuleName: "api1",
		YAML: []byte(`
volumes:
  shared-data:
    driver: local
`),
	}

	frag2 := &types.RenderedModule{
		ModuleName: "api2",
		YAML: []byte(`
volumes:
  shared-data:
    driver: other
`),
	}

	_, err := MergeComposeFragments([]*types.RenderedModule{frag1, frag2}, nil)
	if err == nil {
		t.Fatalf("Expected error for volume block conflict, got nil")
	}

	if !strings.Contains(err.Error(), "conflict detected: volume 'shared-data' is defined multiple times with options") {
		t.Errorf("Expected conflict error message, got: %v", err)
	}
}

func TestMergeComposeFragments_WithMixins(t *testing.T) {
	frag1 := &types.RenderedModule{
		ModuleName: "api1",
		YAML: []byte(`
services:
  api:
    image: api:v1
    environment:
      EXISTING: "true"
`),
	}

	mixins := &types.ProjectMixins{
		Labels: map[string]string{
			"env":        "production",
			"managed-by": "keelo",
		},
		Environment: map[string]string{
			"GLOBAL_VAR": "injected",
		},
	}

	merged, err := MergeComposeFragments([]*types.RenderedModule{frag1}, mixins)
	if err != nil {
		t.Fatalf("Expected successful merge, got error: %v", err)
	}

	output := string(merged)

	// Check Labels Injection
	if !strings.Contains(output, "env: production") || !strings.Contains(output, "managed-by: keelo") {
		t.Errorf("Mixins labels were not correctly injected into the output:\n%s", output)
	}

	// Check Env Injection and Existing Preservation
	if !strings.Contains(output, "GLOBAL_VAR: injected") {
		t.Errorf("Mixins environment was not correctly injected into the output:\n%s", output)
	}
	if !strings.Contains(output, "EXISTING: \"true\"") && !strings.Contains(output, "EXISTING: 'true'") && !strings.Contains(output, "EXISTING: true") {
		t.Errorf("Existing environment variable was lost during mixin injection:\n%s", output)
	}
}
