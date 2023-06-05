package controller

import (
	"context"
	"math/big"
	"os"
	"time"

	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/contract"
	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/internal/gateway"
	"github.com/memoio/backend/internal/storage"
	"github.com/memoio/go-mefs-v2/lib/backend/kv"
	"github.com/memoio/go-mefs-v2/lib/backend/wrap"
)

const (
	metaStorePrefix = "meta"
)

type Controller struct {
	storageApi  gateway.IGateway
	contracts   map[int]*contract.Contract
	storageType storage.StorageType
	cfg         *config.Config
	is          *database.SendStorage
	sp          *database.SendPay
	download    map[string]*big.Int
}

func NewController(path string, cfg *config.Config) *Controller {
	logger.Info("new controller")
	api, ok := ApiMap[path]
	if !ok {
		logger.Error("storage api not support")
		return nil
	}

	ct := contract.NewContract(cfg.Contract)

	opt := kv.DefaultOptions
	bpath := "./datastore/" + api.T.String()
	err := os.MkdirAll(bpath, os.ModePerm)
	if err != nil {
		logger.Error(err)
		return nil
	}
	ds, err := kv.NewBadgerStore(bpath, &opt)
	if err != nil {
		logger.Error(err)
		return nil
	}

	dss := wrap.NewKVStore(metaStorePrefix, ds)

	is := database.NewSender(dss)
	sp := database.NewSenderPay(dss)

	dw := make(map[string]*big.Int)
	return &Controller{
		storageApi:  api.G,
		storageType: api.T,
		contracts:   ct,
		is:          is,
		sp:          sp,
		cfg:         cfg,
		download:    dw,
	}
}

func (c *Controller) UploadToContract() error {
	ticker := time.NewTicker(24 * time.Hour)

	for range ticker.C {
		logger.Info("Upload To Contract")

		scl := c.is.GetAllStorage()
		for _, sc := range scl {

			receipta, err := c.contracts[sc.ChainID].StoreOrderPkg(sc.Address.Hex(), sc.AddHash(), sc.SType, sc.AddSize)
			if err != nil {
				logger.Error(err)
				continue
			}
			receiptd, err := c.contracts[sc.ChainID].StoreOrderPkgExpiration(sc.Address.Hex(), sc.DelHash(), sc.SType, sc.AddSize)
			if err != nil {
				logger.Error(err)
				continue
			}

			add, err := c.contracts[sc.ChainID].CheckTrsaction(context.TODO(), receipta)
			if err != nil {
				logger.Error(err)
				continue
			}

			del, err := c.contracts[sc.ChainID].CheckTrsaction(context.TODO(), receiptd)
			if err != nil {
				logger.Error(err)
				continue
			}

			if add && del {
				err := c.is.ResetStorage(sc.ChainID, sc.Address.Hex(), sc.SType)
				if err != nil {
					logger.Error(err)
					return err
				}
			}
		}

		pcl := c.sp.GetAllStorage()
		for _, pc := range pcl {
			receipt, err := c.contracts[pc.ChainID].StoreOrderPay(pc.Address.Hex(), pc.Hash(), pc.SType, pc.Value, pc.Size)
			if err != nil {
				logger.Error(err)
				continue
			}

			res, err := c.contracts[pc.ChainID].CheckTrsaction(context.TODO(), receipt)
			if err != nil {
				logger.Error(err)
				continue
			}

			if res {
				err := c.sp.ResetPay(pc.ChainID, pc.Address.Hex(), pc.SType)
				if err != nil {
					logger.Error(err)
					return err
				}
			}
		}

	}

	return nil
}
