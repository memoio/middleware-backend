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
	handler *handler
}

func RegistRoutes() Routes {
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	handler := newHandler()

	r := Routes{
		router,
		handler,
	}

	r.registRoute()
	r.registLoginRoute()
	r.registShareRoute()
	r.registFileDnsRoute()
	r.registAccount()
	r.registAdmin()
	r.registStorageRoute()
	return r
}

func (r Routes) registRoute() {
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Use(Cors())
	r.Use(ErrorHandler())
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

func (r Routes) registAccount() {
	r.handler.handleAccount(r.Group("/account"))
}

func (r Routes) registAdmin() {
	r.handler.handleAdmin(r.Group("/admin"))
}

func (r Routes) registStorageRoute() {
	r.handler.handleStorage(r.Group("/mefs"), handlerMefs())
	r.handler.handleStorage(r.Group("/ipfs"), handlerIpfs())
}
