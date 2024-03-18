package cmd

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/da"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
)

var DaCmd = &cli.Command{
	Name:  "da",
	Usage: "middleware data access",
	Subcommands: []*cli.Command{
		daRunCmd,
		daStopCmd,
	},
}

var daRunCmd = &cli.Command{
	Name:  "test",
	Usage: "test data access",
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

		srv, err := NewDataAccessServer(endPoint)
		if err != nil {
			log.Fatalf("new data access server: %s\n", err)
		}

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

var daStopCmd = &cli.Command{
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

func NewDataAccessServer(endpoint string) (*http.Server, error) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	// r.Use(Cors())
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Server")
	})
	da.LoadDAModule(router.Group("/da"))

	submitterSk, err := crypto.HexToECDSA(config.Cfg.DataAccess.SubmitterSecurityKey)
	if err != nil {
		return nil, err
	}

	prover, err := da.NewDataAccessProver("dev", submitterSk)
	if err != nil {
		return nil, err
	}
	go prover.ProveDataAccess()

	return &http.Server{
		Addr:    endpoint,
		Handler: router,
	}, nil
}
