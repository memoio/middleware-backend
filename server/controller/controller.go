package controller

import (
	"context"

	"github.com/memoio/backend/api"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/contract"
	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/internal/gateway/mefs"
	"github.com/memoio/backend/internal/logs"
)

var logger = logs.Logger("controller")

type Controller struct {
	storeID  int
	st       api.StorageType
	store    api.IGateway
	contract api.IContract
	database api.IDataBase
}

func NewController(st api.StorageType, store api.IGateway) (*Controller, error) {
	contract, err := contract.NewContract(config.Cfg.Contract)
	if err != nil {
		return nil, err
	}

	database, err := database.NewDataStore(st.String())
	if err != nil {
		return nil, err
	}
	return &Controller{
		storeID:  -1,
		st:       st,
		store:    store,
		contract: contract,
		database: database,
	}, nil
}

func (c *Controller) changeStore(ctx context.Context, area string) error {
	if c.st == api.MEFS {
		ui, err := c.database.SelectUser(ctx, area)
		if err != nil {
			return err
		}
		store, err := mefs.NewGatewayWith(ui.Api, ui.Token)
		if err != nil {
			return err
		}

		c.store = store
		c.storeID = ui.ID
	}

	return nil
}

func (c *Controller) changeStoreWithID(ctx context.Context, id int) error {
	if c.st == api.MEFS {
		ui, err := c.database.GetUser(ctx, id)
		if err != nil {
			return err
		}
		store, err := mefs.NewGatewayWith(ui.Api, ui.Token)
		if err != nil {
			return err
		}

		c.store = store
		c.storeID = ui.ID
	}

	return nil
}
