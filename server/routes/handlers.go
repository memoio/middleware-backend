package routes

import (
	"github.com/gin-gonic/gin"
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

func (h *handler) handleStorage(r *gin.RouterGroup) {
	// OBJ
	r.POST("/putObject/", h.putObjectHandle)
	r.POST("/getObject/:cid", h.getObjectHandle)
	r.POST("/listObject", h.listObjectsHandle)
	r.POST("/deleteObject", h.deleteObjectHandle)

	r.POST("/getBalance", h.getBalanceHandle)

	// package
	r.POST("/getSpaceInfo", h.getSpaceInfoHandle)
	r.POST("/getTrafficInfo", h.getTrafficInfoHandle)
	r.POST("/getSpaceCheck", h.getSpaceCheckHandle)
	r.POST("/getTrafficCheck", h.getTrafficCheckHandle)
	r.GET("/getSpacePrice", h.getSpacePriceHandle)
	r.GET("/getTrafficPrice", h.getTrafficPriceHandle)
	r.POST("/buySpace", h.buySpaceHandle)
	r.POST("/buyTraffic", h.buyTrafficHandle)
	r.POST("/recharge", h.getApproveTsHash)
	r.POST("/getAllowance", h.getAllowanceHandle)

	r.GET("/getReceipt", h.checkReceiptHandle)

	r.GET("/cashSpace", h.cashSpaceHandle)
	r.GET("/cashTraffic", h.cashTrafficHandle)
}
