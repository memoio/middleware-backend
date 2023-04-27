package cmd

import (
	"fmt"

	"github.com/memoio/backend/config"
	db "github.com/memoio/backend/global/database"
	"github.com/urfave/cli/v2"
)

var InitCmd = &cli.Command{
	Name:  "init",
	Usage: "init middleware daemon",
	Action: func(ctx *cli.Context) error {
		if !db.InitDB() {
			return fmt.Errorf("init database error")
		}
		config.Init()

		return nil
	},
}
