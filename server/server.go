package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/contract"
	"github.com/memoio/backend/gateway"
)

type Server struct {
	Router       *gin.Engine
	Gateway      *gateway.Gateway
	Config       *config.Config
	NonceManager *NonceManager
}

type AuthenticationFaileMessage struct {
	Nonce string
	Error gateway.APIError
}

func NewServer(endpoint string) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(Cors())
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Server")
	})

	nonceManager := NewNonceManager(30*int64(time.Second.Seconds()), 1*int64(time.Minute.Seconds()))

	router.GET("/getnonce", func(c *gin.Context) {
		nonce := nonceManager.GetNonce()
		c.String(http.StatusOK, nonce)
	})

	router.POST("/login", LoginHandler(nonceManager))

	router.POST("/lens/login", LensLoginHandler(nonceManager))

	router.GET("/fresh", FreshHandler())

	router.GET("/pay", PayHandler())

	router.GET("/test", TestHandler())

	router.GET("/db", DBHandler())

	config, err := config.ReadFile("")
	if err != nil {
		log.Fatal("config not right")
		return nil
	}
	g := gateway.NewGateway(config)

	s := &Server{
		Router:       router,
		Gateway:      g,
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
	s.commonregistRoutes(mefs, gateway.MEFS)
	ipfs := s.Router.Group("/ipfs")
	s.commonregistRoutes(ipfs, gateway.IPFS)
}

func (s Server) commonregistRoutes(r *gin.RouterGroup, storage gateway.StorageType) {
	s.addPutobjectRoutes(r, storage)
	s.addGetObjectRoutes(r, storage)
	s.addListObjectRoutes(r, storage)
	s.addGetPriceRoutes(r, storage)
	s.addGetStorageRoutes(r, storage)
	s.addGetBalanceRoutes(r, storage)
}

func TestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Query("address")
		si, _ := contract.GetPkgSize(address)

		c.JSON(http.StatusOK, si)
	}
}
