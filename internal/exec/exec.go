package exec

import (
	"fmt"
	"os"
	"os/exec"
)

// Runner defines the interface for executing Docker Compose commands
type Runner interface {
	Up(composeFile string, detach bool) error
	Down(composeFile string) error
	Logs(composeFile string, follow bool) error
}

// DockerComposeRunner implements Runner by calling the host 'docker compose' command.
type DockerComposeRunner struct{}

// NewDockerComposeRunner creates a new executor.
func NewDockerComposeRunner() *DockerComposeRunner {
	return &DockerComposeRunner{}
}

func (r *DockerComposeRunner) runCmd(args ...string) error {
	cmd := exec.Command("docker", args...)

	// Stream output directly to the user's terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	return nil
}

// Up runs 'docker compose up'
func (r *DockerComposeRunner) Up(composeFile string, detach bool) error {
	args := []string{"compose", "-f", composeFile, "up"}
	if detach {
		args = append(args, "-d")
	}
	return r.runCmd(args...)
}

// Down runs 'docker compose down'
func (r *DockerComposeRunner) Down(composeFile string) error {
	args := []string{"compose", "-f", composeFile, "down"}
	return r.runCmd(args...)
}

// Logs runs 'docker compose logs'
func (r *DockerComposeRunner) Logs(composeFile string, follow bool) error {
	args := []string{"compose", "-f", composeFile, "logs"}
	if follow {
		args = append(args, "-f")
	}
	return r.runCmd(args...)
}
