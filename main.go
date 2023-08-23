package main

import (
	"fmt"
	"os"

	"github.com/memoio/backend/cmd"
	"github.com/urfave/cli/v2"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8081/
//	@BasePath	/
func main() {
	local := make([]*cli.Command, 0, 1)
	local = append(local, cmd.BackendCmd)
	local = append(local, cmd.WalletCmd)
	local = append(local, cmd.VersionCmd)
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
