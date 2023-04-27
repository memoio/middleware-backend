package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/gateway"
)

type Server struct {
	Router       *gin.Engine
	Config       *config.Config
	NonceManager *NonceManager
}

type AuthenticationFaileMessage struct {
	Nonce string
	Error gateway.APIError
}

func NewServer(endpoint string, checkRegistered bool) *http.Server {
	log.Println("Server Start")
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(Cors())
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Server")
	})

	nonceManager := NewNonceManager(30*int64(time.Second.Seconds()), 1*int64(time.Minute.Seconds()))

	router.GET("/challenge", ChallengeHandler(nonceManager))

	router.POST("/login", LoginHandler(nonceManager))

	router.POST("/lens/login", LensLoginHandler(nonceManager, checkRegistered))

	router.GET("/refresh", RefreshHandler())

	config, err := config.ReadFile()
	if err != nil {
		log.Fatal("config not right")
		return nil
	}
	InitAuthConfig(config.SecurityKey, config.Domain, config.LensAPIUrl)

	s := &Server{
		Router:       router,
		Config:       config,
		NonceManager: nonceManager,
	}

	s.registRoute()

	srv := &http.Server{
		Addr:    endpoint,
		Handler: s.Router,
	}

	return srv
}

func (s Server) registRoute() {
	mefs := s.Router.Group("/mefs")
	s.mefsRegistRoutes(mefs)
	ipfs := s.Router.Group("/ipfs")
	s.ipfsRegistRoutes(ipfs)
	account := s.Router.Group("/account")
	s.accountRegistRoutes(account)
	test := s.Router.Group("/test")
	s.testRegistRoutes(test)
}
