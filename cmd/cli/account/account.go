package account

import (
	"github.com/spf13/cobra"

	"github.com/go-rfe/gpwd/cmd/root"
)

// accountCmd represents the account management command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage gpwd account on the server",
	Long:  `This is a cli for account manipulation`,
}

func init() {
	root.AddCommand(accountCmd)
}
