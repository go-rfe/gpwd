package server

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-rfe/gpwd/cmd/root"
	"github.com/go-rfe/gpwd/internal/logging/log"
	"github.com/go-rfe/gpwd/internal/server"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "gpwd passwords store server",
	Long:  `Simple passwords store server`,
	Run: func(cmd *cobra.Command, args []string) {
		config := server.Cfg{}

		if err := viper.Unmarshal(&config); err != nil {
			log.Fatal().Msgf("Failed to read agent config: %s", err)
		}

		server.NewServer(&config).Run()
	},
}

var (
	cfgFile string
)

const (
	defaultTokenLifeSpan = 3600 * time.Second
)

func init() {
	cobra.OnInitialize(InitConfig)

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	root.AddCommand(serverCmd)

	serverCmd.Flags().String("serverAddress", "localhost:8080", "Server listen address")
	cobra.CheckErr(viper.BindPFlag("server_address", serverCmd.Flags().Lookup("serverAddress")))

	serverCmd.Flags().String("certPath", home+"/.gpwd/server.pem", "Server TLS certificate PEM file")
	cobra.CheckErr(viper.BindPFlag("server_cert_path", serverCmd.Flags().Lookup("certPath")))

	serverCmd.Flags().String("keyPath", home+"/.gpwd/server-key.pem", "Server TLS key PEM file")
	cobra.CheckErr(viper.BindPFlag("server_key_path", serverCmd.Flags().Lookup("keyPath")))

	serverCmd.Flags().String("databaseDSN", "", "Server database DSN")
	cobra.CheckErr(viper.BindPFlag("database_dsn", serverCmd.Flags().Lookup("databaseDSN")))

	serverCmd.Flags().Duration("tokenLifespan", defaultTokenLifeSpan, "Server token lifespan")
	cobra.CheckErr(viper.BindPFlag("token_lifespan", serverCmd.Flags().Lookup("tokenLifespan")))

	serverCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.gpwd.yaml)")
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
