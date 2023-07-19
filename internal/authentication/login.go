package auth

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/memoio/go-did/memo"
	"github.com/memoio/go-did/types"
)

type Request struct {
	Address     string `json:"address"`
	RandomToken string `josn:"token"`
	Timestamp   int64  `json:"timestamp"`
	Signature   string `json:"signature"`
}

func Login(did, token string, timestamp int64, signature string) (bool, error) {
	hash := crypto.Keccak256([]byte(did), []byte(token), int64ToBytes(timestamp))
	sig, err := hexutil.Decode(signature)
	if err != nil {
		return false, err
	}

	publicKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return false, err
	}

	ok, err := CheckAuthPermission(did, crypto.PubkeyToAddress(*publicKey).Hex())
	if err != nil || !ok {
		return ok, err
	}

	err = sessionStore.AddSession(did, token, timestamp)
	if err != nil {
		return false, err
	}

	return true, nil
}

func VerifyIdentity(did, token, payload string, requestID int64, signature string) (bool, error) {
	hash := crypto.Keccak256([]byte(did), []byte(token), int64ToBytes(requestID))
	sig, err := hexutil.Decode(signature)
	if err != nil {
		return false, err
	}
	publicKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return false, err
	}

	ok, err := CheckAuthPermission(did, crypto.PubkeyToAddress(*publicKey).Hex())
	if err != nil || !ok {
		return ok, err
	}

	err = sessionStore.VerifySession(did, token, requestID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func CheckAuthPermission(did, address string) (bool, error) {
	resolver, err := memo.NewMemoDIDResolver("dev")
	if err != nil {
		return false, err
	}

	keys, err := resolver.Dereference(did + "#authentication")
	if err != nil {
		return false, err
	}

	var permission = false
	for _, key := range keys {
		permitAddr, err := types.PublicKeyToAddress(key)
		if err != nil {
			continue
		}

		if address == permitAddr.Hex() {
			permission = true
		}
	}

	return permission, nil
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
