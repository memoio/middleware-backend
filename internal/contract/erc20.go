package contract

import (
	"context"
	"math/big"
)

func (c *Contract) Approve(ctx context.Context, sender string, buyValue *big.Int) (string, error) {
	return c.GetTrasaction(ctx, c.erc20, sender, "erc20", "approve", c.storeAddr, buyValue)
}
