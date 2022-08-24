package main

import (
	"github.com/go-rfe/gpwd/cmd/root"

	// Register commands
	_ "github.com/go-rfe/gpwd/cmd/agent"
	_ "github.com/go-rfe/gpwd/cmd/cli/account"
	_ "github.com/go-rfe/gpwd/cmd/cli/secret"
	_ "github.com/go-rfe/gpwd/cmd/server"
)

func main() {
	root.Execute()
}
