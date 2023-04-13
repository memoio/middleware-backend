package contract

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/global"
	"github.com/memoio/contractsv2/go_contracts/erc"
)

type pkgInfo struct {
	Time    uint64
	Kind    uint8
	Buysize *big.Int
	Amount  *big.Int
	State   uint8
}

type storeInfo struct {
	Starttime uint64
	Endtime   uint64
	Kind      uint8
	Buysize   *big.Int
	Amount    *big.Int
	State     uint8
}

func BalanceOf(ctx context.Context, addr string) *big.Int {
	res := new(big.Int)
	client, err := ethclient.DialContext(ctx, global.Endpoint)
	if err != nil {
		return res
	}
	defer client.Close()

	erc20Ins, err := erc.NewERC20(global.ContractAddrV2, client)
	if err != nil {
		return res
	}

	bal, err := erc20Ins.BalanceOf(&bind.CallOpts{
		From: global.GatewayAddr,
	}, common.HexToAddress(addr))
	if err != nil {
		return res
	}
	return res.Set(bal)
}

func GetPkgSize(kind uint8, address string) (global.StorageInfo, error) {
	var out []interface{}
	log.Println(kind)
	err := CallContract(&out, "getPkgSize", global.ContractAddrV2, common.HexToAddress(address), kind)
	if err != nil {
		log.Println(err)
		return global.StorageInfo{}, err
	}

	available := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	free := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	used := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	files := *abi.ConvertType(out[3], new(uint64)).(*uint64)

	log.Println(available, free, used, files)
	return global.StorageInfo{
		Buysize: available.Int64(),
		Free:    free.Int64(),
		Used:    used.Int64(),
		Files:   int(files),
	}, nil
}

func StoreGetPkgInfos() ([]pkgInfo, error) {
	log.Println("StoreGetPkgInfos:")
	var out []interface{}
	err := CallContract(&out, "storeGetPkgInfos", global.ContractAddrV2)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	out0 := *abi.ConvertType(out[0], new([]pkgInfo)).(*[]pkgInfo)

	return out0, nil
}

func StoreOrderPkg(address, mid string, stype uint8, size *big.Int) bool {
	log.Println("StoreOrderPkg:", stype, address, mid, size)
	return sendTransaction("storage", "storeOrderPkg", global.ContractAddrV2, common.HexToAddress(address), mid, stype, size)
}

func StoreOrderPay(address, hash string, amount *big.Int, size *big.Int) bool {
	log.Println("StoreOrderPay:", address, hash, amount, size)
	return sendTransaction("pay", "storeOrderpay", global.ContractAddrV2, common.HexToAddress(address), hash, amount, size)
}

func StoreBuyPkg(address string, pkgid uint64, amount int64, starttime uint64, chainid string) bool {
	log.Println("StoreBuyPkg:", address, pkgid, amount, starttime, chainid)
	a := big.NewInt(amount)
	return sendTransaction("buy", "storeBuyPkg", global.ContractAddrV2, common.HexToAddress(address), pkgid, a, starttime, chainid)
}

func AdminAddPkgInfo(time string, amount string, kind string, size string) bool {
	log.Println("AdminAddPkgInfo:", time, amount, kind, size)
	t, a, k, s := new(big.Int), new(big.Int), new(big.Int), new(big.Int)
	t.SetString(time, 10)
	a.SetString(amount, 10)
	k.SetString(kind, 10)
	s.SetString(size, 10)
	return sendTransaction("buy", "adminAddPkgInfo", global.ContractAddrV2, t.Uint64(), a, uint8(k.Uint64()), s)
}

func StoreGetBuyPkgs(address string) ([]storeInfo, error) {
	log.Println("StoreGetBuyPkgs:")
	var out []interface{}
	err := CallContract(&out, "storeGetBuyPkgs", global.ContractAddrV2, common.HexToAddress(address))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	out0 := *abi.ConvertType(out[0], new([]storeInfo)).(*[]storeInfo)
	return out0, nil
}
