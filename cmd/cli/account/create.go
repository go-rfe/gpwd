package account

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-rfe/gpwd/internal/client/accounts"
	"github.com/go-rfe/gpwd/internal/encryption"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create account using gpwd agent",
	Long:  `cli connects to the agent and sends account info securely.`,
	Run: func(cmd *cobra.Command, args []string) {
		password, err := encryption.AskForSecretInput("Please enter user password:")
		cobra.CheckErr(err)

		ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
		defer cancel()

		client, err := accounts.NewAccountsClient(ctx, viper.GetString("socket_path"))
		cobra.CheckErr(err)

		id, err := client.Create(
			viper.GetString("create_server_address"),
			viper.GetString("create_username"),
			password,
		)
		cobra.CheckErr(err)
		fmt.Println(id)
	},
}

func init() {
	accountCmd.AddCommand(createCmd)

	createCmd.Flags().String("serverAddress", "", "Server address to connect to")
	cobra.CheckErr(viper.BindPFlag("create_server_address", createCmd.Flags().Lookup("serverAddress")))

	createCmd.Flags().String("username", "", "User name")
	cobra.CheckErr(viper.BindPFlag("create_username", createCmd.Flags().Lookup("username")))
}
