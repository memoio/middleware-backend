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
	"github.com/memoio/middleware-contracts/go-contracts/control"
	"github.com/memoio/middleware-contracts/go-contracts/proxy"
)

var logger = logs.Logger("contract")

func createAbi(cabi string) abi.ABI {
	parsed, err := abi.JSON(strings.NewReader(cabi))
	if err != nil {
		logger.Error(err)
	}
	return parsed
}

func getContractABI(name string) abi.ABI {
	switch name {
	case "control":
		return createAbi(control.ControlABI)
	case "proxy":
		return createAbi(proxy.ProxyABI)
	}
	return abi.ABI{}
}

func (c *Contract) CallContract(ctx context.Context, results *[]interface{}, name, method string, args ...interface{}) error {
	client, err := ethclient.DialContext(ctx, c.endpoint)
	if err != nil {
		return err
	}
	defer client.Close()

	logger.Infof("CallContract %s %s %s %s", name, method, args, c.contractAddr)
	if results == nil {
		results = new([]interface{})
	}

	contractABI := getContractABI(name)

	encodeData, err := contractABI.Pack(method, args...)
	if err != nil {
		return err
	}

	msg := ethereum.CallMsg{
		To:   &c.proxyAddr,
		Data: encodeData,
	}

	result, err := client.CallContract(context.TODO(), msg, nil)
	if err != nil {
		return err
	}

	if len(*results) == 0 {
		res, err := contractABI.Unpack(method, result)
		*results = res
		return err
	}
	res := *results
	return contractABI.UnpackIntoInterface(res[0], method, result)
}

func (c *Contract) sendTransaction(ctx context.Context, name, method string, args ...interface{}) (string, error) {
	logger.Info("sendTransaction")
	client, err := ethclient.DialContext(ctx, c.endpoint)
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return "", lerr
	}
	defer client.Close()

	nonce, err := client.PendingNonceAt(ctx, c.gatewayAddr)
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return "", lerr
	}
	logger.Debug("nonce: ", nonce)

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return "", lerr
	}
	logger.Debug("chainID: ", chainID)

	contractABI := getContractABI(name)

	data, err := contractABI.Pack(method, args...)
	if err != nil {
		lerr := logs.ContractError{Message: fmt.Sprint("pack error: ", err)}
		logger.Error(lerr)
		return "", lerr
	}

	privateKey, err := crypto.HexToECDSA(c.gatewaySecretKey)
	if err != nil {
		lerr := logs.ContractError{Message: fmt.Sprintf("Failed to decode gateway sk: %v", err)}
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
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		lerr := logs.ContractError{Message: fmt.Sprintf("Failed to send transaction: %v\n", err)}
		logger.Error(lerr)
		return "", lerr
	}

	return signedTx.Hash().String(), nil
}

func (c *Contract) CheckTrsaction(ctx context.Context, hash string) error {
	client, err := ethclient.DialContext(ctx, c.endpoint)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer client.Close()

	signedTx := common.HexToHash(hash)

	receipt, err := client.TransactionReceipt(ctx, signedTx)
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

	logger.Info("RECEIPT: ", receipt)
	if len(receipt.Logs) != 0 {
		if len(receipt.Logs[0].Topics) == 0 {
			err := logs.ContractError{Message: "no topics"}
			logger.Error(err)
			return err
		}
	}
	return nil
}
