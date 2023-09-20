package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/memoio/backend/internal/wallet"
	"github.com/urfave/cli/v2"
)

var (
	ksp = filepath.Join("./", "keystore")
)

var WalletCmd = &cli.Command{
	Name:  "wallet",
	Usage: "wallet cmd",
	Subcommands: []*cli.Command{
		newWalletCmd,
		listWalletCmd,
		exportCmd,
		importCmd,
	},
}

var newWalletCmd = &cli.Command{
	Name:  "new",
	Usage: "create a new address",
	Action: func(_ *cli.Context) error {
		ks, err := wallet.NewKeyRepo(ksp)
		if err != nil {
			return err
		}
		wl := wallet.New(ks)
		addr, err := wl.WalletNew(context.Background())
		if err != nil {
			return err
		}

		fmt.Println(addr)
		return nil
	},
}

var listWalletCmd = &cli.Command{
	Name:  "list",
	Usage: "list all addresses",
	Action: func(_ *cli.Context) error {
		ks, err := wallet.NewKeyRepo(ksp)
		if err != nil {
			return err
		}
		wl := wallet.New(ks)
		addrs, err := wl.WalletList(context.Background())
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			fmt.Println(addr)
		}
		return nil
	},
}

var exportCmd = &cli.Command{
	Name:  "export",
	Usage: "export address",
	Action: func(ctx *cli.Context) error {
		arg := ctx.Args()

		ks, err := wallet.NewKeyRepo(ksp)
		if err != nil {
			return err
		}
		addr := arg.Get(0)
		wl := wallet.New(ks)
		sk, err := wl.WalletExport(context.Background(), common.HexToAddress(addr))
		if err != nil {
			return err
		}

		fmt.Println(hex.EncodeToString(sk))
		return nil
	},
}

var importCmd = &cli.Command{
	Name:  "import",
	Usage: "import address",
	Action: func(ctx *cli.Context) error {
		arg := ctx.Args()

		ks, err := wallet.NewKeyRepo(ksp)
		if err != nil {
			return err
		}
		sk := arg.Get(0)
		wl := wallet.New(ks)
		skb, err := hex.DecodeString(sk)
		if err != nil {
			return err
		}
		address, err := wl.WalletImport(context.Background(), skb)
		if err != nil {
			return err
		}

		fmt.Println(address)
		return nil
	},
}
