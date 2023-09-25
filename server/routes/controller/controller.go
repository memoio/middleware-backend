package controller

import (
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/contract"
	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/internal/datastore"
	"github.com/memoio/backend/internal/kzg"
	"github.com/memoio/backend/internal/logs"
)

var logger = logs.Logger("controller")

type Controller struct {
	store     api.IGateway
	contract  api.IContract
	database  api.IDataBase
	datastore api.IDataStore
	publickey api.IPublicKey
}

func NewController() (*Controller, error) {
	contract, err := contract.NewContract(config.Cfg.Contract)
	if err != nil {
		return nil, err
	}

	database := database.NewDataBase()

	datastore, err := datastore.NewDataStore()
	if err != nil {
		return nil, err
	}
	publickey, err := kzg.NewKzg()
	if err != nil {
		return nil, err
	}
	return &Controller{
		contract:  contract,
		database:  database,
		datastore: datastore, // datastore used by cashcheck
		publickey: publickey,
	}, nil
}
