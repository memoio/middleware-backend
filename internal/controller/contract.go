package controller

import (
	"context"
	"fmt"
	"math/big"

	"github.com/memoio/backend/internal/contract"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/storage"
)

func chainIdNotSet(chain int) error {
	return logs.ContractError{Message: fmt.Sprintf("chain %d not set", chain)}
}

type Package contract.BuyPackage

type PackageInfo struct {
	Pkgid int
	contract.PackageInfo
}

func (c *Controller) CanWrite(ctx context.Context, chain int, address string, size *big.Int) (bool, error) {
	cs, err := c.CheckStorage(ctx, chain, address, size)
	if err != nil {
		return false, err
	}
	return cs, nil
}

// storage
func (c *Controller) CheckStorage(ctx context.Context, chain int, address string, size *big.Int) (bool, error) {
	si, err := c.GetStorageInfo(ctx, chain, address)
	if err != nil {
		return false, err
	}

	logger.Debug("Avi", si.Buysize+si.Free, "Used", si.Used+size.Int64())
	return si.Buysize+si.Free > si.Used+size.Int64(), nil
}

func (c *Controller) GetStorageInfo(ctx context.Context, chain int, address string) (storage.StorageInfo, error) {
	ct, ok := c.contracts[chain]
	if !ok {
		return storage.StorageInfo{}, chainIdNotSet(chain)
	}

	si, err := ct.GetPkgSize(c.storageType, address)
	if err != nil {
		return storage.StorageInfo{}, err
	}

	cachesize, err := c.is.GetStorage(chain, address, c.storageType)
	if err != nil {
		return storage.StorageInfo{}, err
	}

	si.Used += cachesize.Int64()

	return si, nil
}

// balance
func (c *Controller) GetBalance(ctx context.Context, chain int, address string) (*big.Int, error) {
	ct, ok := c.contracts[chain]
	if !ok {
		return nil, chainIdNotSet(chain)
	}

	balance, err := ct.BalanceOf(ctx, address)
	if err != nil {
		return balance, err
	}

	value, err := c.sp.Size(chain, address, c.storageType)
	if err != nil {
		return balance, err
	}

	return balance.Sub(balance, value), nil
}

func (c *Controller) BuyPackage(chain int, address string, pkg Package) (string, error) {
	ct, ok := c.contracts[chain]
	if !ok {
		logger.Error(chainIdNotSet(chain))
		return "", chainIdNotSet(chain)
	}
	return ct.StoreBuyPkg(address, contract.BuyPackage(pkg))
}

func (c *Controller) GetPackageList(chain int) ([]PackageInfo, error) {
	ct, ok := c.contracts[chain]
	if !ok {
		return nil, chainIdNotSet(chain)
	}
	pi, err := ct.StoreGetPkgInfos()
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

func (c *Controller) GetUserBuyPackages(chain int, address string) ([]contract.UserBuyPackage, error) {
	ct, ok := c.contracts[chain]
	if !ok {
		return nil, chainIdNotSet(chain)
	}
	return ct.StoreGetBuyPkgs(address)
}

func (c *Controller) StoreOrderPkg(address string) error {
	// c.contract.StoreOrderPkg(address)
	return nil
}
