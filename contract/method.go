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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/global"
)

const (
	getPkgSizeABI    = `[{"constant":true,"inputs":[{"name":"to","type":"address"}],"name":"getPkgSize","outputs":[{"name":"used","type":"uint256"},{"name":"available","type":"uint256"},{"name":"total","type":"uint256"},{"name":"expires","type":"uint64"}],"payable":false,"stateMutability":"view","type":"function"}]`
	storeOrderPkgABI = `[{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"hashid","type":"string"},{"name":"size","type":"uint256"}],"name":"storeOrderPkg","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"hashid","type":"string"},{"indexed":false,"name":"size","type":"uint256"},{"indexed":false,"name":"nonce","type":"uint256"}],"name":"StoreOrderPkg","type":"event"}]`
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
	}

	return abi.ABI{}
}

func CallContract(name string, args ...interface{}) ([]byte, error) {
	client, err := ethclient.DialContext(context.TODO(), global.Endpoint)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	contractABI := getContractABI(name)

	encodeData, err := contractABI.Pack(name, args...)
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

func sendTransaction(trtype, name string, args ...interface{}) bool {
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
	tx := types.NewTransaction(nonce, global.ContractAddr, big.NewInt(0), gasLimit, gasPrice, data)

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
		log.Println(err)
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
	}

	if receipt.Status != 1 {
		log.Println("Status not right")
		log.Println(receipt.Logs)
		log.Println(receipt)
		return false
	}

	if len(receipt.Logs) == 0 {
		log.Println("no logs")
		return false
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
