package controller

import (
	"context"
	"math/big"

	"github.com/memoio/backend/api"
)

func (c *Controller) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	balance, err := c.contract.BalanceOf(ctx, address)
	if err != nil {
		return balance, err
	}

	return balance, nil
}

func (c *Controller) GetStorageInfo(ctx context.Context, address string) (api.StorageInfo, error) {
	result := api.StorageInfo{}
	// get contract size
	csize, err := c.getPkgSize(ctx, address)
	if err != nil {
		return result, err
	}

	// get cache size

	return csize, nil
}
