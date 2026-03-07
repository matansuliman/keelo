package cli

import (
	"fmt"
	"os"

	"keelo/internal/compose"
	"keelo/internal/exec"

	"github.com/spf13/cobra"
)

var upDetach bool
var upComposeFile string

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the generated Docker Compose project",
	Long:  `Runs 'docker compose up' using the generated compose file. Defaults to docker-compose.generated.yaml.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(upComposeFile); os.IsNotExist(err) {
			return fmt.Errorf("generated compose file '%s' not found. Run 'keelo render' first", upComposeFile)
		}

		runner := exec.NewDockerComposeRunner()
		fmt.Printf("Starting project via docker compose (file: %s)...\n", upComposeFile)

		if err := runner.Up(upComposeFile, upDetach); err != nil {
			return fmt.Errorf("starting project: %w", err)
		}
		return nil
	},
}

var downComposeFile string

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove the generated Docker Compose project",
	Long:  `Runs 'docker compose down' using the generated compose file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(downComposeFile); os.IsNotExist(err) {
			return fmt.Errorf("generated compose file '%s' not found", downComposeFile)
		}

		runner := exec.NewDockerComposeRunner()
		fmt.Printf("Stopping project via docker compose (file: %s)...\n", downComposeFile)

		if err := runner.Down(downComposeFile); err != nil {
			return fmt.Errorf("stopping project: %w", err)
		}
		return nil
	},
}

var logsFollow bool
var logsComposeFile string

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View logs from the generated Docker Compose project",
	Long:  `Runs 'docker compose logs' using the generated compose file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(logsComposeFile); os.IsNotExist(err) {
			return fmt.Errorf("generated compose file '%s' not found", logsComposeFile)
		}

		runner := exec.NewDockerComposeRunner()
		if err := runner.Logs(logsComposeFile, logsFollow); err != nil {
			return fmt.Errorf("fetching logs: %w", err)
		}
		return nil
	},
}

func init() {
	defaultFile := compose.DefaultOutputFileName()

	upCmd.Flags().BoolVarP(&upDetach, "detach", "d", false, "Detached mode: Run containers in the background")
	upCmd.Flags().StringVarP(&upComposeFile, "file", "f", defaultFile, "Path to the generated compose file")

	downCmd.Flags().StringVarP(&downComposeFile, "file", "f", defaultFile, "Path to the generated compose file")

	logsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "Follow log output")
	logsCmd.Flags().StringVar(&logsComposeFile, "file", defaultFile, "Path to the generated compose file")

	rootCmd.AddCommand(upCmd, downCmd, logsCmd)
}
