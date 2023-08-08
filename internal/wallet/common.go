package wallet

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/memoio/backend/internal/logs"
)

func PrivateToAddr(pi *ecdsa.PrivateKey) (common.Address, error) {
	var addr common.Address

	pk := pi.Public()
	pubKeyECDSA, ok := pk.(*ecdsa.PublicKey)
	if !ok {
		lerr := logs.WalletError{Message: "error casting public key to ECDSA"}
		logger.Error(lerr)
		return addr, lerr
	}

	addr = crypto.PubkeyToAddress(*pubKeyECDSA)
	return addr, nil
}
