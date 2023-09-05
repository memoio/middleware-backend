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

func newHandler() *handler {
	controller, err := controller.NewController()
	if err != nil {
		logger.Panic(err)
	}
	return &handler{
		controller: controller,
	}
}

func handlerMefs() api.IGateway {
	store, err := mefs.NewGateway()
	if err != nil {
		logger.Error("init mefs error:", err)
	}

	return store
}

func handlerIpfs() api.IGateway {
	store, err := ipfs.NewGateway()
	if err != nil {
		logger.Error("init ipfs error:", err)
	}

	return store
}

func (h *handler) handleStorage(r *gin.RouterGroup, store api.IGateway) {
	h.controller.SetStore(store)

	// OBJ
	r.POST("/putObject/", auth.VerifyIdentityHandler, h.putObjectHandle)
	r.GET("/getObject/:cid", auth.VerifyIdentityHandler, h.getObjectHandle)
	r.GET("/listObject", auth.VerifyIdentityHandler, h.listObjectsHandle)
	r.GET("/deleteObject", auth.VerifyIdentityHandler, h.deleteObjectHandle)
}

func (h *handler) handleAccount(r *gin.RouterGroup) {
	// info
	r.GET("/getBalance", auth.VerifyIdentityHandler, h.getBalanceHandle)

	// package
	r.GET("/getSpaceInfo", auth.VerifyIdentityHandler, h.getSpaceInfoHandle)
	r.GET("/getTrafficInfo", auth.VerifyIdentityHandler, h.getTrafficInfoHandle)
	r.GET("/getSpaceCheckHash", auth.VerifyIdentityHandler, h.getSpaceCheckHashHandle)
	r.GET("/getTrafficCheckHash", auth.VerifyIdentityHandler, h.getTrafficCheckHashHandle)
	r.GET("/getSpacePrice", h.getSpacePriceHandle)
	r.GET("/getTrafficPrice", h.getTrafficPriceHandle)
	r.GET("/buySpace", auth.VerifyIdentityHandler, h.buySpaceHandle)
	r.GET("/buyTraffic", auth.VerifyIdentityHandler, h.buyTrafficHandle)
	r.GET("/getApproveTsHash", auth.VerifyIdentityHandler, h.getApproveTsHash)
	r.GET("/getAllowance", auth.VerifyIdentityHandler, h.getAllowanceHandle)

	r.GET("/getReceipt", h.checkReceiptHandle)
}

func (h *handler) handleAdmin(r *gin.RouterGroup) {
	r.GET("/cashSpace", h.cashSpaceHandle)
	r.GET("/cashTraffic", h.cashTrafficHandle)
}
