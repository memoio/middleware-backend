package contract

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/memoio/backend/api"
)

func (c *Contract) Send(ctx context.Context, name, method string, args ...interface{}) (string, error) {
	logger.Info(name, args)
	return c.sendTransaction(ctx, name, method, args...)
}

func (c *Contract) StoreBuyPkg(ctx context.Context, address string, pkg api.BuyPackage) (string, error) {
	logger.Info("StoreBuyPkg:", address, pkg.Pkgid, pkg.Amount, pkg.Starttime, pkg.Chainid)
	a := big.NewInt(pkg.Amount)
	return c.Send(ctx, "proxy", "storeBuyPkg", common.HexToAddress(address), pkg.Pkgid, a, pkg.Starttime, pkg.Chainid)
}
