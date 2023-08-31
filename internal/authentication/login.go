package auth

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/memoio/go-did/memo"
)

type Request struct {
	Address     string `json:"address"`
	RandomNonce string `josn:"nonce"`
	Timestamp   int64  `json:"timestamp"`
	Signature   string `json:"signature"`
}

func Login(did, nonce string, timestamp int64, signature string) (string, string, error) {
	nonceBytes, err := hexutil.Decode(nonce)
	if err != nil {
		return "", "", err
	}

	sig, err := hexutil.Decode(signature)
	if err != nil {
		return "", "", err
	}

	ok, err := CheckAuthPermission(did, sig, []byte(did), nonceBytes, int64ToBytes(timestamp))
	if err != nil || !ok {
		return "", "", err
	}

	err = sessionStore.AddSession(did, nonce, timestamp)
	if err != nil {
		return "", "", err
	}

	accessToken, err := genAccessToken(did)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := genRefreshToken(did)

	return accessToken, refreshToken, nil
}

func VerifyIdentity(did, nonce, hash string, requestID int64, signature string) (bool, error) {
	err := sessionStore.VerifySession(did, nonce, requestID)
	if err != nil {
		return false, err
	}

	nonceBytes, err := hexutil.Decode(nonce)
	if err != nil {
		return false, err
	}

	var message [][]byte
	if hash != "" {
		hashBytes, err := hexutil.Decode(hash)
		if err != nil {
			return false, err
		}
		message = [][]byte{[]byte(did), nonceBytes, hashBytes, int64ToBytes(requestID)}
	} else {
		message = [][]byte{[]byte(did), nonceBytes, int64ToBytes(requestID)}
	}

	sig, err := hexutil.Decode(signature)
	if err != nil {
		return false, err
	}

	ok, err := CheckAuthPermission(did, sig, message...)
	if err != nil || !ok {
		return ok, err
	}

	return true, nil
}

func CheckAuthPermission(did string, sig []byte, message ...[]byte) (bool, error) {
	resolver, err := memo.NewMemoDIDResolver("dev")
	if err != nil {
		return false, err
	}

	keys, err := resolver.Dereference(did + "#authentication")
	if err != nil {
		return false, err
	}

	for _, key := range keys {
		ok, _ := key.VerifySignature(sig, message...)
		if ok {
			return true, nil
		}
	}

	return false, nil
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
