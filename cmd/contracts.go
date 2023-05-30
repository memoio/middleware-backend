package cmd

import (
	"errors"
	"log"

	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/contract"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
)

var ContractCmd = &cli.Command{
	Name:  "contract",
	Usage: "contract command",
	Subcommands: []*cli.Command{
		setPkgCmd,
	},
}

var setPkgCmd = &cli.Command{
	Name:  "setpkg",
	Usage: "set package",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "time",
			Aliases: []string{"t"},
			Usage:   "time",
		},
		&cli.StringFlag{
			Name:    "amount",
			Aliases: []string{"a"},
			Usage:   "amount",
		},
		&cli.StringFlag{
			Name:    "kind",
			Aliases: []string{"k"},
			Usage:   "kind",
		},
		&cli.StringFlag{
			Name:    "size",
			Aliases: []string{"s"},
			Usage:   "size",
		},
	},
	Action: func(ctx *cli.Context) error {
		time := ctx.String("time")
		amount := ctx.String("amount")
		chainid := ctx.Int("chainid")
		kind := ctx.String("kind")
		size := ctx.String("size")
		cf, err := config.ReadFile()
		if err != nil {
			return err
		}

		contracts := contract.NewContract(cf.Contract)
		ct, ok := contracts[chainid]
		if !ok {
			return xerrors.Errorf("%s is not set", chainid)
		}
		flag := ct.AdminAddPkgInfo(time, amount, kind, size)
		if flag {
			log.Println("set package success!")
		} else {
			return errors.New("set package falid")
		}
		return nil
	},
}
