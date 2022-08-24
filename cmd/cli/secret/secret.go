package secret

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-rfe/gpwd/cmd/root"
)

// secretCmd represents the store command
var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage gpwd secrets",
	Long:  `This is a cli for secret manipulation on the agent side`,
}

func init() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	root.AddCommand(secretCmd)

	// A socket file for agent connection
	secretCmd.PersistentFlags().String("socketPath", home+"/.gpwd/gpwd.sock", "Agent socket path")
	cobra.CheckErr(viper.BindPFlag("socket_path", secretCmd.PersistentFlags().Lookup("socketPath")))

	// An agent response timeout
	secretCmd.PersistentFlags().Duration("timeout", 10*time.Second, "Agent timeout")
	cobra.CheckErr(viper.BindPFlag("timeout", secretCmd.PersistentFlags().Lookup("timeout")))

	// An agent certificate for auth
	secretCmd.Flags().String("certPath", home+"/.gpwd/agent.pem", "Agent TLS certificate PEM file")
	cobra.CheckErr(viper.BindPFlag("cert_path", secretCmd.Flags().Lookup("certPath")))
}
