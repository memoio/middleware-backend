package auth

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/xerrors"
)

type Request struct {
	Address     string `json:"address"`
	RandomToken string `josn:"token"`
	Timestamp   int64  `json:"timestamp"`
	Signature   string `json:"signature"`
}

func Login(address, token string, chainID int, timestamp int64, signature string) (bool, error) {
	hash := crypto.Keccak256([]byte(address), []byte(token), int64ToBytes(timestamp))
	sig, err := hexutil.Decode(signature)
	if err != nil {
		return false, err
	}

	publicKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return false, err
	}

	if address != crypto.PubkeyToAddress(*publicKey).Hex() {
		return false, xerrors.Errorf("The signature cannot match address")
	}

	err = sessionStore.AddSession(address, token, chainID, timestamp)
	if err != nil {
		return false, err
	}

	return true, nil
}

func VerifyIdentity(token string, chainID int, requestID int64, signature string) (string, error) {
	hash := crypto.Keccak256([]byte(token), int64ToBytes(requestID))
	sig, err := hexutil.Decode(signature)
	if err != nil {
		return "", err
	}
	publicKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return "", err
	}

	address := crypto.PubkeyToAddress(*publicKey).Hex()

	err = sessionStore.VerifySession(address, token, chainID, requestID)
	if err != nil {
		return "", err
	}

	return address, nil
}

func int64ToBytes(v int64) []byte {
	return []byte{
		byte(0xff & v),
		byte(0xff & (v >> 8)),
		byte(0xff & (v >> 16)),
		byte(0xff & (v >> 24)),
		byte(0xff & (v >> 32)),
		byte(0xff & (v >> 40)),
		byte(0xff & (v >> 48)),
		byte(0xff & (v >> 56)),
	}
}
