package controller

import (
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/contract"
	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/internal/logs"
)

var logger = logs.Logger("controller")

type Controller struct {
	st       api.StorageType
	store    api.IGateway
	contract api.IContract
	database api.IDataBase
}

func NewController(st api.StorageType, store api.IGateway) *Controller {
	contract := contract.NewContract(config.Cfg.Contract)
	database := database.DataBase
	return &Controller{
		st:       st,
		store:    store,
		contract: contract,
		database: database,
	}
}
