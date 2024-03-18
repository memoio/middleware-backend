package cmd

import "github.com/urfave/cli/v2"

var CommonCmd = []*cli.Command{
	BackendCmd,
	DaCmd,
	WalletCmd,
	VersionCmd,
	UserCmd,
}
