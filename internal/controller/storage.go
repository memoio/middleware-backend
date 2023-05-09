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

func PutObject(ctx context.Context, path, address, object string, r io.Reader, opts ObjectOptions) (PutObjectResult, error) {
	result := PutObjectResult{}

	api, ok := ApiMap[path]
	if !ok {
		logger.Error("storage api not support")
		return result, logs.ControllerError{Message: "storage api not support"}
	}

	// Check if it is possible to write
	cw, err := CanWrite(ctx, api.T, address, big.NewInt(opts.Size))
	if err != nil {
		return result, err
	}

	if !cw {
		logger.Error("Insufficient space or balance")
		return result, logs.StorageError{Message: "insufficient space or balance"}
	}

	// put obejct
	oi, err := api.G.PutObject(ctx, address, object, r, gateway.ObjectOptions(opts))
	if err != nil {
		return result, err
	}
	result.Mid = oi.Cid
	return result, nil
}

func GetObject(ctx context.Context, path, mid string, w io.Writer, opts ObjectOptions) (GetObjectResult, error) {
	result := GetObjectResult{}
	api, ok := ApiMap[path]
	if !ok {
		logger.Error("storage api not support")
		return result, logs.ControllerError{Message: "storage api not support"}
	}

	obi, err := api.G.GetObjectInfo(ctx, mid)
	if err != nil {
		return result, err
	}

	err = api.G.GetObject(ctx, mid, w, gateway.ObjectOptions(opts))
	if err != nil {
		return result, err
	}

	result.Name = obi.Name
	result.CType = obi.CType
	result.Size = obi.Size

	return result, nil
}

func ListObjects(ctx context.Context, path, address string) (ListObjectsResult, error) {
	result := ListObjectsResult{}

	api, ok := ApiMap[path]
	if !ok {
		logger.Error("storage api not support")
		return result, logs.ControllerError{Message: "storage api not support"}
	}

	loi, err := api.G.ListObjects(ctx, address)
	if err != nil {
		return result, err
	}

	result.Address = address
	result.Storage = api.T.String()

	for _, oi := range loi {
		result.Objects = append(result.Objects, ObjectInfoResult{
			Name:        oi.Name,
			Size:        oi.Size,
			Mid:         oi.Cid,
			ModTime:     oi.ModTime,
			UserDefined: oi.UserDefined,
		})
	}

	return result, nil
}
