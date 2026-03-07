package exec

import (
	"testing"
)

func TestNewDockerComposeRunner(t *testing.T) {
	runner := NewDockerComposeRunner()
	if runner == nil {
		t.Fatal("Expected runner to be created, got nil")
	}
}

// Full execution tests require docker compose to be installed and mock environments.
// For MVP, we verify the initialization logic.
