package controller

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
	com "github.com/memoio/contractsv2/common"
)

func (c *Controller) canWrite(ctx context.Context, address string, size uint64, msg api.SignMessage) error {
	err := verifySign(common.HexToAddress(address), msg)
	if err != nil {
		return err
	}

	return c.checkUpSize(ctx, address, size, msg.Size)
}

func (c *Controller) canRead(ctx context.Context, address string, size uint64, msg api.SignMessage) error {
	err := verifySign(common.HexToAddress(address), msg)
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

func verifySign(buyer common.Address, msg api.SignMessage) error {
	hash := com.GetCashCheckHash(msg.StorePayAddr, msg.Seller, msg.Size, msg.Nonce)

	publicKeyHash := buyer.Bytes()[1:]

	publicKey, err := crypto.Ecrecover(nil, publicKeyHash)
	if err != nil {
		lerr := logs.ControllerError{Message: err.Error()}
		logger.Error(lerr)
		return lerr
	}

	signB, err := hex.DecodeString(msg.Sign)
	if err != nil {
		lerr := logs.ControllerError{Message: err.Error()}
		logger.Error(lerr)
		return lerr
	}

	res := crypto.VerifySignature(publicKey, hash, signB)

	if !res {
		lerr := logs.ControllerError{Message: "verify sign not success "}
		logger.Error(lerr)
		return lerr
	}
	return nil
}
