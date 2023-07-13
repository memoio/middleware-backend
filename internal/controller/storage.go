package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	err := c.CanWrite(ctx, chain, address, big.NewInt(opts.Size))
	if err != nil {
		return result, err
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
		Public:     opts.Public,
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

	obi, err := c.checkAccess(ctx, chain, address, mid)
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

func (c *Controller) GetObjectPublic(ctx context.Context, chain int, mid string, w io.Writer, opts ObjectOptions) (GetObjectResult, error) {
	result := GetObjectResult{}
	obi, err := c.checkAccessPublic(ctx, chain, mid)
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

func (c *Controller) DeleteObject(ctx context.Context, address string, id int) error {
	fi, err := c.GetObjectById(ctx, id)
	if err != nil {
		return err
	}

	if fi.Address != address {
		return logs.ControllerError{"no access to delete"}
	}

	bucket := address + fmt.Sprint(fi.ChainID)
	err = c.storageApi.DeleteObject(ctx, bucket, fi.Name)
	if err != nil {
		return err
	}

	err = database.Delete(fi.ChainID, address, fi.Mid, c.storageType)
	if err != nil {
		return err
	}

	return c.is.DelStorage(fi.ChainID, address, c.storageType, big.NewInt(fi.Size), fi.Mid)
}

func (c *Controller) getPrice(size int64) *big.Int {
	price := c.cfg.Storage.TrafficCost
	p := big.NewInt(price)

	p.Mul(p, big.NewInt(size))

	return p
}

func (c *Controller) canRead(ctx context.Context, address string, chain int, size int64) error {
	flowsize, err := c.GetFlowSize(ctx, chain, address)
	if err != nil {
		return err
	}

	used := flowsize.Used
	used.Add(used, big.NewInt(size))
	cachesize, err := c.sp.Size(chain, address, c.storageType)
	if err != nil {
		return err
	}

	used.Add(used, cachesize)
	if used.Cmp(flowsize.Free) > 0 {
		balance, err := c.GetBalance(ctx, chain, address)
		if err != nil {
			return nil
		}

		trafficCost := c.cfg.Storage.TrafficCost

		needpay := big.NewInt(trafficCost)
		needpay.Mul(needpay, big.NewInt(size))

		if balance.Cmp(needpay) < 0 {
			return logs.ControllerError{Message: fmt.Sprintf("balance not enough, balance=%d needpay=%d", balance, needpay)}
		}
		return nil
	}

	return nil
}

func (c *Controller) UpdateFlowSize(ctx context.Context, chain int, address, mid string, size *big.Int) error {
	value := c.getPrice(size.Int64())
	err := c.sp.AddPay(chain, address, c.storageType, size, value, mid)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (c *Controller) checkAccess(ctx context.Context, chain int, address string, mid string) (database.FileInfo, error) {
	result := database.FileInfo{}
	obi, err := c.GetObjectInfo(ctx, chain, mid)
	if err != nil {
		return result, err
	}
	for key, fi := range obi {
		if key == address {
			return fi, nil
		}
	}
	err = logs.ControllerError{Message: "no access"}
	return result, err
}

func (c *Controller) checkAccessPublic(ctx context.Context, chain int, mid string) (database.FileInfo, error) {
	result := database.FileInfo{}
	obi, err := c.GetObjectInfoPublic(ctx, chain, mid)
	if err != nil {
		return result, err
	}
	for _, fi := range obi {
		log.Println(fi)
		if fi.Public {
			return fi, nil
		}
	}
	err = logs.ControllerError{Message: "no access"}
	return result, err
}
