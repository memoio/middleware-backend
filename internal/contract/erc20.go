package contract

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/contractsv2/go_contracts/erc"
)

func (c *Contract) Approve(ctx context.Context, at, sender string, buyValue *big.Int) (string, error) {
	if at == "store" {
		return c.GetTrasaction(ctx, c.erc20, sender, "erc20", "approve", c.storeAddr, buyValue)
	}
	if at == "read" {
		return c.GetTrasaction(ctx, c.erc20, sender, "erc20", "approve", c.readAddr, buyValue)
	}
	return "", logs.ContractError{Message: "type not right"}
}

func (c *Contract) Allowance(ctx context.Context, at, buyer string) (*big.Int, error) {
	var res *big.Int
	client, err := ethclient.DialContext(ctx, c.endpoint)
	if err != nil {
		return res, err
	}
	defer client.Close()

	erc20Ins, err := erc.NewERC20(c.erc20, client)
	if err != nil {
		return res, err
	}
	if at == "store" {
		res, err = erc20Ins.Allowance(&bind.CallOpts{From: c.contractAddr}, common.HexToAddress(buyer), c.storeAddr)
		if err != nil {
			return res, err
		}
	}

	if at == "read" {
		res, err = erc20Ins.Allowance(&bind.CallOpts{From: c.contractAddr}, common.HexToAddress(buyer), c.readAddr)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}
