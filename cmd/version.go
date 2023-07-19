package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

const Version = "0.2.0"

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
