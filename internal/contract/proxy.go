package contract

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

func (c *Contract) BuySpace(ctx context.Context, buyer string, size uint64) (string, error) {
	return c.GetTrasaction(ctx, c.proxyAddr, buyer, "proxy", "buySpace", size, durationDay, common.HexToAddress(buyer))
}

func (c *Contract) BuyTraffic(ctx context.Context, buyer string, size uint64) (string, error) {
	return c.GetTrasaction(ctx, c.proxyAddr, buyer, "proxy", "buyTraffic", size, common.HexToAddress(buyer))
}
