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
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(upComposeFile); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: Generated compose file '%s' not found. Run 'tool render' first.\n", upComposeFile)
			os.Exit(1)
		}

		runner := exec.NewDockerComposeRunner()
		fmt.Printf("Starting project via docker compose (file: %s)...\n", upComposeFile)

		if err := runner.Up(upComposeFile, upDetach); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting project: %v\n", err)
			os.Exit(1)
		}
	},
}

var downComposeFile string

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove the generated Docker Compose project",
	Long:  `Runs 'docker compose down' using the generated compose file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(downComposeFile); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: Generated compose file '%s' not found.\n", downComposeFile)
			os.Exit(1)
		}

		runner := exec.NewDockerComposeRunner()
		fmt.Printf("Stopping project via docker compose (file: %s)...\n", downComposeFile)

		if err := runner.Down(downComposeFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error stopping project: %v\n", err)
			os.Exit(1)
		}
	},
}

var logsFollow bool
var logsComposeFile string

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View logs from the generated Docker Compose project",
	Long:  `Runs 'docker compose logs' using the generated compose file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(logsComposeFile); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: Generated compose file '%s' not found.\n", logsComposeFile)
			os.Exit(1)
		}

		runner := exec.NewDockerComposeRunner()
		if err := runner.Logs(logsComposeFile, logsFollow); err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching logs: %v\n", err)
			os.Exit(1)
		}
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
