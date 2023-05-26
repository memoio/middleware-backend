package controller

import (
	"context"
	"math/big"

	"github.com/memoio/backend/internal/contract"
	"github.com/memoio/backend/internal/storage"
)

type Package contract.BuyPackage

type PackageInfo struct {
	Pkgid int
	contract.PackageInfo
}

func (c *Controller) CanWrite(ctx context.Context, address string, size *big.Int) (bool, error) {
	cs, err := c.CheckStorage(ctx, address, size)
	if err != nil {
		return false, err
	}
	return cs, nil
}

// storage
func (c *Controller) CheckStorage(ctx context.Context, address string, size *big.Int) (bool, error) {
	si, err := c.GetStorageInfo(ctx, address)
	if err != nil {
		return false, err
	}

	logger.Debug("Avi", si.Buysize+si.Free, "Used", si.Used+size.Int64())
	return si.Buysize+si.Free > si.Used+size.Int64(), nil
}

func (c *Controller) GetStorageInfo(ctx context.Context, address string) (storage.StorageInfo, error) {
	si, err := c.contract.GetPkgSize(c.storageType, address)
	if err != nil {
		return storage.StorageInfo{}, err
	}

	cachesize, err := c.is.GetStorage(address, c.storageType)
	if err != nil {
		return storage.StorageInfo{}, err
	}

	si.Used += cachesize.Int64()

	return si, nil
}

// balance
func (c *Controller) GetBalance(ctx context.Context, address string) (*big.Int, error) {
	balance, err := c.contract.BalanceOf(ctx, address)
	if err != nil {
		return balance, err
	}

	value, err := c.sp.Size(address, c.storageType)
	if err != nil {
		return balance, err
	}

	return balance.Sub(balance, value), nil
}

func (c *Controller) BuyPackage(address string, pkg Package) bool {
	return c.contract.StoreBuyPkg(address, contract.BuyPackage(pkg))
}

func (c *Controller) GetPackageList() ([]PackageInfo, error) {
	pi, err := c.contract.StoreGetPkgInfos()
	if err != nil {
		return nil, err
	}
	var pl []PackageInfo
	for i, p := range pi {
		pl = append(pl, PackageInfo{
			Pkgid:       i + 1,
			PackageInfo: p,
		})
	}

	return pl, nil
}

func (c *Controller) GetUserBuyPackages(address string) ([]contract.UserBuyPackage, error) {
	return c.contract.StoreGetBuyPkgs(address)
}

func (c *Controller) StoreOrderPkg(address string) error {
	// c.contract.StoreOrderPkg(address)
	return nil
}
