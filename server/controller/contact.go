package controller

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/utils"
)

func (c *Controller) CheckReceipt(ctx context.Context, receipt string) error {
	return c.contract.CheckTrsaction(ctx, receipt)
}

func (c *Controller) SpacePayInfo(ctx context.Context, address string) (IPayPayment, error) {
	out, err := c.contract.Call(ctx, "proxy", "spacePayInfo", common.HexToAddress(address))
	if err != nil {
		return IPayPayment{}, err
	}

	out0 := *abi.ConvertType(out[0], new(IPayPayment)).(*IPayPayment)

	return out0, err
}

func (c *Controller) Allowance(ctx context.Context, at, address string) (*big.Int, error) {
	return c.contract.Allowance(ctx, at, address)
}

func (c *Controller) TrafficPayInfo(ctx context.Context, address string) (IPayPayment, error) {
	out, err := c.contract.Call(ctx, "proxy", "trafficPayInfo", common.HexToAddress(address))
	if err != nil {
		return IPayPayment{}, err
	}

	out0 := *abi.ConvertType(out[0], new(IPayPayment)).(*IPayPayment)

	return out0, err
}

func (c *Controller) CashSpace(ctx context.Context, buyer string) (string, error) {
	check := c.getSpaceCheck(ctx, buyer)
	sender, err := utils.GetSeller(ctx)
	if err != nil {
		return "", err
	}
	return c.contract.CashSpaceCheck(ctx, sender, check.Nonce, check.CheckSize.Uint64(), api.DurationDay, check.Sign)
}

func (c *Controller) CashTraffic(ctx context.Context, buyer string) (string, error) {
	check := c.getTrafficCheck(ctx, buyer)
	sender, err := utils.GetSeller(ctx)
	if err != nil {
		return "", err
	}
	return c.contract.CashTrafficCheck(ctx, sender, check.Nonce, check.CheckSize.Uint64(), check.Sign)
}

func (c *Controller) GetStorePayHash(ctx context.Context, address string, checksize uint64) (string, error) {
	pi, err := c.SpacePayInfo(ctx, address)
	if err != nil {
		return "", err
	}

	return c.contract.GetStorePayHash(ctx, checksize, pi.Nonce), nil
}

func (c *Controller) GetReadPayHash(ctx context.Context, address string, checksize uint64) (string, error) {
	pi, err := c.TrafficPayInfo(ctx, address)
	if err != nil {
		return "", err
	}

	return c.contract.GetReadPayHash(ctx, checksize, pi.Nonce), nil
}

func (c *Controller) BuySpace(ctx context.Context, address string, size uint64) (string, error) {
	return c.contract.BuySpace(ctx, address, size)
}

func (c *Controller) BuyTraffic(ctx context.Context, address string, size uint64) (string, error) {
	return c.contract.BuyTraffic(ctx, address, size)
}

func (c *Controller) Approve(ctx context.Context, at, buyer string, value *big.Int) (string, error) {
	return c.contract.Approve(ctx, at, buyer, value)
}
