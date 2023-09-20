package datastore

import (
	"context"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
)

var _ api.IDataStore = (*DataStore)(nil)

var logger = logs.Logger("datastore")

type DataStore struct {
	*CashCheck
}

func NewDataStore() (*DataStore, error) {
	res := &DataStore{}

	opt := DefaultOptions
	bpath := "./datastore/"
	err := os.MkdirAll(bpath, os.ModePerm)
	if err != nil {
		logger.Error(err)
		return res, err
	}
	ds, err := NewBadgerStore(bpath, &opt)
	if err != nil {
		logger.Error(err)
		return res, err
	}

	// create paycheck with ds
	cp := NewCheckPay(ds)
	return &DataStore{cp}, nil
}

func (d *DataStore) Upload(ctx context.Context, info api.CheckInfo) error {
	return d.check(ctx, SPACE, info)
}

func (d *DataStore) Download(ctx context.Context, info api.CheckInfo) error {
	return d.check(ctx, TRAFFIC, info)
}

func (d *DataStore) GetSpaceInfo(ctx context.Context, buyer string) (api.CheckInfo, error) {
	return d.getCheck(ctx, SPACE, common.HexToAddress(buyer))
}

func (d *DataStore) GetTrafficInfo(ctx context.Context, buyer string) (api.CheckInfo, error) {
	return d.getCheck(ctx, TRAFFIC, common.HexToAddress(buyer))
}
