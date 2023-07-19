package cmd

import (
	"context"
	"fmt"

	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/contract"
	"github.com/urfave/cli/v2"
)

var ContractCmd = &cli.Command{
	Name:  "contract",
	Usage: "contract command",
	Subcommands: []*cli.Command{
		setPkgCmd,
		checkReceipt,
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
		&cli.IntFlag{
			Name:    "chainid",
			Aliases: []string{"c"},
			Usage:   "chainid",
			Value:   985,
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

		contracts := contract.NewContracts(cf.Contracts)
		ct, ok := contracts[chainid]
		if !ok {
			return fmt.Errorf("%d is not set", chainid)
		}
		receipt, err := ct.AdminAddPkgInfo(time, amount, kind, size)
		if err != nil {
			return err
		}

		fmt.Println(receipt)
		return nil
	},
}

var checkReceipt = &cli.Command{
	Name:  "check",
	Usage: "check [chainid] [receipt]",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "chainid",
			Aliases: []string{"c"},
			Usage:   "chainid",
			Value:   985,
		},
	},
	Action: func(ctx *cli.Context) error {
		chainid := ctx.Int("chainid")

		receipt := ctx.Args().Get(0)

		cf, err := config.ReadFile()
		if err != nil {
			return err
		}

		contracts := contract.NewContracts(cf.Contracts)
		ct, ok := contracts[chainid]
		if !ok {
			return fmt.Errorf("%d is not set", chainid)
		}
		err = ct.CheckTrsaction(context.Background(), receipt)
		if err != nil {
			return err
		}

		fmt.Println("check success")
		return nil
	},
}
