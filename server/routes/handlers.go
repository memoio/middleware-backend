package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/api"
	auth "github.com/memoio/backend/internal/authentication"
	"github.com/memoio/backend/internal/gateway/ipfs"
	"github.com/memoio/backend/internal/gateway/mefs"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/server/routes/controller"
)

var logger = logs.Logger("routes")

type handler struct {
	controller *controller.Controller
}

func newHandler(store api.IGateway, path string) *handler {
	controller, err := controller.NewController(store, path)
	if err != nil {
		logger.Panic(err)
	}
	return &handler{
		controller: controller,
	}
}

func handlerMefs() *handler {
	store, err := mefs.NewGateway()
	if err != nil {
		logger.Error("init mefs error:", err)
	}

	return newHandler(store, "mefs")
}

func handlerIpfs() *handler {
	store, err := ipfs.NewGateway()
	if err != nil {
		logger.Error("init ipfs error:", err)
	}

	return newHandler(store, "ipfs")
}

func (ro Routes) handleStorage(r *gin.RouterGroup, h *handler) {

	// OBJ
	r.POST("/putObject/", auth.VerifyIdentityHandler, h.putObjectHandle)
	r.POST("/getObject/:cid", auth.VerifyIdentityHandler, h.getObjectHandle)
	r.POST("/listObject", auth.VerifyIdentityHandler, h.listObjectsHandle)
	r.POST("/deleteObject", auth.VerifyIdentityHandler, h.deleteObjectHandle)
}

func (h *handler) handleAccount(r *gin.RouterGroup) {
	// info
	r.POST("/getBalance", auth.VerifyIdentityHandler, h.getBalanceHandle)

	// package
	r.POST("/getSpaceInfo", auth.VerifyIdentityHandler, h.getSpaceInfoHandle)
	r.POST("/getTrafficInfo", auth.VerifyIdentityHandler, h.getTrafficInfoHandle)
	r.POST("/getSpaceCheckHash", auth.VerifyIdentityHandler, h.getSpaceCheckHashHandle)
	r.POST("/getTrafficCheckHash", auth.VerifyIdentityHandler, h.getTrafficCheckHashHandle)
	r.GET("/getSpacePrice", h.getSpacePriceHandle)
	r.GET("/getTrafficPrice", h.getTrafficPriceHandle)
	r.POST("/buySpace", auth.VerifyIdentityHandler, h.buySpaceHandle)
	r.POST("/buyTraffic", auth.VerifyIdentityHandler, h.buyTrafficHandle)
	r.POST("/getApproveTsHash", auth.VerifyIdentityHandler, h.getApproveTsHash)
	r.POST("/getAllowance", auth.VerifyIdentityHandler, h.getAllowanceHandle)

	r.GET("/getReceipt", h.checkReceiptHandle)

	r.GET("/cashSpace", h.cashSpaceHandle)
	r.GET("/cashTraffic", h.cashTrafficHandle)
}
