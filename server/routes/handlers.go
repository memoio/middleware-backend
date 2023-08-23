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
	// OBJ
	r.POST("/putOBJ/", auth.VerifyIdentityHandler, h.putObjectHandle)
	r.GET("/getOBJ/:cid", auth.VerifyIdentityHandler, h.getObjectHandle)
	r.GET("/listOBJ", auth.VerifyIdentityHandler, h.listObjectsHandle)
	r.GET("/deleteOBJ", auth.VerifyIdentityHandler, h.deleteObjectHandle)

	// info
	r.GET("/getBalance", auth.VerifyIdentityHandler, h.getBalanceHandle)

	// package
	r.GET("/getSpace", auth.VerifyIdentityHandler, h.getSpace)
	r.GET("/getTraffic", auth.VerifyIdentityHandler, h.getTraffic)
	r.GET("/getSpaceHash", auth.VerifyIdentityHandler, h.spaceHash)
	r.GET("/getTrafficHash", auth.VerifyIdentityHandler, h.trafficHash)
	r.GET("/getSpacePrice", h.spacePrice)
	r.GET("/getTrafficPrice", h.trafficPrice)
	r.GET("/buySpace", auth.VerifyIdentityHandler, h.BuySpace)
	r.GET("/buyTraffic", auth.VerifyIdentityHandler, h.BuyTraffic)
	r.GET("/getApproveHash", auth.VerifyIdentityHandler, h.Approve)
	r.GET("/getAllowance", auth.VerifyIdentityHandler, h.allowance)

	r.GET("/cashSpace", h.cashSpace)
	r.GET("/cashTraffic", h.cashTraffic)
	r.GET("/getReceipt", h.checkReceipt)

}
