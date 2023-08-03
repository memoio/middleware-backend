package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Routes struct {
	*gin.Engine
}

func RegistRoutes() Routes {
	router := gin.Default()
	r := Routes{
		router,
	}

	r.registRoute()
	r.registLoginRoute()
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

func (r Routes) registLoginRoute() {

}

func (r Routes) registStorageRoute() {
	Init()
	handleStorage(r.Group("/mefs"), handlerMap["mefs"])
	handleStorage(r.Group("/ipfs"), handlerMap["ipfs"])
}
