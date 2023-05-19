package controller

import (
	"os"

	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/contract"
	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/internal/gateway"
	"github.com/memoio/backend/internal/storage"
	"github.com/memoio/go-mefs-v2/lib/backend/kv"
	"github.com/memoio/go-mefs-v2/lib/backend/wrap"
)

const (
	metaStorePrefix = "meta"
)

type Controller struct {
	storageApi  gateway.IGateway
	contract    *contract.Contract
	storageType storage.StorageType
	is          *database.SendStorage
}

func NewController(path string, cf *config.Config) *Controller {
	logger.Info("new controller")
	api, ok := ApiMap[path]
	if !ok {
		logger.Error("storage api not support")
		return nil
	}

	ct := contract.NewContract(cf.Contract)

	opt := kv.DefaultOptions
	bpath := "./datastore/" + api.T.String()
	err := os.MkdirAll(bpath, os.ModePerm)
	if err != nil {
		logger.Error(err)
		return nil
	}
	ds, err := kv.NewBadgerStore(bpath, &opt)
	if err != nil {
		logger.Error(err)
		return nil
	}

	dss := wrap.NewKVStore(metaStorePrefix, ds)

	is := database.NewSender(dss)
	return &Controller{
		storageApi:  api.G,
		storageType: api.T,
		contract:    ct,
		is:          is,
	}
}
