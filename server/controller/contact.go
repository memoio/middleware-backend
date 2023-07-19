package controller

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/memoio/backend/api"
)

func (c *Controller) getPkgSize(ctx context.Context, address string) (api.StorageInfo, error) {
	result := api.StorageInfo{}
	out, err := c.contract.Call("getPkgSize", common.HexToAddress(address), uint8(c.st))
	if err != nil {
		return result, err
	}

	available := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	free := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	used := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	files := *abi.ConvertType(out[3], new(uint64)).(*uint64)

	si := api.StorageInfo{
		Storage: c.st.String(),
		Buysize: available.Int64(),
		Free:    free.Int64(),
		Used:    used.Int64(),
		Files:   int(files),
	}

	return si, nil
}

func (c *Controller) GetPackageList(ctx context.Context) ([]packageInfos, error) {
	out, err := c.contract.Call("storeGetPkgInfos")
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

func (c *Controller) GetFlowSize(ctx context.Context, address string) (flowSize, error) {
	result := flowSize{}

	out, err := c.contract.Call("flowSize", common.HexToAddress(address))
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

func (c *Controller) GetUserBuyPackages(ctx context.Context, address string) ([]userBuyPackage, error) {
	out, err := c.contract.Call("storeGetBuyPkgs", common.HexToAddress(address))
	if err != nil {
		return nil, err
	}

	result := *abi.ConvertType(out[0], new([]userBuyPackage)).(*[]userBuyPackage)

	return result, nil
}

func (c *Controller) CheckReceipt(ctx context.Context, receipt string) error {
	return c.contract.CheckTrsaction(ctx, receipt)
}

func (c *Controller) BuyPackage(ctx context.Context, address string, pkg api.BuyPackage) (string, error) {
	return c.contract.StoreBuyPkg(ctx, address, pkg)
}
