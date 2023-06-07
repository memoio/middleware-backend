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
	stop        chan struct{}
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

func (c *Controller) Start() {
	go c.stratTask()
}

func (c *Controller) Stop() {
	c.stop <- struct{}{}
	time.Sleep(2 * time.Second)
}

func (c *Controller) stratTask() error {
	return c.uploadToContract()
}

func (c *Controller) uploadToContract() error {
	ticker := time.NewTicker(24 * time.Hour)

	for {
		select {
		case <-ticker.C:
			logger.Info("Upload To Contract")
			err := c.UploadStorage()
			if err != nil {
				logger.Error("upload ", err)
				return err
			}
			err = c.UploadTraffic()
			if err != nil {
				return err
			}
		case <-c.stop:
			logger.Info("tase stop")
			ticker.Stop()
			return nil
		}
	}
}

func (c *Controller) UploadStorage() error {
	scl := c.is.GetAllStorage()
	for _, sc := range scl {
		err := c.UploadAddStorage(sc)
		if err != nil {
			continue
		}

		err = c.UploadDelStorage(sc)
		if err != nil {
			continue
		}

		err = c.is.ResetStorage(sc.ChainID, sc.Address.Hex(), sc.SType)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Controller) UploadAddStorage(sc *database.StorageCheck) error {
	receipta, err := c.contracts[sc.ChainID].StoreOrderPkg(sc.Address.Hex(), sc.AddHash(), sc.SType, sc.AddSize)
	if err != nil {
		return err
	}

	err = c.contracts[sc.ChainID].CheckTrsaction(context.TODO(), receipta)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) UploadDelStorage(sc *database.StorageCheck) error {
	receiptd, err := c.contracts[sc.ChainID].StoreOrderPkgExpiration(sc.Address.Hex(), sc.DelHash(), sc.SType, sc.AddSize)
	if err != nil {
		return err
	}

	err = c.contracts[sc.ChainID].CheckTrsaction(context.TODO(), receiptd)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) UploadTraffic() error {
	pcl := c.sp.GetAllStorage()
	for _, pc := range pcl {
		receipt, err := c.contracts[pc.ChainID].StoreOrderPay(pc.Address.Hex(), pc.Hash(), pc.SType, pc.Value, pc.Size)
		if err != nil {
			continue
		}

		err = c.contracts[pc.ChainID].CheckTrsaction(context.TODO(), receipt)
		if err != nil {
			continue
		}

		err = c.sp.ResetPay(pc.ChainID, pc.Address.Hex(), pc.SType)
		if err != nil {
			return err
		}
	}
	return nil
}
