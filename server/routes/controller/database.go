package controller

import (
	"context"

	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/gateway/mefs"
	"github.com/memoio/backend/internal/logs"
)

func (c *Controller) getObjectInfoById(ctx context.Context, id int) (api.FileInfo, error) {
	result := api.FileInfo{}
	oi, err := c.database.GetObjectInfoById(ctx, id)
	if err != nil {
		return result, err
	}

	fi := oi.(api.FileInfo)
	if fi == result {
		lerr := logs.ControllerError{Message: "file not exist"}
		logger.Error(lerr)
		return result, lerr
	}
	return fi, nil
}

func (c *Controller) changeStore(ctx context.Context, area string) error {
	st := c.store.GetStoreType(ctx)
	if st == api.MEFS {
		ui, err := c.database.SelectUser(ctx, area)
		if err != nil {
			return err
		}
		store, err := mefs.NewGatewayWith(ui)
		if err != nil {
			return err
		}

		c.store = store
	}

	return nil
}

func (c *Controller) storeFileInfo(ctx context.Context, fi api.FileInfo, ci api.CheckInfo) error {
	err := c.datastore.Upload(ctx, ci)
	if err != nil {
		return err
	}
	return c.database.PutObject(ctx, fi)
}

func (c *Controller) getObjectInfo(ctx context.Context, address, mid string) (api.FileInfo, error) {
	result := api.FileInfo{}
	st := c.store.GetStoreType(ctx)
	oi, err := c.database.GetObjectInfo(ctx, address, mid, st)
	if err != nil {
		return result, err
	}

	fi := oi.(api.FileInfo)
	if fi == result {
		lerr := logs.DataBaseError{Message: "file not exist"}
		logger.Error(lerr)
		return result, lerr
	}
	return fi, nil
}
