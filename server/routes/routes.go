package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/internal/logs"
)

var logger = logs.Logger("routes")

type Routes struct {
	*gin.Engine
}

func RegistRoutes() Routes {
	router := gin.Default()
	r := Routes{
		router,
	}

	r.registRoute()

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

func (r Routes) registStorageRoute() {
	handleStorage(r.Group("/mefs"), handlerMefs())
	handleStorage(r.Group("/ipfs"), handlerIpfs())
}
