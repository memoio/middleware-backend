package controller

import (
	"context"
	"math/big"

	"github.com/memoio/backend/internal/contract"
	"github.com/memoio/backend/internal/storage"
)

func CanWrite(ctx context.Context, st storage.StorageType, address string, size *big.Int) (bool, error) {
	cs, err := CheckStorage(ctx, st, address, size)
	if err != nil {
		return false, err
	}
	return cs, nil
}

// storage
func CheckStorage(ctx context.Context, st storage.StorageType, address string, size *big.Int) (bool, error) {
	si, err := GetStorageInfo(ctx, st, address)
	if err != nil {
		return false, err
	}

	logger.Debug("Avi", si.Buysize+si.Free, "Used", si.Used+size.Int64())
	return si.Buysize+si.Free > si.Used+size.Int64(), nil
}

func GetStorageInfo(ctx context.Context, st storage.StorageType, address string) (storage.StorageInfo, error) {
	si, err := contract.GetPkgSize(st, address)
	if err != nil {
		return storage.StorageInfo{}, err
	}

	return si, nil
}
