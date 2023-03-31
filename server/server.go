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
	Gateway      *gateway.Gateway
	Config       *config.Config
	NonceManager *NonceManager
}

type AuthenticationFaileMessage struct {
	Nonce string
	Error gateway.APIError
}

func NewServer(endpoint string) *http.Server {
	log.Println("Server Start")

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
	s.commonRegistRoutes(mefs, gateway.MEFS)
	ipfs := s.Router.Group("/ipfs")
	s.commonRegistRoutes(ipfs, gateway.IPFS)
	// test := s.Router.Group("/test")
	// s.testregistRoutes(test)
}

func (s Server) commonRegistRoutes(r *gin.RouterGroup, storage gateway.StorageType) {
	s.addPutobjectRoutes(r, storage)
	s.addGetObjectRoutes(r, storage)
	s.addListObjectRoutes(r, storage)
	s.addGetPriceRoutes(r, storage)
	s.addGetStorageRoutes(r, storage)
	s.addGetBalanceRoutes(r, storage)
}

func (s Server) testregistRoutes(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/storage", func(c *gin.Context) {
		address := c.Query("address")
		si, err := s.Gateway.GetPkgSize(c.Request.Context(), address)
		if err != nil {
			c.JSON(516, err)
			return
		}

		c.JSON(http.StatusOK, si)
	})

	p.POST("/put", func(c *gin.Context) {
		address := c.Query("address")
		hashid := c.Query("hash")

		err := s.Gateway.TestPutobject(c.Request.Context(), address, hashid, 1024)
		if err != nil {
			log.Println("TEST: ", err)
			c.JSON(520, err.Error())
			return
		}

		si, err := s.Gateway.GetPkgSize(c.Request.Context(), address)
		if err != nil {
			c.JSON(516, err)
			return
		}

		c.JSON(http.StatusOK, si)
	})

	p.POST("/update", func(c *gin.Context) {
		address := c.Query("address")
		hashid := c.Query("hash")

		si, err := s.Gateway.TestUpdatePkg(c.Request.Context(), address, hashid, 1024)
		if err != nil {
			log.Println("TEST: ", err)
			c.JSON(520, err.Error())
			return
		}

		c.JSON(http.StatusOK, si)
	})
}
