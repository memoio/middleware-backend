package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	db "github.com/memoio/backend/global/db"
	"github.com/memoio/backend/server"
	"github.com/urfave/cli/v2"
)

var BackendCmd = &cli.Command{
	Name:  "daemon",
	Usage: "middleware daemon",
	Subcommands: []*cli.Command{
		runCmd,
		stopCmd,
	},
}

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "run server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "endpoint",
			Aliases: []string{"e"},
			Usage:   "input your endpoint",
			Value:   ":8080",
		},
	},
	Action: func(ctx *cli.Context) error {
		endPoint := ctx.String("endpoint")
		srv := server.NewServer(endPoint)

		db.NewCron()

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")

		cctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(cctx); err != nil {
			log.Fatal("Server forced to shutdown: ", err)
		}

		log.Println("Server exiting")
		return nil
	},
}

var stopCmd = &cli.Command{
	Name:  "stop",
	Usage: "stop server",
	Action: func(ctx *cli.Context) error {
		return nil
	},
}
