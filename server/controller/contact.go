package controller

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
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
 
func (c *Controller) TrafficPayInfo(ctx context.Context, address string) (IPayPayment, error) {
	out, err := c.contract.Call(ctx, "proxy", "trafficPayInfo", common.HexToAddress(address))
	if err != nil {
		return IPayPayment{}, err
	}

	out0 := *abi.ConvertType(out[0], new(IPayPayment)).(*IPayPayment)

	return out0, err
}
