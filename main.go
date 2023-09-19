package main

import (
	"fmt"
	"os"

	"github.com/memoio/backend/cmd"
	"github.com/urfave/cli/v2"
)

//	@title			MiddleWare API
//	@version		1.0
//	@description	This is a middleware server.

// @host		localhost:8081
// @BasePath	/
func main() {
	local := make([]*cli.Command, 0, 1)
	local = append(local, cmd.CommonCmd...)
	app := cli.App{
		Commands: local,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Show application version",
			},
		},
		Action: func(_ *cli.Context) error {
			fmt.Println(cmd.Version + "+" + cmd.BuildFlag)
			return nil
		},
	}
	app.Setup()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", err) // nolint:errcheck
		os.Exit(1)
	}
}
