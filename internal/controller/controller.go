package controller

import (
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/contract"
	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/internal/gateway"
	"github.com/memoio/backend/internal/storage"
)

type Controller struct {
	write       *database.WriteCheck
	storageApi  gateway.IGateway
	contract    *contract.Contract
	storageType storage.StorageType
}

func NewController(path string, cf *config.Config) *Controller {
	api, ok := ApiMap[path]
	if !ok {
		logger.Error("storage api not support")
		return nil
	}

	ct := contract.NewContract(cf.Contract)

	return &Controller{
		write:       database.NewWriteCheck(),
		storageApi:  api.G,
		storageType: api.T,
		contract:    ct,
	}
}
