package contract

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/global"
	"github.com/memoio/contractsv2/go_contracts/erc"
)

func BalanceOf(addr common.Address) *big.Int {
	res := new(big.Int)
	client, err := ethclient.DialContext(context.TODO(), global.Endpoint)
	if err != nil {
		return res
	}
	defer client.Close()

	erc20Ins, err := erc.NewERC20(global.ContractAddr, client)
	if err != nil {
		return res
	}

	bal, err := erc20Ins.BalanceOf(&bind.CallOpts{
		From: global.GatewayAddr,
	}, addr)
	if err != nil {
		return res
	}
	return res.Set(bal)
}

func GetPkgSize(address string) (global.StorageInfo, error) {
	result, err := CallContract("getPkgSize", common.HexToAddress(address))
	if err != nil {
		return global.StorageInfo{}, err
	}

	if len(result) != 128 {
		return global.StorageInfo{}, fmt.Errorf("result not right %d", len(result))
	}

	available := new(big.Int)
	available.SetBytes(result[0:32])
	free := new(big.Int)
	free.SetBytes(result[32:64])
	used := new(big.Int)
	used.SetBytes(result[64:96])
	files := new(big.Int)
	files.SetBytes(result[96:])

	log.Println(available, free, used, files)
	return global.StorageInfo{
		Available: available.String(),
		Free:      free.String(),
		Used:      used.String(),
		Files:     files.String(),
	}, nil
}

func CallContract(name string, args interface{}) ([]byte, error) {
	client, err := ethclient.DialContext(context.TODO(), global.Endpoint)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	contractABI := getContractABI(name)

	encodeData, err := contractABI.Pack(name, args)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &global.ContractAddr,
		Data: encodeData,
	}
	result, err := client.CallContract(context.TODO(), msg, nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getContractABI(name string) abi.ABI {
	switch name {
	case "getPkgSize":
		return createAbi(getPkgSizeAbi)
	}

	return abi.ABI{}
}
