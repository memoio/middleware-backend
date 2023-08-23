package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/docs"
	auth "github.com/memoio/backend/internal/authentication"
	"github.com/memoio/backend/internal/filedns"
	"github.com/memoio/backend/internal/share"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Routes struct {
	*gin.Engine
}

func RegistRoutes() Routes {
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r := Routes{
		router,
	}

	r.registRoute()
	r.registLoginRoute()
	r.registShareRoute()
	r.registFileDnsRoute()
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
	auth.LoadAuthModule(r.Group("/"))
}

func (r Routes) registShareRoute() {
	share.LoadShareModule(r.Group("/"))
}

func (r Routes) registFileDnsRoute() {
	filedns.LoadFileDnsModule(r.Group("/"))
}

func (r Routes) registStorageRoute() {
	handleStorage(r.Group("/mefs"), handlerMefs())
	handleStorage(r.Group("/ipfs"), handlerIpfs())
}
