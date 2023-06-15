package controller

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/memoio/backend/internal/contract"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/storage"
)

type flowSize struct {
	Used *big.Int
	Free *big.Int
}

type userBuyPackage struct {
	Starttime uint64
	Endtime   uint64
	Kind      uint8
	Buysize   *big.Int
	Amount    *big.Int
	State     uint8
}

type Package contract.BuyPackage

type packageInfo struct {
	Time    uint64
	Kind    uint8
	Buysize *big.Int
	Amount  *big.Int
	State   uint8
}

type packageInfos struct {
	Pkgid int
	packageInfo
}

func chainIdNotSet(chain int) error {
	return logs.ContractError{Message: fmt.Sprintf("chain %d not set", chain)}
}

func (c *Controller) CanWrite(ctx context.Context, chain int, address string, size *big.Int) error {
	err := c.CheckStorage(ctx, chain, address, size)
	if err != nil {
		return err
	}
	return nil
}

// storage
func (c *Controller) CheckStorage(ctx context.Context, chain int, address string, size *big.Int) error {
	si, err := c.GetStorageInfo(ctx, chain, address)
	if err != nil {
		return err
	}

	logger.Info("Avi: ", si.Buysize+si.Free, "Used: ", si.Used+size.Int64())
	if si.Buysize+si.Free < si.Used+size.Int64() {
		err = logs.StorageError{Message: "insufficient space or balance"}
		return err
	}
	return nil
}

func (c *Controller) GetStorageInfo(ctx context.Context, chain int, address string) (storage.StorageInfo, error) {
	ct, err := c.getContract(chain)
	if err != nil {
		return storage.StorageInfo{}, err
	}

	out, err := ct.Get("getPkgSize", common.HexToAddress(address), uint8(c.storageType))
	if err != nil {
		return storage.StorageInfo{}, err
	}

	available := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	free := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	used := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	files := *abi.ConvertType(out[3], new(uint64)).(*uint64)

	si := storage.StorageInfo{
		Storage: c.storageType.String(),
		Buysize: available.Int64(),
		Free:    free.Int64(),
		Used:    used.Int64(),
		Files:   int(files),
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
	ct, err := c.getContract(chain)
	if err != nil {
		return nil, err
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
	ct, err := c.getContract(chain)
	if err != nil {
		return "", err
	}
	return ct.StoreBuyPkg(address, contract.BuyPackage(pkg))
}

func (c *Controller) GetPackageList(chain int) ([]packageInfos, error) {
	ct, err := c.getContract(chain)
	if err != nil {
		return nil, err
	}

	out, err := ct.Get("storeGetPkgInfos")
	if err != nil {
		return nil, err
	}

	result := *abi.ConvertType(out[0], new([]packageInfo)).(*[]packageInfo)

	var pl []packageInfos
	for i, p := range result {
		pl = append(pl, packageInfos{
			Pkgid:       i + 1,
			packageInfo: p,
		})
	}

	return pl, nil
}

func (c *Controller) GetUserBuyPackages(chain int, address string) ([]userBuyPackage, error) {
	ct, err := c.getContract(chain)
	if err != nil {
		return nil, err
	}

	out, err := ct.Get("storeGetBuyPkgs", common.HexToAddress(address))
	if err != nil {
		return nil, err
	}

	result := *abi.ConvertType(out[0], new([]userBuyPackage)).(*[]userBuyPackage)

	return result, nil
}

func (c *Controller) StoreOrderPkg(address string) error {
	// c.contract.StoreOrderPkg(address)
	return nil
}

func (c *Controller) CheckReceipt(ctx context.Context, chain int, hash string) error {
	ct, err := c.getContract(chain)
	if err != nil {
		return err
	}
	return ct.CheckTrsaction(ctx, hash)
}

func (c *Controller) getContract(chain int) (*contract.Contract, error) {
	ct, ok := c.contracts[chain]
	if !ok {
		return nil, chainIdNotSet(chain)
	}

	err := ct.CheckContract()
	if err != nil {
		return nil, err
	}

	return ct, nil
}

func (c *Controller) GetFlowSize(ctx context.Context, chain int, address string) (flowSize, error) {
	result := flowSize{}
	ct, err := c.getContract(chain)
	if err != nil {
		return result, err
	}

	out, err := ct.Get("flowSize", common.HexToAddress(address))
	if err != nil {
		return result, err
	}

	usesize := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	freesize := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return flowSize{
		Used: usesize,
		Free: freesize,
	}, nil
}
