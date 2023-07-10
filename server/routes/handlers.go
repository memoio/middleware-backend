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
	r.GET("/storageinfo", auth.VerifyIdentityHandler, h.getStorageInfoHandle)
	r.GET("/flowsize", auth.VerifyIdentityHandler, h.getFlowSize)

	// package
	r.GET("/pkginfos", h.getPkgInfos)
	r.GET("/getbuypkgs", h.getBuyPackages)
	r.GET("/buypkg", h.buyPackage)

	r.GET("/receipt", h.checkReceipt)
}
