package contract

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/internal/logs"
)

var logger = logs.Logger("contract")

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
		return err
	}
	defer client.Close()
	// logger.Info(c.contractAddr, c.endpoint, c.gatewayAddr, c.gatewaySecretKey)
	logger.Info("CallContract ", name)
	if results == nil {
		results = new([]interface{})
	}

	contractABI := getContractABI(name)

	encodeData, err := contractABI.Pack(name, args...)
	if err != nil {
		return err
	}

	logger.Info("packed!")
	msg := ethereum.CallMsg{
		To:   &c.contractAddr,
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

func (c *Contract) sendTransaction(trtype, name string, args ...interface{}) (string, error) {
	logger.Info("sendTransaction")
	client, err := ethclient.DialContext(context.TODO(), c.endpoint)
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return "", lerr
	}
	defer client.Close()

	nonce, err := client.PendingNonceAt(context.TODO(), c.gatewayAddr)
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return "", lerr
	}
	logger.Debug("nonce: ", nonce)

	chainID, err := client.NetworkID(context.TODO())
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return "", lerr
	}
	logger.Debug("chainID: ", chainID)

	contractABI := getContractABI(name)

	data, err := contractABI.Pack(name, args...)
	if err != nil {
		lerr := logs.ContractError{Message: fmt.Sprint("pack error: ", err)}
		logger.Error(lerr)
		return "", lerr
	}

	privateKey, err := crypto.HexToECDSA(c.gatewaySecretKey)
	if err != nil {
		lerr := logs.ContractError{Message: fmt.Sprintf("Failed to decode private key: %v", err)}
		logger.Error(lerr)
		return "", lerr
	}

	gasLimit := uint64(300000)
	gasPrice := big.NewInt(1000)
	tx := types.NewTransaction(nonce, c.contractAddr, big.NewInt(0), gasLimit, gasPrice, data)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		lerr := logs.ContractError{Message: fmt.Sprint("Failed to dSignTx", err)}
		logger.Error(lerr)
		return "", lerr
	}
	err = client.SendTransaction(context.TODO(), signedTx)
	if err != nil {
		lerr := logs.ContractError{Message: fmt.Sprintf("Failed to send transaction: %v\n", err)}
		logger.Error(lerr)
		return "", lerr
	}

	return signedTx.Hash().String(), nil
}

func (c *Contract) CheckTrsaction(ctx context.Context, hash string) error {
	client, err := ethclient.DialContext(context.TODO(), c.endpoint)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer client.Close()

	signedTx := common.HexToHash(hash)

	receipt, err := client.TransactionReceipt(context.TODO(), signedTx)
	if err != nil {
		logger.Error("receipt:", err)
		return err
	}

	return checkResult(receipt)
}

func checkResult(receipt *types.Receipt) error {
	if receipt.Status != 1 {
		err := logs.ContractError{Message: "Status not right"}
		logger.Error(err)
		logger.Error(receipt.Logs)
		logger.Error(receipt)
		return err
	}

	logger.Debug("RECEIPT: ", receipt)

	if len(receipt.Logs[0].Topics) == 0 {
		err := logs.ContractError{Message: "no topics"}
		logger.Error(err)
		return err
	}

	return nil
}
