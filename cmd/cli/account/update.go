package account

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-rfe/gpwd/internal/client/accounts"
	"github.com/go-rfe/gpwd/internal/encryption"
	"github.com/go-rfe/gpwd/internal/logging/log"
)

// updateCmd represents the create command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update account using gpwd agent",
	Long:  `cli connects to the agent and updates account`,
	Run: func(cmd *cobra.Command, args []string) {
		password, err := encryption.AskForSecretInput("Please enter user password:")
		cobra.CheckErr(err)

		ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
		defer cancel()

		switch {
		case viper.GetString("update_server_address") == "":
			log.Fatal().Msg("Server address is required")
		case viper.GetString("update_username") == "":
			log.Fatal().Msg("Username is required")
		}

		client, err := accounts.NewAccountsClient(ctx, viper.GetString("socket_path"))
		cobra.CheckErr(err)

		cobra.CheckErr(client.Update(
			viper.GetString("update_server_address"),
			viper.GetString("update_username"),
			password,
		))
	},
}

func init() {
	accountCmd.AddCommand(updateCmd)

	updateCmd.Flags().String("serverAddress", "", "Server address to connect to")
	cobra.CheckErr(viper.BindPFlag("update_server_address", updateCmd.Flags().Lookup("serverAddress")))

	updateCmd.Flags().String("username", "", "User name")
	cobra.CheckErr(viper.BindPFlag("update_username", updateCmd.Flags().Lookup("username")))
}
