package routes

import (
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/gateway/ipfs"
	"github.com/memoio/backend/internal/gateway/mefs"
	"github.com/memoio/backend/server/controller"
)

func init() {
	handlerMap = make(map[string]handler)

	initMefs()
	initIpfs()
}

func initMefs() {
	store, err := mefs.NewGateway()
	if err != nil {
		logger.Error("init mefs error:", err)
	}

	control := controller.NewController(api.MEFS, store)
	handlerMap["mefs"] = handler{controller: control}
}

func initIpfs() {
	store, err := ipfs.NewGateway()
	if err != nil {
		logger.Error("init ipfs error:", err)
	}
	control := controller.NewController(api.IPFS, store)
	handlerMap["ipfs"] = handler{controller: control}
}
