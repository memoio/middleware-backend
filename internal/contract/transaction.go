package contract

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/storage"
)

func (c *Contract) Send(name string, args ...interface{}) (string, error) {
	logger.Info(name, args)
	return c.sendTransaction(name, args...)
}

func (c *Contract) StoreOrderPkgExpiration(address, mid string, st storage.StorageType, size *big.Int) (string, error) {
	logger.Info("storeOrderPkgExpiration:", st, address, mid, size)
	return c.sendTransaction("storeOrderPkgExpiration", common.HexToAddress(address), mid, uint8(st), size)
}

func (c *Contract) StoreBuyPkg(ctx context.Context, address string, pkg api.BuyPackage) (string, error) {
	logger.Info("StoreBuyPkg:", address, pkg.Pkgid, pkg.Amount, pkg.Starttime, pkg.Chainid)
	a := big.NewInt(pkg.Amount)
	return c.sendTransaction("storeBuyPkg", common.HexToAddress(address), pkg.Pkgid, a, pkg.Starttime, pkg.Chainid)
}

func (c *Contract) AdminAddPkgInfo(time string, amount string, kind string, size string) (string, error) {
	logger.Info("AdminAddPkgInfo:", time, amount, kind, size)
	t, a, s := new(big.Int), new(big.Int), new(big.Int)
	t.SetString(time, 10)
	a.SetString(amount, 10)
	s.SetString(size, 10)
	k := storage.StringToStorageType(kind)
	return c.sendTransaction("adminAddPkgInfo", t.Uint64(), a, uint8(k), s)
}

func (c *Contract) StoreOrderPkg(address, mid string, st storage.StorageType, size *big.Int) (string, error) {
	logger.Info("StoreOrderPkg:", st, address, mid, size)
	return c.sendTransaction("storeOrderPkg", common.HexToAddress(address), mid, uint8(st), size)
}

func (c *Contract) StoreOrderPay(address, hash string, st storage.StorageType, amount *big.Int, size *big.Int) (string, error) {
	logger.Info("StoreOrderPay:", address, hash, amount, size)
	return c.sendTransaction("storeOrderpay", common.HexToAddress(address), hash, uint8(st), amount, size)
}

func (c *Contract) FlowOrderPay(address, hash string, st storage.StorageType, amount, size *big.Int) (string, error) {
	logger.Info("FlowOrderPay:", address, hash, amount, size)
	return c.sendTransaction("flowOrderpay", common.HexToAddress(address), hash, uint8(st), amount, size)
}
