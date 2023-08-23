package controller

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/memoio/backend/internal/logs"
)

func (c *Controller) canWrite(ctx context.Context, address string, size uint64, sign string, pi IPayPayment) error {
	checksize, err := c.checkSize(ctx, "up", address, size, pi)
	if err != nil {
		return err
	}

	return c.verifySign(ctx, "store", common.HexToAddress(address), sign, checksize, pi.Nonce)
}

func (c *Controller) canRead(ctx context.Context, address string, size uint64, sign string, pi IPayPayment) error {
	checksize, err := c.checkSize(ctx, "down", address, size, pi)
	if err != nil {
		return err
	}

	return c.verifySign(ctx, "read", common.HexToAddress(address), sign, checksize, pi.Nonce)
}

func (c *Controller) checkSize(ctx context.Context, sizetype, address string, size uint64, pi IPayPayment) (uint64, error) {
	var res, ci uint64
	var err error
	if sizetype == "up" {
		ci, err = c.getUpCacheInfo(ctx, address)
		if err != nil {
			return res, err
		}
	} else {
		ci, err = c.getDownCacheInfo(ctx, address)
		if err != nil {
			return res, err
		}
	}

	checksize := ci + size

	if checksize > pi.FreeByte+pi.SizeByte {
		lerr := logs.ControllerError{Message: fmt.Sprintf("space not enough, have %d, need %d", pi.FreeByte+pi.SizeByte, checksize)}
		logger.Error(lerr)
		return res, lerr
	}

	return checksize, nil
}

func (c *Controller) getUpCacheInfo(ctx context.Context, address string) (uint64, error) {
	return c.database.GetUpSize(ctx, address)
}

func (c *Controller) getDownCacheInfo(ctx context.Context, address string) (uint64, error) {
	return c.database.GetDownSize(ctx, address)
}

func (c *Controller) verifySign(ctx context.Context, at string, buyer common.Address, sign string, checksize uint64, nonce *big.Int) error {
	var hashs string
	if at == "store" {
		hashs = c.contract.GetStorePayHash(ctx, checksize, nonce)
	}
	if at == "read" {
		hashs = c.contract.GetReadPayHash(ctx, checksize, nonce)
	}
	hash, err := hexutil.Decode(hashs)
	if err != nil {
		lerr := logs.ControllerError{Message: err.Error()}
		logger.Error(lerr)
		return lerr
	}

	sig := hexutil.MustDecode(sign)
	if sig[64] == 27 || sig[64] == 28 {
		sig[64] -= 27
	}

	publicKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		lerr := logs.ControllerError{Message: err.Error()}
		logger.Error(lerr)
		return lerr
	}

	addr := crypto.PubkeyToAddress(*publicKey)

	if strings.Compare(buyer.Hex(), addr.Hex()) != 0 {
		lerr := logs.ControllerError{Message: "verify sign not success"}
		logger.Error(lerr)
		return lerr
	}

	return nil
}
