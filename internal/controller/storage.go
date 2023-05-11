package controller

import (
	"context"
	"encoding/json"
	"io"
	"math/big"

	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/internal/gateway"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/utils"
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

	userdefine, err := json.Marshal(oi.UserDefined)
	if err != nil {
		return result, err
	}
	fi := database.FileInfo{
		Address:    address,
		Name:       object,
		Mid:        oi.Cid,
		SType:      c.storageType,
		Size:       oi.Size,
		ModTime:    oi.ModTime,
		UserDefine: string(userdefine),
	}

	res, err := database.Put(fi)
	if err != nil || !res {
		return result, logs.StorageError{Message: "write to database error, err"}
	}

	result.Mid = oi.Cid

	return result, nil
}

func (c *Controller) GetObject(ctx context.Context, mid string, w io.Writer, opts ObjectOptions) (GetObjectResult, error) {
	result := GetObjectResult{}

	obi, err := c.GetObjectInfo(ctx, mid)
	if err != nil {
		return result, err
	}

	err = c.storageApi.GetObject(ctx, mid, w, gateway.ObjectOptions(opts))
	if err != nil {
		return result, err
	}

	result.Name = obi.Name
	result.CType = utils.TypeByExtension(obi.Name)
	result.Size = obi.Size

	return result, nil
}
