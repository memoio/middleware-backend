package controller

import (
	"context"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
)

func (c *Controller) canWrite(ctx context.Context, address, sign string, size uint64) (api.CheckInfo, error) {
	ci, err := c.getSpaceCheckInfo(ctx, address, size)
	if err != nil {
		return api.CheckInfo{}, err
	}
	if strings.HasPrefix(sign, "0x") {
		sign = "0x" + sign
	}
	sig, err := hexutil.Decode(sign)
	if err != nil {
		return api.CheckInfo{}, logs.ControllerError{Message: err.Error()}
	}
	if sig[64] == 27 || sig[64] == 28 {
		sig[64] -= 27
	}
	ci.Sign = sig
	err = c.verifySign(ctx, "space", ci)
	if err != nil {
		return api.CheckInfo{}, err
	}
	return ci, nil
}

func (c *Controller) canRead(ctx context.Context, address, sign string, size uint64) (api.CheckInfo, error) {
	ci, err := c.getTrafficCheckInfo(ctx, address, size)
	if err != nil {
		return api.CheckInfo{}, err
	}

	sig := hexutil.MustDecode(sign)
	if sig[64] == 27 || sig[64] == 28 {
		sig[64] -= 27
	}
	ci.Sign = sig
	err = c.verifySign(ctx, "traffic", ci)
	if err != nil {
		return api.CheckInfo{}, err
	}
	return ci, nil
}

func (c *Controller) verifySign(ctx context.Context, ct string, ci api.CheckInfo) error {
	var hashs string
	if ct == "space" {
		hashs = c.contract.GetSapceCheckHash(ctx, ci.FileSize.Uint64(), ci.Nonce)
	} else {
		hashs = c.contract.GetTrafficCheckHash(ctx, ci.FileSize.Uint64(), ci.Nonce)
	}
	hash, err := hexutil.Decode(hashs)
	if err != nil {
		lerr := logs.ControllerError{Message: err.Error()}
		logger.Error(lerr)
		return lerr
	}

	publicKey, err := crypto.SigToPub(hash, ci.Sign)
	if err != nil {
		lerr := logs.ControllerError{Message: err.Error()}
		logger.Error(lerr)
		return lerr
	}

	addr := crypto.PubkeyToAddress(*publicKey)

	if strings.Compare(ci.Buyer.Hex(), addr.Hex()) != 0 {
		lerr := logs.ControllerError{Message: "verify sign not success"}
		logger.Error(lerr)
		return lerr
	}

	return nil
}
