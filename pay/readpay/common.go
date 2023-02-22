package readpay

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/memoio/go-mefs-v2/build"
	"github.com/memoio/go-mefs-v2/lib/types"
	"golang.org/x/xerrors"
)

var (
	opSk         = ""
	contractAddr = common.HexToAddress("")
	opAddr       = common.HexToAddress("")
)

func generateCheck(fromAddr, toAddr common.Address, nonce uint64) (*Check, error) {
	c := &Check{
		ContractAddr: contractAddr,
		OwnerAddr:    opAddr,
		ToAddr:       toAddr,
		Nonce:        nonce,
		Value:        big.NewInt(types.DefaultReadPrice * build.DefaultSegSize * 1024 * 40),
		FromAddr:     fromAddr,
	}

	skECDSA, err := crypto.HexToECDSA(opSk)
	if err != nil {
		return nil, xerrors.Errorf("convert to ECDSA err: %w", err)
	}

	sigByte, err := crypto.Sign(c.Hash(), skECDSA)
	if err != nil {
		return nil, xerrors.Errorf("sign paycheck error: %w", err)
	}
	c.Sig = sigByte

	return c, nil
}
