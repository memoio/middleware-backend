package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"

	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/internal/gateway"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/utils"
)

var logger = logs.Logger("controller")

type ObjectOptions gateway.ObjectOptions

func (c *Controller) PutObject(ctx context.Context, chain int, address, object string, r io.Reader, opts ObjectOptions) (PutObjectResult, error) {
	result := PutObjectResult{}

	// Check if it is possible to write
	cw, err := c.CanWrite(ctx, chain, address, big.NewInt(opts.Size))
	if err != nil {
		return result, err
	}

	if !cw {
		logger.Error("Insufficient space or balance")
		return result, logs.StorageError{Message: "insufficient space or balance"}
	}

	// put obejct
	bucket := address + fmt.Sprint(chain)
	oi, err := c.storageApi.PutObject(ctx, bucket, object, r, gateway.ObjectOptions(opts))
	if err != nil {
		return result, err
	}

	userdefine, err := json.Marshal(oi.UserDefined)
	if err != nil {
		return result, err
	}
	fi := database.FileInfo{
		ChainID:    chain,
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

	err = c.is.AddStorage(chain, address, c.storageType, big.NewInt(oi.Size), oi.Cid)
	if err != nil {
		return result, err
	}

	result.Mid = oi.Cid

	return result, nil
}

func (c *Controller) GetObject(ctx context.Context, chain int, address, mid string, w io.Writer, opts ObjectOptions) (GetObjectResult, error) {
	result := GetObjectResult{}

	obi, err := c.GetObjectInfo(ctx, chain, address, mid)
	if err != nil {
		return result, err
	}

	balance, err := c.GetBalance(ctx, chain, address)
	if err != nil {
		return result, err
	}

	trafficCost := c.cfg.Storage.TrafficCost

	needpay := big.NewInt(trafficCost)
	needpay.Mul(needpay, big.NewInt(obi.Size))

	if balance.Cmp(needpay) < 0 {
		return result, logs.ControllerError{Message: fmt.Sprintf("balance not enough, balance=%d needpay=%d", balance, needpay)}
	}

	err = c.storageApi.GetObject(ctx, mid, w, gateway.ObjectOptions(opts))
	if err != nil {
		return result, err
	}

	result.Name = obi.Name
	result.CType = utils.TypeByExtension(obi.Name)
	result.Size = obi.Size

	value := c.getPrice(result.Size)
	err = c.sp.AddPay(chain, address, c.storageType, big.NewInt(result.Size), value, obi.Mid)

	if err != nil {
		logger.Error(err)
		return result, err
	}

	return result, nil
}

func (c *Controller) DeleteObject(ctx context.Context, chain int, address, mid string) error {
	fi, err := c.GetObjectInfo(ctx, chain, address, mid)
	if err != nil {
		return err
	}

	bucket := address + fmt.Sprint(chain)
	err = c.storageApi.DeleteObject(ctx, bucket, fi.Name)
	if err != nil {
		return err
	}

	res, err := database.Delete(chain, address, mid, c.storageType)
	if err != nil {
		return err
	}

	if !res {
		lerr := logs.ConfigError{Message: "delete object failed"}
		logger.Error(lerr.Error())
		return lerr
	}

	return c.is.DelStorage(chain, address, c.storageType, big.NewInt(fi.Size), fi.Mid)
}

func (c *Controller) getPrice(size int64) *big.Int {
	price := c.cfg.Storage.TrafficCost
	p := big.NewInt(price)

	p.Mul(p, big.NewInt(size))

	return p
}
