package agent

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-rfe/gpwd/cmd/root"
	"github.com/go-rfe/gpwd/internal/agent"
	"github.com/go-rfe/gpwd/internal/encryption"
	"github.com/go-rfe/gpwd/internal/logging/log"
)

const (
	minimumPasswordLength = 32
)

// agentCmd represents the agent command for starting gpwd local agent
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "gpwd passwords store agent",
	Long:  `Simple passwords store agent with SQLite backend`,
	Run: func(cmd *cobra.Command, args []string) {
		var password []byte
		var err error
		if viper.GetString("master_password") == "" {
			password, err = encryption.AskForSecretInput("Please type your master password:")
			if err != nil {
				log.Fatal().Msgf("Failed to read master password: %s", err)
			}
			if len(string(password)) < minimumPasswordLength {
				log.Fatal().Msgf("minimum master password length is %d characters", minimumPasswordLength)
			}
		} else {
			password = []byte(viper.GetString("master_password"))
		}
		viper.Set("master_password", password)

		config := agent.Cfg{}

		if err := viper.Unmarshal(&config); err != nil {
			log.Fatal().Msgf("Failed to read agent config: %s", err)
		}

		log.Debug().Msgf("Using socket: %s", config.SocketPath)
		agent.NewAgent(&config).Run()
	},
}

const (
	defaultSyncInterval = 10 * time.Second
)

var (
	cfgFile string
)

func init() {
	cobra.OnInitialize(InitConfig)

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	root.AddCommand(agentCmd)

	agentCmd.Flags().String("socketPath", home+"/.gpwd/gpwd.sock", "Agent socket path")
	cobra.CheckErr(viper.BindPFlag("socket_path", agentCmd.Flags().Lookup("socketPath")))

	agentCmd.Flags().Duration("syncInterval", defaultSyncInterval, "Time interval for sync data to server")
	cobra.CheckErr(viper.BindPFlag("sync_interval", agentCmd.Flags().Lookup("syncInterval")))

	agentCmd.Flags().String("storePath", home+"/.gpwd/", "Agent storage path")
	cobra.CheckErr(viper.BindPFlag("store_path", agentCmd.Flags().Lookup("storePath")))

	agentCmd.Flags().String("certPath", home+"/.gpwd/agent.pem", "Agent TLS certificate PEM file")
	cobra.CheckErr(viper.BindPFlag("agent_cert_path", agentCmd.Flags().Lookup("certPath")))

	agentCmd.Flags().String("keyPath", home+"/.gpwd/agent-key.pem", "Agent TLS key PEM file")
	cobra.CheckErr(viper.BindPFlag("agent_key_path", agentCmd.Flags().Lookup("keyPath")))

	agentCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.gpwd.yaml)")
}

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gpwd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gpwd")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info().Msgf("Using config file: %s", viper.ConfigFileUsed())
	}
}
