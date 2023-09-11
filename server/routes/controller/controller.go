package controller

import (
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/contract"
	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/internal/datastore"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/server/routes/controller/kzg"
)

var logger = logs.Logger("controller")

type Controller struct {
	store     api.IGateway
	contract  api.IContract
	database  api.IDataBase
	datastore api.IDataStore
	publickey api.IPublicKey
}

func NewController(store api.IGateway, path string) (*Controller, error) {
	contract, err := contract.NewContract(config.Cfg.Contract)
	if err != nil {
		return nil, err
	}

	database := database.NewDataBase()

	datastore, err := datastore.NewDataStore(path)
	if err != nil {
		return nil, err
	}
	publickey, err := kzg.GenKey()
	if err != nil {
		return nil, err
	}
	return &Controller{
		contract:  contract,
		database:  database,
		datastore: datastore,
		publickey: publickey,
		store:     store,
	}, nil
}
