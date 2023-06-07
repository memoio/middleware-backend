package contract

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/storage"
	"github.com/memoio/contractsv2/go_contracts/erc"
)

type BuyPackage struct {
	Pkgid     uint64
	Amount    int64
	Starttime uint64
	Chainid   string
}

type PackageInfo struct {
	Time    uint64
	Kind    uint8
	Buysize *big.Int
	Amount  *big.Int
	State   uint8
}

type UserBuyPackage struct {
	Starttime uint64
	Endtime   uint64
	Kind      uint8
	Buysize   *big.Int
	Amount    *big.Int
	State     uint8
}

type Contract struct {
	contractAddr     common.Address
	endpoint         string
	gatewayAddr      common.Address
	gatewaySecretKey string
}

func NewContract(cfc map[int]config.ContractConfig) map[int]*Contract {
	res := make(map[int]*Contract)

	for chainid, cfg := range cfc {
		res[chainid] = &Contract{
			contractAddr:     common.HexToAddress(cfg.ContractAddr),
			endpoint:         cfg.Endpoint,
			gatewayAddr:      common.HexToAddress(cfg.GatewayAddr),
			gatewaySecretKey: cfg.GatewaySecretKey,
		}
	}

	return res
}

func (c *Contract) BalanceOf(ctx context.Context, addr string) (*big.Int, error) {
	res := new(big.Int)
	client, err := ethclient.DialContext(ctx, c.endpoint)
	if err != nil {
		return res, err
	}
	defer client.Close()

	erc20Ins, err := erc.NewERC20(c.contractAddr, client)
	if err != nil {
		return res, err
	}

	bal, err := erc20Ins.BalanceOf(&bind.CallOpts{
		From: c.gatewayAddr,
	}, common.HexToAddress(addr))
	if err != nil {
		return res, err
	}
	return res.Set(bal), nil
}

func (c *Contract) GetPkgSize(st storage.StorageType, address string) (storage.StorageInfo, error) {
	var out []interface{}
	err := c.CallContract(&out, "getPkgSize", common.HexToAddress(address), uint8(st))
	if err != nil {
		logger.Error(err)
		return storage.StorageInfo{}, logs.ContractError{Message: err.Error()}
	}

	available := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	free := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	used := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	files := *abi.ConvertType(out[3], new(uint64)).(*uint64)

	return storage.StorageInfo{
		Storage: st.String(),
		Buysize: available.Int64(),
		Free:    free.Int64(),
		Used:    used.Int64(),
		Files:   int(files),
	}, nil
}

func (c *Contract) StoreGetPkgInfos() ([]PackageInfo, error) {
	logger.Info("StoreGetPkgInfos:")
	var out []interface{}
	err := c.CallContract(&out, "storeGetPkgInfos")
	if err != nil {
		logger.Error(err)
		return nil, logs.ContractError{Message: err.Error()}
	}

	out0 := *abi.ConvertType(out[0], new([]PackageInfo)).(*[]PackageInfo)

	return out0, nil
} 

func (c *Contract) StoreGetBuyPkgs(address string) ([]UserBuyPackage, error) {
	logger.Info("StoreGetBuyPkgs:")
	var out []interface{}
	err := c.CallContract(&out, "storeGetBuyPkgs", common.HexToAddress(address))
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	out0 := *abi.ConvertType(out[0], new([]UserBuyPackage)).(*[]UserBuyPackage)
	return out0, nil
}

func (c *Contract) GetStoreAllSize() *big.Int {
	var out []interface{}
	err := c.CallContract(&out, "getStoreAllSize")
	if err != nil {
		logger.Error(err)
		return nil
	}

	available := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	return available
}
