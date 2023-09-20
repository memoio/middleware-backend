package cmd

import "github.com/urfave/cli/v2"

var CommonCmd = []*cli.Command{
	BackendCmd,
	WalletCmd,
	VersionCmd,
	UserCmd,
}
