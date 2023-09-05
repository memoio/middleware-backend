package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/memoio/backend/server"
	"github.com/mitchellh/go-homedir"
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

		opts := server.ServerOption{
			Endpoint: endPoint,
		}

		srv := server.NewServer(opts)

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

		pidpath, err := homedir.Expand("./")
		if err != nil {
			return nil
		}

		pid := os.Getpid()
		pids := []byte(strconv.Itoa(pid))
		err = os.WriteFile(path.Join(pidpath, "pid"), pids, 0644)
		if err != nil {
			return err
		}

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
	Action: func(_ *cli.Context) error {
		pidpath, err := homedir.Expand("./")
		if err != nil {
			return nil
		}

		pd, _ := ioutil.ReadFile(path.Join(pidpath, "pid"))

		err = kill(string(pd))
		if err != nil {
			return err
		}
		log.Println("gateway gracefully exit...")

		return nil
	},
}

func kill(pid string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("kill", "-15", pid).Run()
	case "windows":
		return exec.Command("taskkill", "/F", "/T", "/PID", pid).Run()
	default:
		return fmt.Errorf("unsupported platform %s", runtime.GOOS)
	}
}
