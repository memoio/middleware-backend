package cmd

import (
	"github.com/urfave/cli/v2"
)

var ContractCmd = &cli.Command{
	Name:        "contract",
	Usage:       "contract command",
	Subcommands: []*cli.Command{},
}
