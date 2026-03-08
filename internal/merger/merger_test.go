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

func TestMergeComposeFragments_NestedEnvVars(t *testing.T) {
	frag1 := &types.RenderedModule{
		ModuleName: "app1",
		YAML: []byte(`
services:
  myapp:
    environment:
      SHARED_ENV: "app1_value"
      APP1_ONLY: "true"
`),
	}

	frag2 := &types.RenderedModule{
		ModuleName: "app2",
		YAML: []byte(`
services:
  myapp:
    environment:
      SHARED_ENV: "app2_value"
      APP2_ONLY: "true"
`),
	}

	merged, err := MergeComposeFragments([]*types.RenderedModule{frag1, frag2}, nil)
	if err != nil {
		t.Fatalf("Expected successful merge, got error: %v", err)
	}

	output := string(merged)

	// Ensure both unique env vars exist
	if !strings.Contains(output, "APP1_ONLY: \"true\"") {
		t.Errorf("Missing APP1_ONLY env var in merged output:\n%s", output)
	}
	if !strings.Contains(output, "APP2_ONLY: \"true\"") {
		t.Errorf("Missing APP2_ONLY env var in merged output:\n%s", output)
	}

	// SHARED_ENV will be overwritten by the last fragment (frag2) due to map merging
	if !strings.Contains(output, "SHARED_ENV: app2_value") {
		t.Errorf("Expected SHARED_ENV to be 'app2_value', got something else in merged output:\n%s", output)
	}
}

func TestMergeComposeFragments_StrategicMerge(t *testing.T) {
	frag1 := &types.RenderedModule{
		ModuleName: "api1",
		YAML: []byte(`
services:
  api:
    image: api:v1
    environment:
      VAR1: "true"
`),
	}

	frag2 := &types.RenderedModule{
		ModuleName: "api2",
		YAML: []byte(`
services:
  api:
    image: api:v2
    environment:
      VAR2: "true"
`),
	}

	merged, err := MergeComposeFragments([]*types.RenderedModule{frag1, frag2}, nil)
	if err != nil {
		t.Fatalf("Expected successful merge, got error: %v", err)
	}

	output := string(merged)

	// Ensure both contributing env vars exist (Nested Merging)
	if !strings.Contains(output, "VAR1: \"true\"") {
		t.Errorf("Missing VAR1 in merged output:\n%s", output)
	}
	if !strings.Contains(output, "VAR2: \"true\"") {
		t.Errorf("Missing VAR2 in merged output:\n%s", output)
	}

	// Scalar field (image) should be overwritten by the last fragment
	if !strings.Contains(output, "image: api:v2") {
		t.Errorf("Expected image api:v2, got something else in merged output:\n%s", output)
	}
}

func TestMergeComposeFragments_VolumeMerge(t *testing.T) {
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
    driver_opts:
      type: nfs
`),
	}

	merged, err := MergeComposeFragments([]*types.RenderedModule{frag1, frag2}, nil)
	if err != nil {
		t.Fatalf("Expected successful merge, got error: %v", err)
	}

	output := string(merged)
	if !strings.Contains(output, "driver: local") {
		t.Errorf("Missing driver in merged volumes:\n%s", output)
	}
	if !strings.Contains(output, "driver_opts:") {
		t.Errorf("Missing driver_opts in merged volumes:\n%s", output)
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

func TestMergeComposeFragments_CommentPreservation(t *testing.T) {
	frag := &types.RenderedModule{
		ModuleName: "postgres",
		YAML: []byte("# Global service comment\n" +
			"services:\n" +
			"  # The database\n" +
			"  app-postgres:\n" +
			"    image: postgres:15 # Use pg15 strictly\n"),
	}

	merged, err := MergeComposeFragments([]*types.RenderedModule{frag}, nil)
	if err != nil {
		t.Fatalf("Expected successful merge, got error: %v", err)
	}

	output := string(merged)
	if !strings.Contains(output, "# Global service comment") {
		t.Errorf("Missing global comment in merged output:\n%s", output)
	}
	if !strings.Contains(output, "# The database") {
		t.Errorf("Missing service comment in merged output:\n%s", output)
	}
	if !strings.Contains(output, "# Use pg15 strictly") {
		t.Errorf("Missing inline line comment in merged output:\n%s", output)
	}
}
