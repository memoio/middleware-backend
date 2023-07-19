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

func (c *Controller) canWrite(ctx context.Context, address, sign string, size uint64, msg api.SignMessage) error {
	err := c.checkSize(ctx, address, size, msg.Size)
	if err != nil {
		return err
	}

	verifySign(common.HexToAddress(address), sign, msg)
	return nil
}

func canRead() error {
	return nil
}

func (c *Controller) checkSize(ctx context.Context, address string, size, checksize uint64) error {
	pi, err := c.SpacePayInfo(ctx, address)
	if err != nil {
		return err
	}

	ci, err := c.getCacheInfo(ctx, address)
	if err != nil {
		return err
	}

	remain := pi.FreeByte + pi.SizeByte - ci

	if remain < size {
		lerr := logs.ControllerError{Message: fmt.Sprintf("space not enough, have %d, need %d", remain, size)}
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

func (c *Controller) getCacheInfo(ctx context.Context, address string) (uint64, error) {
	return c.database.GetSize(ctx, address)
}

func verifySign(buyer common.Address, sign string, msg api.SignMessage) error {
	hash := com.GetCashCheckHash(msg.StorePayAddr, msg.Seller, msg.Size, msg.Nonce)

	publicKeyHash := buyer.Bytes()[1:]

	publicKey, err := crypto.Ecrecover(nil, publicKeyHash)
	if err != nil {
		lerr := logs.ControllerError{Message: err.Error()}
		logger.Error(lerr)
		return lerr
	}

	signB, err := hex.DecodeString(sign)
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
