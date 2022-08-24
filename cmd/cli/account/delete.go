package account

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-rfe/gpwd/internal/client/accounts"
)

// deleteCmd represents the create command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete account using gpwd agent",
	Long:  "cli connects to the agent and removes account",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
		defer cancel()

		client, err := accounts.NewAccountsClient(ctx, viper.GetString("socket_path"))
		cobra.CheckErr(err)

		cobra.CheckErr(client.Delete())
	},
}

func init() {
	accountCmd.AddCommand(deleteCmd)
}
