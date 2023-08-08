package routes

import (
	"github.com/gin-gonic/gin"
	auth "github.com/memoio/backend/internal/authentication"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/server/controller"
)

var logger = logs.Logger("routes")

type handler struct {
	controller *controller.Controller
}

func handleStorage(r *gin.RouterGroup, h handler) {
	// store
	r.POST("/", auth.VerifyIdentityHandler, h.putObjectHandle)
	r.GET("/:cid", auth.VerifyIdentityHandler, h.getObjectHandle)
	r.GET("/listobjects", auth.VerifyIdentityHandler, h.listObjectsHandle)
	r.GET("/delete", auth.VerifyIdentityHandler, h.deleteObjectHandle)

	// info
	r.GET("/balance", auth.VerifyIdentityHandler, h.getBalanceHandle)

	// package
	r.GET("/space", auth.VerifyIdentityHandler, h.getSpace)
	r.GET("/traffic", auth.VerifyIdentityHandler, h.getTraffic)
	r.GET("/spacehash", auth.VerifyIdentityHandler, h.spaceHash)
	r.GET("/traffichash", auth.VerifyIdentityHandler, h.trafficHash)
	r.GET("/buyspace", auth.VerifyIdentityHandler, h.BuySpace)
	r.GET("/buytraffic", auth.VerifyIdentityHandler, h.BuyTraffic)
	r.GET("/approve", auth.VerifyIdentityHandler, h.Approve)
	r.GET("/allowance", auth.VerifyIdentityHandler, h.allowance)

	r.GET("/cashspace", h.cashSpace)
	r.GET("/cashtraffic", h.cashTraffic)
	r.GET("/receipt", h.checkReceipt)

}
