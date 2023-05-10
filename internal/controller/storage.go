package controller

import (
	"context"
	"io"
	"math/big"

	"github.com/memoio/backend/internal/gateway"
	"github.com/memoio/backend/internal/logs"
)

var logger = logs.Logger("controller")

type ObjectOptions gateway.ObjectOptions

func (c *Controller) PutObject(ctx context.Context, address, object string, r io.Reader, opts ObjectOptions) (PutObjectResult, error) {
	result := PutObjectResult{}

	// Check if it is possible to write
	cw, err := c.CanWrite(ctx, address, big.NewInt(opts.Size))
	if err != nil {
		return result, err
	}

	if !cw {
		logger.Error("Insufficient space or balance")
		return result, logs.StorageError{Message: "insufficient space or balance"}
	}

	// put obejct
	oi, err := c.storageApi.PutObject(ctx, address, object, r, gateway.ObjectOptions(opts))
	if err != nil {
		return result, err
	}
	result.Mid = oi.Cid

	return result, nil
}

func (c *Controller) GetObject(ctx context.Context, mid string, w io.Writer, opts ObjectOptions) (GetObjectResult, error) {
	result := GetObjectResult{}

	obi, err := c.storageApi.GetObjectInfo(ctx, mid)
	if err != nil {
		return result, err
	}

	err = c.storageApi.GetObject(ctx, mid, w, gateway.ObjectOptions(opts))
	if err != nil {
		return result, err
	}

	result.Name = obi.Name
	result.CType = obi.CType
	result.Size = obi.Size

	return result, nil
}


