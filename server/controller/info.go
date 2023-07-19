package controller

import (
	"context"
	"math/big"
)

func (c *Controller) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	balance, err := c.contract.BalanceOf(ctx, address)
	if err != nil {
		return balance, err
	}

	return balance, nil
}
