package contract

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/wallet"
)

func (c *Contract) Send(ctx context.Context, sender, name, method string, args ...interface{}) (string, error) {
	return c.sendTransaction(ctx, sender, name, method, args...)
}

func (c *Contract) sendTransaction(ctx context.Context, sender, name, method string, args ...interface{}) (string, error) {
	logger.Info("sendTransaction")
	client, err := ethclient.DialContext(ctx, c.endpoint)
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return "", lerr
	}
	defer client.Close()

	nonce, err := client.PendingNonceAt(ctx, common.HexToAddress(sender))
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return "", lerr
	}
	logger.Info("nonce: ", nonce)

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

	ks, err := wallet.NewKeyRepo(ksp)
	if err != nil {
		lerr := logs.ContractError{Message: fmt.Sprint("get wallet err: ", err)}
		logger.Error(lerr)
		return "", err
	}
	wl := wallet.New(ks)
	sk, err := wl.WalletExport(ctx, common.HexToAddress(sender))
	if err != nil {
		lerr := logs.ContractError{Message: fmt.Sprint("get sk err: ", err)}
		logger.Error(lerr)
		return "", err
	}
	privateKey, err := crypto.ToECDSA(sk)
	if err != nil {
		lerr := logs.ContractError{Message: fmt.Sprintf("Failed to decode gateway sk: %v", err)}
		logger.Error(lerr)
		return "", lerr
	}

	gasLimit := uint64(300000)
	gasPrice := big.NewInt(1000)
	tx := types.NewTransaction(nonce, c.proxyAddr, big.NewInt(0), gasLimit, gasPrice, data)

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
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error("receipt:", lerr)
		return lerr
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

func (c *Contract) GetTrasaction(ctx context.Context, contract common.Address, sender, name, method string, args ...interface{}) (api.Transaction, error) {
	res := api.Transaction{}
	client, err := ethclient.DialContext(ctx, c.endpoint)
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return res, lerr
	}
	defer client.Close()

	nonce, err := client.PendingNonceAt(ctx, common.HexToAddress(sender))
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return res, lerr
	}
	logger.Debug("nonce: ", nonce)

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return res, lerr
	}
	logger.Info("chainID: ", chainID)

	contractABI := getContractABI(name)

	data, err := contractABI.Pack(method, args...)
	if err != nil {
		lerr := logs.ContractError{Message: fmt.Sprint("pack error: ", err)}
		logger.Error(lerr)
		return res, lerr
	}

	gasLimit := uint64(300000)
	gasPrice := big.NewInt(1000)
	// tx := types.NewTransaction(nonce, contract, big.NewInt(0), gasLimit, gasPrice, data)
	tx := types.LegacyTx{
		Nonce:    nonce,
		To:       &contract,
		Value:    big.NewInt(0),
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	}

	return api.Transaction(tx), nil
}

func (c *Contract) SendTx(ctx context.Context, hash string) error {
	client, err := ethclient.DialContext(ctx, c.endpoint)
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return lerr
	}

	defer client.Close()

	signedTxBytes, err := hex.DecodeString(hash)
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return lerr
	}

	var signedTx = new(types.Transaction)
	err = signedTx.UnmarshalBinary(signedTxBytes)
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return lerr
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		lerr := logs.ContractError{Message: fmt.Sprintf("Failed to send transaction: %v\n", err)}
		logger.Error(lerr)
		return lerr
	}

	return nil
}
