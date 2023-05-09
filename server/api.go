package server

import (
	"log"

	"github.com/memoio/backend/internal/gateway"
	"github.com/memoio/backend/internal/gateway/ipfs"
	"github.com/memoio/backend/internal/gateway/mefs"
	"github.com/memoio/backend/internal/storage"
)

var ApiMap map[string]Api

type Api struct {
	G gateway.IGateway
	T storage.StorageType
}

func init() {
	loadApiMap()
}

func loadApiMap() {
	ApiMap = make(map[string]Api)

	mefs, err := mefs.NewGateway()
	if err != nil {
		log.Println("load mefs ap failed")
		return
	}
	ApiMap["/mefs"] = Api{G: mefs, T: storage.MEFS}

	ipfs, err := ipfs.NewGateway()
	if err != nil {
		log.Println("load mefs ap failed")
		return
	}
	ApiMap["/ipfs"] = Api{G: ipfs, T: storage.IPFS}
}
