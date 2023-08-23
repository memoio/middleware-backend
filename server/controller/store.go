package controller

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/utils"
)

func (c *Controller) PutObject(ctx context.Context, address, object string, r io.Reader, opts ObjectOptions) (PutObjectResult, error) {
	result := PutObjectResult{}
	pi, err := c.SpacePayInfo(ctx, address)
	if err != nil {
		return result, err
	}

	err = c.canWrite(ctx, address, uint64(opts.Size), opts.Sign, pi)
	if err != nil {
		return result, err
	}

	err = c.changeStore(ctx, opts.Area)
	if err != nil {
		return result, err
	}

	oi, err := c.store.PutObject(ctx, address, object, r, api.ObjectOptions(opts))
	if err != nil {
		return result, err
	}

	userdefine, err := json.Marshal(oi.UserDefined)
	if err != nil {
		return result, err
	}

	fi := api.FileInfo{
		Address:    address,
		Name:       object,
		Mid:        oi.Cid,
		SType:      c.st,
		Size:       oi.Size,
		ModTime:    oi.ModTime,
		UserDefine: string(userdefine),
		UserID:     c.storeID,
	}

	err = c.storeFileInfo(ctx, fi, opts.Sign, pi.Nonce)
	if err != nil {
		c.store.DeleteObject(ctx, address, oi.Name)
		return result, err
	}

	result.Mid = oi.Cid

	return result, nil
}

func (c *Controller) GetObject(ctx context.Context, address, mid string, w io.Writer, opts ObjectOptions) (GetObjectResult, error) {
	result := GetObjectResult{}

	ob, err := c.GetObjectInfo(ctx, address, mid)
	if err != nil {
		return result, err
	}

	pi, err := c.TrafficPayInfo(ctx, address)
	if err != nil {
		return result, err
	}

	err = c.canRead(ctx, address, uint64(ob.Size), opts.Sign, pi)
	if err != nil {
		return result, err
	}

	err = c.changeStoreWithID(ctx, ob.UserID)
	if err != nil {
		return result, err
	}

	err = c.store.GetObject(ctx, mid, w, api.ObjectOptions(opts))
	if err != nil {
		return result, err
	}

	result.Name = ob.Name
	result.CType = utils.TypeByExtension(ob.Name)
	result.Size = ob.Size

	err = c.storeCacheInfo(ctx, ob, opts.Sign, pi.Nonce)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (c *Controller) GetObjectInfo(ctx context.Context, address, mid string) (api.FileInfo, error) {
	return c.getObjectInfo(ctx, address, mid, c.st)
}

func (c *Controller) ListObjects(ctx context.Context, address string) (ListObjectsResult, error) {
	return c.listobjects(ctx, address)
}

func (c *Controller) GetObjectInfoById(ctx context.Context, id int) (api.FileInfo, error) {
	return c.getObjectInfoById(ctx, id)
}

func (c *Controller) DeleteObject(ctx context.Context, address string, id int) error {
	oi, err := c.GetObjectInfoById(ctx, id)
	if err != nil {
		return err
	}

	if address != oi.Address {
		lerr := logs.ControllerError{Message: "address not right"}
		logger.Error(lerr)
		return lerr
	}

	err = c.store.DeleteObject(ctx, address, oi.Name)
	if err != nil {
		if strings.Contains(err.Error(), "not exist") {
			return c.deleteObject(ctx, id)
		}
		return err
	}

	return c.deleteObject(ctx, id)
}
