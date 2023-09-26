package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	auth "github.com/memoio/backend/internal/authentication"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/market"
	"github.com/memoio/backend/internal/share"
)

var logger = logs.Logger("routes")

type Routes struct {
	*gin.Engine
}

func RegistRoutes(checkRegistered bool) Routes {
	router := gin.Default()
	r := Routes{
		router,
	}

	r.registRoute()

	r.registLoginRoute(checkRegistered)
	r.registMarketRoute()
	r.registShareRoute()

	r.registStorageRoute()
	return r
}

func (r Routes) registRoute() {
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Use(Cors())
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Server")
	})
}

func (r Routes) registLoginRoute(checkRegistered bool) {
	auth.LoadAuthModule(r.Group("/"), checkRegistered)
}

func (r Routes) registMarketRoute() {
	market.LoadNFTMarketModule(r.Group("nft"))
}

func (r Routes) registShareRoute() {
	share.LoadShareModule(r.Group("/"))
}

func (r Routes) registStorageRoute() {
	handleStorage(r.Group("/mefs"), handlerMefs())
	handleStorage(r.Group("/ipfs"), handlerIpfs())
}
