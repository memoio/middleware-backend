package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/config"
	auth "github.com/memoio/backend/internal/authentication"
	"github.com/memoio/backend/internal/controller"
	"github.com/memoio/backend/internal/logs"
)

type Server struct {
	Router     *gin.Engine
	Config     *config.Config
	Controller *controller.Controller
}

type ServerOption struct {
	Endpoint        string
	CheckRegistered bool
}

type AuthenticationFaileMessage struct {
	Nonce string
	Error logs.APIError
}

func NewServer(opt ServerOption) *http.Server {
	log.Println("Server Start")
	gin.SetMode(gin.ReleaseMode)

	config, err := config.ReadFile()
	if err != nil {
		log.Fatal("config not right")
		return nil
	}

	router := gin.Default()

	s := &Server{
		Config: config,
		Router: router,
	}

	s.registRoute()

	srv := &http.Server{
		Addr:    opt.Endpoint,
		Handler: s.Router,
	}

	return srv
}

func (s Server) registRoute() {
	// add storage routes

	s.Router.Use(Cors())
	s.Router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Server")
	})

	s.registLogin()
	s.registController()
}

func (s Server) registLogin() {
	auth.LoadAuthRouter(s.Router.Group("/"))
}

func (s Server) registController() {
	for k := range controller.ApiMap {
		r := s.Router.Group(k)
		ct := controller.NewController(r.BasePath(), s.Config)
		s.Controller = ct

		s.StorageRegistRoutes(r)
		s.accountRegistRoutes(r)
	}
}
