package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tool",
	Short: "A modular composition layer for Docker Compose",
	Long: `A CLI tool that acts as a package manager and renderer for Compose applications.
It allows you to declare modules for infrastructure components and generates a final docker-compose.yml.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
