package secret

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-rfe/gpwd/internal/client/secrets"
	"github.com/go-rfe/gpwd/internal/encryption"
)

// getCmd represents the create command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get secret from gpwd agent",
	Long:  `cli connects to the agent and gets data by ID`,
	Run: func(cmd *cobra.Command, args []string) {
		var password []byte

		ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
		defer cancel()

		client, err := secrets.NewSecretsClient(ctx, viper.GetString("socket_path"))
		cobra.CheckErr(err)

		secret, err := client.Get(viper.GetString("get_id"))
		cobra.CheckErr(err)

		if viper.GetString("master_password") == "" {
			password, err = encryption.AskForSecretInput("Please type your master password:")
			cobra.CheckErr(err)
		} else {
			password = []byte(viper.GetString("master_password"))
		}

		_, decrypt, err := encryption.GetCrypto(password)
		cobra.CheckErr(err)

		var labels []byte
		if secret.GetLabels() != nil {
			labels, err = json.Marshal(secret.GetLabels())
			cobra.CheckErr(err)
		}

		var decryptedData []byte

		if secret.GetData() != nil {
			decryptedData, err = decrypt(secret.GetData())
			cobra.CheckErr(err)
		}

		dataFilePath := viper.GetString("data_file_path")
		if dataFilePath != "" {
			err := os.WriteFile(dataFilePath, decryptedData, 0600)
			cobra.CheckErr(err)
		} else {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.Escape)
			_, err = fmt.Fprintln(w, "ID", "\t", "Labels", "\t", "Created At", "\t", "Updated At", "\t", "Deleted At", "\t", "Data")
			cobra.CheckErr(err)

			createdAt := ""
			if secret.CreatedAt != nil {
				createdAt = secret.CreatedAt.AsTime().String()
			}

			updatedAt := ""
			if secret.UpdatedAt != nil {
				updatedAt = secret.UpdatedAt.AsTime().String()
			}

			deletedAt := ""
			if secret.DeletedAt != nil {
				deletedAt = secret.DeletedAt.AsTime().String()
			}

			_, err = fmt.Fprintln(w, secret.GetID(), "\t", string(labels), "\t", createdAt, "\t", updatedAt, "\t", deletedAt, "\t", encryption.ToBase64(decryptedData))
			cobra.CheckErr(err)
			cobra.CheckErr(w.Flush())
		}
	},
}

func init() {
	secretCmd.AddCommand(getCmd)

	getCmd.Flags().String("id", "", "Secret ID")
	cobra.CheckErr(viper.BindPFlag("get_id", getCmd.Flags().Lookup("id")))

	getCmd.Flags().String("dataFilePath", "", "File path to save data to")
	cobra.CheckErr(viper.BindPFlag("data_file_path", getCmd.Flags().Lookup("dataFilePath")))
}
