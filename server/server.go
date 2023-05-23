package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/controller"
	"github.com/memoio/backend/internal/logs"
)

type Server struct {
	Router       *gin.Engine
	Config       *config.Config
	NonceManager *NonceManager
	Controller   *controller.Controller
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

	InitAuthConfig(config.SecurityKey, config.Domain, config.LensAPIUrl)

	nonceManager := NewNonceManager(30*int64(time.Second.Seconds()), 1*int64(time.Minute.Seconds()))
	router := gin.Default()

	s := &Server{
		Config:       config,
		NonceManager: nonceManager,
		Router:       router,
	}

	s.registRoute()

	if opt.CheckRegistered {
		s.registLensLogin()
	}

	// go s.Controller.UploadToContract(context.TODO())

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

	s.Router.GET("/challenge", ChallengeHandler(s.NonceManager))

	s.Router.POST("/login", LoginHandler(s.NonceManager))

	s.Router.GET("/refresh", RefreshHandler())
}

func (s Server) registLensLogin() {
	s.Router.POST("/lens/login", LensLoginHandler(s.NonceManager, true))
}

func (s Server) registController() {
	for k := range controller.ApiMap {
		r := s.Router.Group(k)
		ct := controller.NewController(r.BasePath(), s.Config)
		s.Controller = ct

		go s.Controller.UploadToContract()
		s.StorageRegistRoutes(r)
		s.accountRegistRoutes(r)
	}
}
