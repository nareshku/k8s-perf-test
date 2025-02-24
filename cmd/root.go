package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	configFile string
	duration   time.Duration
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-perf-test",
	Short: "Kubernetes API Performance Testing Tool",
	Long: `A performance testing tool for Kubernetes API server that executes 
concurrent API calls across different resources and users.

Example:
  k8s-perf-test --config config.yaml --duration 5m`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPerformanceTest(configFile, duration)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "config.yaml", "path to config file")
	rootCmd.PersistentFlags().DurationVar(&duration, "duration", 5*time.Minute, "duration to run the test (e.g., 1h, 5m, 30s)")

	// Mark config flag as required
	rootCmd.MarkPersistentFlagRequired("config")
}
