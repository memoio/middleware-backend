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

func NewServer(endpoint string, checkRegistered bool) *http.Server {
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

	router.POST("/lens/login", LensLoginHandler(nonceManager, checkRegistered))

	router.GET("/fresh", FreshHandler())

	config, err := config.ReadFile("")
	if err != nil {
		log.Fatal("config not right")
		return nil
	}
	InitAuthConfig(config.SecurityKey, config.Domain, config.LensAPIUrl)
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
	account := s.Router.Group("/account")
	s.accountRegistRoutes(account)
	// test := s.Router.Group("/test")
	// s.testRegistRoutes(test)
}

func (s Server) commonRegistRoutes(r *gin.RouterGroup, storage gateway.StorageType) {
	s.addPutobjectRoutes(r, storage)
	s.addGetObjectRoutes(r, storage)
	s.addListObjectRoutes(r, storage)
	s.addGetPriceRoutes(r, storage)

}

func (s Server) accountRegistRoutes(r *gin.RouterGroup) {
	s.addGetBalanceRoutes(r)
	s.addGetStorageRoutes(r)
	s.addBuyPkgRoutes(r)
	s.addGetPkgListRoutes(r)
	s.addGetBuyPkgRoutes(r)
}

func (s Server) testRegistRoutes(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/storage", func(c *gin.Context) {
		address := c.Query("address")
		si, err := s.Gateway.GetPkgSize(c.Request.Context(), gateway.MEFS, address)
		if err != nil {
			c.JSON(516, err.Error())
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

		si, err := s.Gateway.GetPkgSize(c.Request.Context(), gateway.MEFS, address)
		if err != nil {
			c.JSON(516, err)
			return
		}

		c.JSON(http.StatusOK, si)
	})

	p.POST("/update", func(c *gin.Context) {
		address := c.Query("address")
		hashid := c.Query("hash")

		si, err := s.Gateway.TestUpdatePkg(c.Request.Context(), gateway.MEFS, address, hashid, 1024)
		if err != nil {
			log.Println("TEST: ", err)
			c.JSON(520, err.Error())
			return
		}

		c.JSON(http.StatusOK, si)
	})

	p.GET("/balance", func(c *gin.Context) {
		address := c.Query("address")

		balance := contract.BalanceOf(c.Request.Context(), address)
		c.JSON(http.StatusOK, BalanceResponse{Address: address, Balance: balance.String()})
	})

	p.GET("/pay", func(c *gin.Context) {
		address := c.Query("address")
		hashid := c.Query("hash")

		p := s.Gateway.TestPay(c.Request.Context(), address, hashid, 1, 1024)

		c.JSON(http.StatusOK, p)
	})

	p.GET("/pkginfo", func(c *gin.Context) {
		result, err := contract.StoreGetPkgInfos()
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}
		c.JSON(http.StatusOK, result)
	})
	p.GET("/setpkg", func(c *gin.Context) {
		time := c.Query("time")
		amount := c.Query("amount")
		kind := c.Query("kind")
		size := c.Query("size")

		flag := contract.AdminAddPkgInfo(time, amount, kind, size)
		c.JSON(http.StatusOK, flag)
	})
	p.GET("buypkg", func(c *gin.Context) {
		address := c.Query("address")
		chainId := c.Query("chainid")
		times := time.Now()
		flag := contract.StoreBuyPkg(address, 2, 1, uint64(times.Second()), chainId)
		c.JSON(http.StatusOK, flag)
	})
	p.GET("getbuypkg", func(c *gin.Context) {
		address := c.Query("address")
		pi, err := contract.StoreGetBuyPkgs(address)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, pi)
	})

}
