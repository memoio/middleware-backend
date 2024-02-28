package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/docs"
	auth "github.com/memoio/backend/internal/authentication"
	"github.com/memoio/backend/internal/da"
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
	swaghost := config.Cfg.SwagHost
	if swaghost != "" {
		docs.SwaggerInfo.Host = swaghost
	}
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r := Routes{
		router,
	}

	r.registRoute()
	r.registLoginRoute()
	r.registShareRoute()
	r.registFileDnsRoute()
	// r.registAccount()
	r.registStorageRoute()
	r.registDARoute()
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

func (r Routes) registDARoute() {
	da.LoadDAModule(r.Group("/da"))
}

// func (r Routes) registAccount() {
// 	account.LoadAccountModule(r.Group("/account"))
// }

func (r Routes) registStorageRoute() {
	h := newHandler()
	h.handleStorage(r.Group("/mefs", auth.VerifyAccessTokenHandler, LoadMefsHandler()))
	// h.handleStorage(r.Group("/mefs", testLoadAddress(), LoadMefsHandler()))
	h.handleStorage(r.Group("/ipfs", auth.VerifyAccessTokenHandler, LoadIpfsHandler()))
}

// func testLoadAddress() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		ctx.Set("address", ctx.Query("address"))
// 	}
// }
