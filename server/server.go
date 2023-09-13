package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/internal/filedns"
	"github.com/memoio/backend/server/routes"
)

type ServerOption struct {
	Endpoint string
}

func NewServer(opt ServerOption) *http.Server {
	log.Println("Init File Dns")
	filedns.InitFileDns(filedns.DefaultSearcherOpts)

	dumper, err := filedns.NewMfileDumper("dev")
	if err != nil {
		panic(err.Error())
	}
	go dumper.DumpMfileDID()

	log.Println("Server Start")
	gin.SetMode(gin.ReleaseMode)

	// register routes
	router := routes.RegistRoutes()

	// start server
	srv := &http.Server{
		Addr:    opt.Endpoint,
		Handler: router,
	}

	return srv
}
