package secret

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-rfe/gpwd/internal/client/secrets"
	"github.com/go-rfe/gpwd/internal/encryption"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create secret using gpwd agent",
	Long: `cli connects to the agent and sends secret data securely.
Data could be loaded from file or provided inline`,
	Run: func(cmd *cobra.Command, args []string) {
		var data []byte
		var err error

		dataFilePath := viper.GetString("create_data_file_path")
		if dataFilePath != "" {
			file, err := os.ReadFile(dataFilePath)
			cobra.CheckErr(err)

			data = file
		} else {
			data, err = encryption.AskForSecretInput("Please enter secret data inline:")
			cobra.CheckErr(err)
		}

		labels := viper.GetStringSlice("create_labels")

		ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
		defer cancel()

		client, err := secrets.NewSecretsClient(ctx, viper.GetString("socket_path"))
		cobra.CheckErr(err)

		id, err := client.Create(data, labels)
		cobra.CheckErr(err)
		fmt.Println(id)
	},
}

func init() {
	secretCmd.AddCommand(createCmd)

	createCmd.Flags().String("dataFromFile", "", "File path to get data from")
	cobra.CheckErr(viper.BindPFlag("create_data_file_path", createCmd.Flags().Lookup("dataFromFile")))

	createCmd.Flags().StringSlice("labels", nil, "Labels key=value, pairs")
	cobra.CheckErr(viper.BindPFlag("create_labels", createCmd.Flags().Lookup("labels")))
}
