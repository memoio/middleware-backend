package cmd

import (
	"context"
	"fmt"

	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/database"
	"github.com/urfave/cli/v2"
)

var UserCmd = &cli.Command{
	Name:  "mefs",
	Usage: "mefs options",
	Subcommands: []*cli.Command{
		addCmd,
		stopCmd,
	},
}

var addCmd = &cli.Command{
	Name:  "add",
	Usage: "add user",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "api",
			Aliases: []string{"a"},
			Usage:   "input api",
		},
		&cli.StringFlag{
			Name:    "token",
			Aliases: []string{"t"},
			Usage:   "input api",
		},
		&cli.StringFlag{
			Name:    "area",
			Aliases: []string{"area"},
			Usage:   "input api",
			Value:   "default",
		},
	},
	Action: func(ctx *cli.Context) error {
		area := ctx.String("area")
		apis := ctx.String("api")
		token := ctx.String("token")
		if area == "" {
			fmt.Println("area is nil")
			return nil
		}
		if apis == "" {
			fmt.Println("api is nil")
			return nil
		}
		if token == "" {
			fmt.Println("token is nil")
			return nil
		}

		db := database.NewDataBase()
		ui := api.USerInfo{
			Area:  area,
			Api:   apis,
			Token: token,
		}
		err := db.AddUser(context.Background(), ui)
		if err != nil {
			return err
		}
		return nil
	},
}
