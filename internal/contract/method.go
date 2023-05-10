package contract

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/internal/logs"
)

var logger = logs.Logger("contract")

const (
	payTopic     = "0xc0e3b3bf3b856068b6537f07e399954cb5abc4fade906ee21432a8ded3c36ec8"
	storageTopic = "0x63fbca6586cb6d6fcf9fe8ab7daf3ffaf7fdad8f5d2ab29109fe71599b10d800"
	buyTopic     = "0x9393f0a0a85953b7957a62d1ced4afd964332dad208249e1db83ce254babfccc"
	delTopic     = "0xbcc253ceed59fcdc9a5bad89f7886c6c4561b5f245b4e99c7d1dea0c397807ed"
)

const (
	getPkgSizeABI              = `[{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"kind","type":"uint8"}],"name":"getPkgSize","outputs":[{"name":"used","type":"uint256"},{"name":"available","type":"uint256"},{"name":"total","type":"uint256"},{"name":"expires","type":"uint64"}],"payable":false,"stateMutability":"view","type":"function"}]`
	storeOrderPkgExpirationABI = `[{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"hashid","type":"string"},{"name":"kind","type":"uint8"},{"name":"size","type":"uint256"}],"name":"storeOrderPkgExpiration","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"hashid","type":"string"},{"indexed":false,"name":"size","type":"uint256"},{"indexed":false,"name":"nonce","type":"uint256"}],"name":"StoreOrderExpirationed","type":"event"}]`
	storeOrderPkgABI           = `[{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"hashid","type":"string"},{"name":"kind","type":"uint8"},{"name":"size","type":"uint256"}],"name":"storeOrderPkg","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"hashid","type":"string"},{"indexed":false,"name":"size","type":"uint256"},{"indexed":false,"name":"nonce","type":"uint256"}],"name":"StoreOrderPkg","type":"event"}]`
	storeOrderPayABI           = `[{"constant":false,"inputs":[{"name":"_addr","type":"address"},{"name":"_str","type":"string"},{"name":"kind","type":"uint8"},{"name":"amount","type":"uint256"},{"name":"size","type":"uint256"}],"name":"storeOrderpay","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
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
		logger.Error(err)
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

func (c *Contract) CallContract(results *[]interface{}, name string, args ...interface{}) error {
	client, err := ethclient.DialContext(context.TODO(), c.endpoint)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer client.Close()

	logger.Info("CallContract ", name)
	if results == nil {
		results = new([]interface{})
	}

	contractABI := getContractABI(name)

	encodeData, err := contractABI.Pack(name, args...)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("packed!")
	msg := ethereum.CallMsg{
		To:   &c.contractAddr,
		Data: encodeData,
	}

	result, err := client.CallContract(context.TODO(), msg, nil)
	if err != nil {
		logger.Error(err)
		return err
	}

	if len(*results) == 0 {
		res, err := contractABI.Unpack(name, result)
		*results = res
		logger.Error(err)
		return err
	}
	res := *results
	return contractABI.UnpackIntoInterface(res[0], name, result)

}

func (c *Contract) sendTransaction(trtype, name string, args ...interface{}) bool {
	logger.Info("sendTransaction")
	client, err := ethclient.DialContext(context.TODO(), c.endpoint)
	if err != nil {
		logger.Error(err)
		return false
	}
	defer client.Close()

	nonce, err := client.PendingNonceAt(context.TODO(), c.gatewayAddr)
	if err != nil {
		logger.Error(err)
		return false
	}
	logger.Debug("nonce: ", nonce)

	chainID, err := client.NetworkID(context.TODO())
	if err != nil {
		logger.Error(err)
		return false
	}
	logger.Debug("chainID: ", chainID)

	contractABI := getContractABI(name)

	data, err := contractABI.Pack(name, args...)
	if err != nil {
		logger.Error("pack error: ", err)
		return false
	}

	privateKey, err := crypto.HexToECDSA(c.gatewaySecretKey)
	if err != nil {
		logger.Error("Failed to decode private key: %v\n", err)
		return false
	}

	gasLimit := uint64(300000)
	gasPrice := big.NewInt(1000)
	tx := types.NewTransaction(nonce, c.contractAddr, big.NewInt(0), gasLimit, gasPrice, data)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return false
	}
	err = client.SendTransaction(context.TODO(), signedTx)
	if err != nil {
		logger.Errorf("Failed to send transaction: %v\n", err)
		return false
	}

	logger.Info("waiting tx complete...")
	time.Sleep(30 * time.Second)

	receipt, err := client.TransactionReceipt(context.TODO(), signedTx.Hash())
	if err != nil {
		logger.Error("receipt:", err)
		return false
	}

	return checkResult(trtype, receipt)
}

func checkResult(trtype string, receipt *types.Receipt) bool {
	var topic string
	switch trtype {
	case "pay":
		topic = payTopic
	case "storage":
		topic = storageTopic
	case "buy":
		topic = buyTopic
	case "delpkg":
		topic = delTopic
	}

	if receipt.Status != 1 {
		logger.Error("Status not right")
		logger.Error(receipt.Logs)
		logger.Error(receipt)
		return false
	}

	if len(receipt.Logs) == 0 {
		logger.Error("no logs")
		return trtype == "setpkg"
	}

	if len(receipt.Logs[0].Topics) == 0 {
		logger.Error("no topics")
		return false
	}

	if receipt.Logs[0].Topics[0].String() != topic {
		logger.Error("topic not right: ", receipt.Logs[0].Topics[0].String())
		return false
	}

	return true
}
