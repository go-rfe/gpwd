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
)

// listCmd represents the create command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list secrets using gpwd agent",
	Long:  `cli connects to the agent and list secrets`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
		defer cancel()

		client, err := secrets.NewSecretsClient(ctx, viper.GetString("socket_path"))
		cobra.CheckErr(err)

		secrets, err := client.List()
		cobra.CheckErr(err)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.Escape)
		_, err = fmt.Fprintln(w, "ID", "\t", "Labels", "\t", "Created At", "\t", "Updated At")
		cobra.CheckErr(err)

		var labels []byte
		for _, secret := range secrets {
			if secret.Labels != nil {
				labels, err = json.Marshal(secret.Labels)
				cobra.CheckErr(err)
			}

			updatedAt := ""
			if secret.UpdatedAt != nil {
				updatedAt = secret.GetUpdatedAt().AsTime().String()
			}
			_, err = fmt.Fprintln(w, secret.ID, "\t",
				string(labels), "\t",
				secret.GetCreatedAt().AsTime().String(), "\t",
				updatedAt)
			cobra.CheckErr(err)
		}
		cobra.CheckErr(w.Flush())
	},
}

func init() {
	secretCmd.AddCommand(listCmd)
}
