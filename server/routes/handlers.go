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

var handlerMap map[string]handler

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

	r.GET("/receipt", h.checkReceipt)
}
