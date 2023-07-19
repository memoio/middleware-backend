package controller

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
	com "github.com/memoio/contractsv2/common"
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

func (c *Controller) cashCheck(ctx context.Context, address, sign string, msg api.SignMessage) error {
	paymentInfo, err := c.SpacePayInfo(ctx, address)
	if err != nil {
		return err
	}

	hash := com.GetCashCheckHash(msg.StorePayAddr, msg.Seller, msg.Size, paymentInfo.Nonce)
	verify := crypto.VerifySignature(common.HexToAddress(address).Bytes(), hash, []byte(sign))
	if !verify {
		lerr := logs.ContractError{Message: "sign verify failed"}
		logger.Error(lerr)
		return lerr
	}
	return nil
}
