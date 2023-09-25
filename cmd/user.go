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
		addUSerCmd,
		listUSerCmd,
		deleteUSerCmd,
	},
}

var addUSerCmd = &cli.Command{
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
			Usage:   "input token",
		},
		&cli.StringFlag{
			Name:    "area",
			Aliases: []string{"ar"},
			Usage:   "input area",
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

var listUSerCmd = &cli.Command{
	Name:  "list",
	Usage: "list user",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "area",
			Aliases: []string{"ar"},
			Usage:   "input area",
		},
	},
	Action: func(ctx *cli.Context) error {
		area := ctx.String("area")

		db := database.NewDataBase()
		uis, err := db.ListUsers(context.TODO(), area)
		if err != nil {
			return err
		}
		for _, ui := range uis {
			fmt.Println(ui)
		}
		return nil
	},
}

var deleteUSerCmd = &cli.Command{
	Name:  "del",
	Usage: "delete user",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "id",
			Aliases: []string{"i"},
			Usage:   "user id",
			Value:   -1,
		},
	},
	Action: func(ctx *cli.Context) error {
		id := ctx.Int("id")
		if id == -1 {
			fmt.Println("id is nil")
			return nil
		}
		db := database.NewDataBase()
		err := db.DeleteUser(context.TODO(), id)
		if err != nil {
			return err
		}
		return nil
	},
}
