package cmd

import (
	"fmt"

	db "github.com/memoio/backend/global/db"
	"github.com/urfave/cli/v2"
)

var InitCmd = &cli.Command{
	Name:  "init",
	Usage: "init middleware daemon",
	Action: func(ctx *cli.Context) error {
		if !db.InitDB() {
			return fmt.Errorf("init database error")
		}

		return nil
	},
}
