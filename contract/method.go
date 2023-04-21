package contract

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/global"
)

const (
	getPkgSizeABI              = `[{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"kind","type":"uint8"}],"name":"getPkgSize","outputs":[{"name":"used","type":"uint256"},{"name":"available","type":"uint256"},{"name":"total","type":"uint256"},{"name":"expires","type":"uint64"}],"payable":false,"stateMutability":"view","type":"function"}]`
	storeOrderPkgExpirationABI = `[{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"hashid","type":"string"},{"name":"kind","type":"uint8"},{"name":"size","type":"uint256"}],"name":"storeOrderPkgExpiration","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"hashid","type":"string"},{"indexed":false,"name":"size","type":"uint256"},{"indexed":false,"name":"nonce","type":"uint256"}],"name":"StoreOrderExpirationed","type":"event"}]`
	storeOrderPkgABI           = `[{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"hashid","type":"string"},{"name":"kind","type":"uint8"},{"name":"size","type":"uint256"}],"name":"storeOrderPkg","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"hashid","type":"string"},{"indexed":false,"name":"size","type":"uint256"},{"indexed":false,"name":"nonce","type":"uint256"}],"name":"StoreOrderPkg","type":"event"}]`
	storeOrderPayABI           = `[{"constant":false,"inputs":[{"name":"_addr","type":"address"},{"name":"_str","type":"string"},{"name":"_uint1","type":"uint256"},{"name":"_uint2","type":"uint256"}],"name":"storeOrderpay","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
	buyPkgABI                  = `[{"constant":false,"inputs":[{"name":"pkgId","type":"uint64"},{"name":"amount","type":"uint256"},{"name":"starttime","type":"uint64"}],"name":"buyPkg","outputs":[],"payable":true,"stateMutability":"payable","type":"function"}]`
	storeBuyPkgABI             = `[{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"pkgId","type":"uint64"},{"name":"amount","type":"uint256"},{"name":"starttime","type":"uint64"},{"name":"chainId","type":"string"}],"name":"storeBuyPkg","outputs":[],"payable":true,"stateMutability":"payable","type":"function"}]`
	storeGetPkgInfosABI        = `[{"constant":true,"inputs":[],"name":"storeGetPkgInfos","outputs":[{"components":[{"name":"time","type":"uint64"},{"name":"kind","type":"uint8"},{"name":"buysize","type":"uint256"},{"name":"amount","type":"uint256"},{"name":"state","type":"uint8"}],"name":"","type":"tuple[]"}],"payable":false,"stateMutability":"view","type":"function"}]`
	adminAddPkgInfoABI         = `[{"constant":false,"inputs":[{"name":"time","type":"uint64"},{"name":"amount","type":"uint256"},{"name":"kind","type":"uint8"},{"name":"buysize","type":"uint256"}],"name":"adminAddPkgInfo","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
	storeGetBuyPkgsABI         = `[{"constant":false,"inputs":[{"name":"to","type":"address"}],"name":"storeGetBuyPkgs","outputs":[{"components":[{"name":"starttime","type":"uint64"},{"name":"endtime","type":"uint64"},{"name":"kind","type":"uint8"},{"name":"buysize","type":"uint256"},{"name":"amount","type":"uint256"},{"name":"state","type":"uint8"}],"name":"","type":"tuple[]"}],"payable":false,"stateMutability":"view","type":"function"}]`
	getStoreAllSizeABI         = `[{"constant": true,"inputs": [], "name": "getStoreAllSize","outputs": [{"name": "","type": "uint256"}],"payable": false,"stateMutability": "view","type": "function"}]`
)

func createAbi(cabi string) abi.ABI {
	parsed, err := abi.JSON(strings.NewReader(cabi))
	if err != nil {
		fmt.Println(err)
	}
	return parsed
}

func getContractABI(name string) abi.ABI {
	switch name {
	case "getPkgSize":
		return createAbi(getPkgSizeABI)
	case "storeOrderPkg":
		return createAbi(storeOrderPkgABI)
	case "storeOrderpay":
		return createAbi(storeOrderPayABI)
	case "buyPkg":
		return createAbi(buyPkgABI)
	case "storeBuyPkg":
		return createAbi(storeBuyPkgABI)
	case "storeGetPkgInfos":
		return createAbi(storeGetPkgInfosABI)
	case "adminAddPkgInfo":
		return createAbi(adminAddPkgInfoABI)
	case "storeGetBuyPkgs":
		return createAbi(storeGetBuyPkgsABI)
	case "storeOrderPkgExpiration":
		return createAbi(storeOrderPkgExpirationABI)
	case "getStoreAllSize":
		return createAbi(getStoreAllSizeABI)
	}

	return abi.ABI{}
}

func CallContract(results *[]interface{}, name string, contract common.Address, args ...interface{}) error {
	client, err := ethclient.DialContext(context.TODO(), global.Endpoint)
	if err != nil {
		return err
	}
	defer client.Close()

	log.Println("connected eth")
	if results == nil {
		results = new([]interface{})
	}

	contractABI := getContractABI(name)

	encodeData, err := contractABI.Pack(name, args...)
	if err != nil {
		return err
	}
	log.Println("packed!")
	msg := ethereum.CallMsg{
		To:   &contract,
		Data: encodeData,
	}

	result, err := client.CallContract(context.TODO(), msg, nil)
	if err != nil {
		return err
	}

	if len(*results) == 0 {
		res, err := contractABI.Unpack(name, result)
		*results = res
		return err
	}
	res := *results
	return contractABI.UnpackIntoInterface(res[0], name, result)

}

func sendTransaction(trtype, name string, contract common.Address, args ...interface{}) bool {
	log.Println("sendTransaction")
	client, err := ethclient.DialContext(context.TODO(), global.Endpoint)
	if err != nil {
		log.Println("sendt error: ", err)
		return false
	}
	defer client.Close()

	log.Println(args...)
	nonce, err := client.PendingNonceAt(context.TODO(), global.GatewayAddr)
	if err != nil {
		return false
	}
	log.Println("nonce: ", nonce)

	chainID, err := client.NetworkID(context.TODO())
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("chainID: ", chainID)

	contractABI := getContractABI(name)

	data, err := contractABI.Pack(name, args...)
	if err != nil {
		log.Println("pack error: ", err)
		return false
	}

	privateKey, err := crypto.HexToECDSA(global.GatewaySecretKey)
	if err != nil {
		log.Printf("Failed to decode private key: %v\n", err)
		return false
	}

	gasLimit := uint64(300000)
	gasPrice := big.NewInt(1000)
	tx := types.NewTransaction(nonce, contract, big.NewInt(0), gasLimit, gasPrice, data)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return false
	}
	err = client.SendTransaction(context.TODO(), signedTx)
	if err != nil {
		log.Printf("Failed to send transaction: %v\n", err)
		return false
	}

	log.Println("waiting tx complete...")
	time.Sleep(30 * time.Second)

	receipt, err := client.TransactionReceipt(context.TODO(), signedTx.Hash())
	if err != nil {
		log.Println("receipt:", err)
		return false
	}

	return checkResult(trtype, receipt)
}

func checkResult(trtype string, receipt *types.Receipt) bool {
	var topic string
	switch trtype {
	case "pay":
		topic = global.PayTopic
	case "storage":
		topic = global.StorageTopic
	case "buy":
		topic = global.BuyTopic
	case "delpkg":
		topic = global.DelTopic
	}

	if receipt.Status != 1 {
		log.Println("Status not right")
		log.Println(receipt.Logs)
		log.Println(receipt)
		return false
	}

	if len(receipt.Logs) == 0 {
		log.Println("no logs")
		return trtype == "setpkg"
	}

	if len(receipt.Logs[0].Topics) == 0 {
		log.Println("no topics")
		return false
	}

	if receipt.Logs[0].Topics[0].String() != topic {
		log.Println("topic not right: ", receipt.Logs[0].Topics[0].String())
		return false
	}

	return true
}
