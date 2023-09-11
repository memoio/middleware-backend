package contract

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/contractsv2/go_contracts/erc"
)

func (c *Contract) ApproveTsHash(ctx context.Context, pt api.PayType, sender string, buyValue *big.Int) (api.Transaction, error) {
	var approveAddr common.Address

	switch pt {
	case api.ReadPay:
		approveAddr = c.readAddr
	case api.StorePay:
		approveAddr = c.storeAddr
	default:
		return api.Transaction{}, logs.ContractError{Message: "type not right"}
	}

	return c.GetTrasaction(ctx, c.erc20, sender, "erc20", "approve", approveAddr, buyValue)

}

func (c *Contract) Allowance(ctx context.Context, pt api.PayType, buyer string) (*big.Int, error) {
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

	var allowAddr common.Address

	switch pt {
	case api.StorePay:
		allowAddr = c.storeAddr
	case api.ReadPay:
		allowAddr = c.readAddr
	}

	res, err = erc20Ins.Allowance(&bind.CallOpts{From: c.contractAddr}, common.HexToAddress(buyer), allowAddr)
	if err != nil {
		return res, err
	}

	return res, nil
}
