package root

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/go-rfe/gpwd/internal/logging"
)

var (
	logLevel string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gpwd",
	Short: "Simple password store",
	Long: `gpwd implements both client and server side infrastructure
to securely store password and any sensitive information.
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func AddCommand(command *cobra.Command) {
	rootCmd.AddCommand(command)
}

func init() {
	cobra.OnInitialize(setLogLevel)

	// Global log level
	rootCmd.PersistentFlags().StringVar(&logLevel, "logLevel", "ERROR", "log level (DEBUG|INFO|WARNING|ERROR)")
}

func setLogLevel() {
	logging.LogLevel(logLevel)
}
