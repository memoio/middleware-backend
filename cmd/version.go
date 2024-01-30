package cmd

import (
	"fmt"

	"github.com/memoio/backend/api"
	"github.com/urfave/cli/v2"
)

const Version = api.Version

var BuildFlag string

var VersionCmd = &cli.Command{
	Name:    "version",
	Usage:   "print version",
	Aliases: []string{"V"},
	Action: func(ctx *cli.Context) error {
		fmt.Println(Version + "+" + BuildFlag)
		return nil
	},
}
