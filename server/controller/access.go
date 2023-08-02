package controller

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
)

func (c *Controller) canWrite(ctx context.Context, address string, size uint64, msg api.SignMessage, nonce *big.Int) error {
	err := c.verifySign(ctx, common.HexToAddress(address), msg, nonce)
	if err != nil {
		return err
	}

	return c.checkUpSize(ctx, address, size, msg.Size)
}

func (c *Controller) canRead(ctx context.Context, address string, size uint64, msg api.SignMessage) error {
	pi, err := c.TrafficPayInfo(ctx, address)
	if err != nil {
		return err
	}

	err = c.verifySign(ctx, common.HexToAddress(address), msg, pi.Nonce)
	if err != nil {
		return err
	}

	return c.checkDownSize(ctx, address, size, msg.Size)
}

func (c *Controller) checkUpSize(ctx context.Context, address string, size, checksize uint64) error {
	pi, err := c.SpacePayInfo(ctx, address)
	if err != nil {
		return err
	}

	ci, err := c.getUpCacheInfo(ctx, address)
	if err != nil {
		return err
	}

	if checksize > pi.FreeByte+pi.SizeByte {
		lerr := logs.ControllerError{Message: fmt.Sprintf("space not enough, have %d, need %d", pi.FreeByte+pi.SizeByte, checksize)}
		logger.Error(lerr)
		return lerr
	}

	if ci+size > checksize {
		lerr := logs.ControllerError{Message: fmt.Sprintf("checksize not enough, have %d, need %d", checksize, ci+size)}
		logger.Error(lerr)
		return lerr
	}

	return nil
}

func (c *Controller) checkDownSize(ctx context.Context, address string, size, checksize uint64) error {
	pi, err := c.TrafficPayInfo(ctx, address)
	if err != nil {
		return err
	}

	ci, err := c.getDownCacheInfo(ctx, address)
	if err != nil {
		return err
	}

	if checksize > pi.FreeByte+pi.SizeByte {
		lerr := logs.ControllerError{Message: fmt.Sprintf("space not enough, have %d, need %d", pi.FreeByte+pi.SizeByte, checksize)}
		logger.Error(lerr)
		return lerr
	}

	if ci+size > checksize {
		lerr := logs.ControllerError{Message: fmt.Sprintf("checksize not enough, have %d, need %d", checksize, ci+size)}
		logger.Error(lerr)
		return lerr
	}

	return nil
}

func (c *Controller) getUpCacheInfo(ctx context.Context, address string) (uint64, error) {
	return c.database.GetUpSize(ctx, address)
}

func (c *Controller) getDownCacheInfo(ctx context.Context, address string) (uint64, error) {
	return c.database.GetDownSize(ctx, address)
}

func (c *Controller) verifySign(ctx context.Context, buyer common.Address, msg api.SignMessage, nonce *big.Int) error {
	hashs := c.contract.GetStorePayHash(ctx, msg.Size, nonce)

	hash, err := hexutil.Decode(hashs)
	if err != nil {
		lerr := logs.ControllerError{Message: err.Error()}
		logger.Error(lerr)
		return lerr
	}

	sig := hexutil.MustDecode(msg.Sign)
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
		lerr := logs.ControllerError{Message: "verify sign not success "}
		logger.Error(lerr)
		return lerr
	}

	return nil
}
