package contract

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/wallet"
	com "github.com/memoio/contractsv2/common"
	"github.com/memoio/middleware-contracts/go-contracts/proxy"
)

func (c *Contract) BuySpace(ctx context.Context, buyer string, size uint64) (string, error) {
	return c.GetTrasaction(ctx, c.proxyAddr, buyer, "proxy", "buySpace", size, api.DurationDay, common.HexToAddress(buyer))
}

func (c *Contract) BuyTraffic(ctx context.Context, buyer string, size uint64) (string, error) {
	return c.GetTrasaction(ctx, c.proxyAddr, buyer, "proxy", "buyTraffic", size, common.HexToAddress(buyer))
}

func (c *Contract) CashSpaceCheck(ctx context.Context, check api.CheckInfo) (string, error) {
	client, err := ethclient.DialContext(ctx, c.endpoint)
	if err != nil {
		return "", err
	}
	defer client.Close()
	proxyIns, err := proxy.NewProxy(c.proxyAddr, client)
	if err != nil {
		return "", err
	}
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return "", lerr
	}
	logger.Debug("chainID: ", chainID)

	sk, err := getSk(ctx, c.seller.String())
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return "", lerr
	}
	txAuth, err := com.MakeAuth(chainID, sk)
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return "", lerr
	}
	tx, err := proxyIns.CashSpaceCheck(txAuth, check.Nonce, check.FileSize.Uint64(), api.DurationDay, check.Sign)
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return "", lerr
	}
	return tx.Hash().String(), nil
}

func (c *Contract) CashTrafficCheck(ctx context.Context, check api.CheckInfo) (string, error) {
	client, err := ethclient.DialContext(ctx, c.endpoint)
	if err != nil {
		return "", err
	}
	defer client.Close()
	proxyIns, err := proxy.NewProxy(c.proxyAddr, client)
	if err != nil {
		return "", err
	}
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return "", lerr
	}
	logger.Debug("chainID: ", chainID)

	sk, err := getSk(ctx, c.seller.String())
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return "", lerr
	}
	txAuth, err := com.MakeAuth(chainID, sk)
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return "", lerr
	}
	tx, err := proxyIns.CashTrafficCheck(txAuth, check.Nonce, check.FileSize.Uint64(), check.Sign)
	if err != nil {
		lerr := logs.ContractError{Message: err.Error()}
		logger.Error(lerr)
		return "", lerr
	}
	return tx.Hash().String(), nil
}

func getSk(ctx context.Context, sender string) (string, error) {
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

	return hex.EncodeToString(sk), nil
}
