package account

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-rfe/gpwd/internal/client/accounts"
)

// getCmd represents the create command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get account from gpwd agent",
	Long:  "cli connects to the agent and gets account",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("timeout"))
		defer cancel()

		client, err := accounts.NewAccountsClient(ctx, viper.GetString("socket_path"))
		cobra.CheckErr(err)

		account, err := client.Get()
		cobra.CheckErr(err)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.Escape)
		_, err = fmt.Fprintln(w, "ID", "\t", "Server", "\t", "Username")
		cobra.CheckErr(err)

		_, err = fmt.Fprintln(w, account.GetID(), "\t", account.GetServerAddress(), "\t", account.GetUserName())
		cobra.CheckErr(err)
		cobra.CheckErr(w.Flush())

	},
}

func init() {
	accountCmd.AddCommand(getCmd)
}
