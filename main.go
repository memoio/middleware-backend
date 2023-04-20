package main

import (
	"fmt"
	"os"

	"github.com/memoio/backend/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	local := make([]*cli.Command, 0, 1)
	local = append(local, cmd.BackendCmd)
	local = append(local, cmd.InitCmd)
	local = append(local, cmd.ContractCmd)
	app := cli.App{
		Commands: local,
	}
	app.Setup()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", err) // nolint:errcheck
		os.Exit(1)
	}
}
