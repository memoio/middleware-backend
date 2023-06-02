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

func (c *Contract) StoreOrderPkg(address, mid string, st storage.StorageType, size *big.Int) bool {
	logger.Info("StoreOrderPkg:", st, address, mid, size)
	return c.sendTransaction("storage", "storeOrderPkg", common.HexToAddress(address), mid, uint8(st), size)
}

func (c *Contract) StoreOrderPay(address, hash string, st storage.StorageType, amount *big.Int, size *big.Int) bool {
	logger.Info("StoreOrderPay:", address, hash, amount, size)
	return c.sendTransaction("pay", "storeOrderpay", common.HexToAddress(address), hash, uint8(st), amount, size)
}

func (c *Contract) StoreBuyPkg(address string, pkg BuyPackage) bool {
	logger.Info("StoreBuyPkg:", address, pkg.Pkgid, pkg.Amount, pkg.Starttime, pkg.Chainid)
	a := big.NewInt(pkg.Amount)
	return c.sendTransaction("buy", "storeBuyPkg", common.HexToAddress(address), pkg.Pkgid, a, pkg.Starttime, pkg.Chainid)
}

func (c *Contract) AdminAddPkgInfo(time string, amount string, kind string, size string) bool {
	logger.Info("AdminAddPkgInfo:", time, amount, kind, size)
	t, a, s := new(big.Int), new(big.Int), new(big.Int)
	t.SetString(time, 10)
	a.SetString(amount, 10)
	s.SetString(size, 10)
	k := storage.StringToStorageType(kind)
	return c.sendTransaction("setpkg", "adminAddPkgInfo", t.Uint64(), a, uint8(k), s)
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

func (c *Contract) StoreOrderPkgExpiration(address, mid string, st storage.StorageType, size *big.Int) bool {
	logger.Info("storeOrderPkgExpiration:", st, address, mid, size)
	return c.sendTransaction("delpkg", "storeOrderPkgExpiration", common.HexToAddress(address), mid, uint8(st), size)
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
