package controller

import (
	"context"
	"io"

	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/utils"
)

func (c *Controller) PutObject(ctx context.Context, address, object string, r io.Reader, opts ObjectOptions) (PutObjectResult, error) {
	result := PutObjectResult{}
	err := canWrite()
	if err != nil {
		return result, err
	}

	oi, err := c.store.PutObject(ctx, address, object, r, api.ObjectOptions(opts))
	if err != nil {
		return result, err
	}

	err = storeFileInfo()
	if err != nil {
		c.store.DeleteObject(ctx, address, oi.Cid)
		return result, err
	}

	result.Mid = oi.Cid

	return result, nil
}

func (c *Controller) GetObject(ctx context.Context, address, mid string, w io.Writer, opts ObjectOptions) (GetObjectResult, error) {
	result := GetObjectResult{}

	err := canRead()
	if err != nil {
		return result, err
	}

	ob, err := c.GetObjectInfo(ctx, address, mid)
	if err != nil {
		return result, err
	}

	err = c.store.GetObject(ctx, mid, w, api.ObjectOptions(opts))
	if err != nil {
		return result, err
	}

	err = storeFlowSize()
	if err != nil {
		return result, err
	}

	result.Name = ob.Name
	result.CType = utils.TypeByExtension(ob.Name)
	result.Size = ob.Size

	return result, nil
}

func (c *Controller) GetObjectInfo(ctx context.Context, address, mid string) (database.FileInfo, error) {
	return c.getObjectInfo(ctx, address, mid, c.st)
}

func (c *Controller) ListObjects(ctx context.Context, address string) (ListObjectsResult, error) {
	return c.listobjects(ctx, address)
}

func (c *Controller) GetObjectInfoById(ctx context.Context, id int) (database.FileInfo, error) {
	return c.getObjectInfoById(ctx, id)
}

func (c *Controller) DeleteObject(ctx context.Context, address, mid string) error {
	oi, err := c.getObjectInfo(ctx, address, mid, c.st)
	if err != nil {
		return err
	}

	err = c.store.DeleteObject(ctx, address, oi.Name)
	if err != nil {
		return err
	}

	return c.deleteObject(ctx, address, mid)
}
