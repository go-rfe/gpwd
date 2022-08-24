package secret

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-rfe/gpwd/internal/client/secrets"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete secret using gpwd agent",
	Long:  `cli connects to the agent and removes secret data by ID`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
		defer cancel()

		client, err := secrets.NewSecretsClient(ctx, viper.GetString("socket_path"))
		cobra.CheckErr(err)

		cobra.CheckErr(client.Delete(viper.GetString("delete_id")))
	},
}

func init() {
	secretCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().String("id", "", "Secret ID")
	cobra.CheckErr(viper.BindPFlag("delete_id", deleteCmd.Flags().Lookup("id")))
}
