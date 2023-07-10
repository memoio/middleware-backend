package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/config"
	auth "github.com/memoio/backend/internal/authentication"
	"github.com/memoio/backend/internal/share"
	"github.com/memoio/backend/server/routes"
)

type Server struct {
	Router routes.Routes
	Config *config.Config
}

type ServerOption struct {
	Endpoint        string
	CheckRegistered bool
}

// type AuthenticationFaileMessage struct {
// 	Nonce string
// 	Error logs.APIError
// }

func NewServer(opt ServerOption) *http.Server {
	log.Println("Server Start")
	gin.SetMode(gin.ReleaseMode)

	config, err := config.ReadFile()
	if err != nil {
		log.Fatal("config not right ", err)
		return nil
	}

	auth.InitAuthConfig(config.SecurityKey, config.Domain, config.LensAPIUrl)

	router := routes.RegistRoutes()

	s := &Server{
		Config: config,
		Router: router,
	}

	s.registRoute(opt.CheckRegistered)

	srv := &http.Server{
		Addr:    opt.Endpoint,
		Handler: s.Router,
	}

	return srv
}

func (s Server) registRoute(checkRegistered bool) {
	s.registLogin(checkRegistered)
	s.registShare()
}

func (s Server) registLogin(checkRegistered bool) {
	auth.LoadAuthModule(s.Router.Group("/"), checkRegistered)
}

func (s Server) registShare() {
	share.LoadAuthModule(s.Router.Group("/"))
}
