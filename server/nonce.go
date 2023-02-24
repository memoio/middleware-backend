package server

import(
	"sync"
	"time"
	"crypto/rand"
	"encoding/hex"
)

type NonceManager struct {
	handledNonce   map[string]int64
	handlingNonce  map[string]int64
	modifyMutex    sync.Mutex

	ExpireEpoch     int64
	ModifyEpoch    int64
	LastModifyTime int64
}

func NewNonceManager(expireEpoch int64, modifyEpoch int64) *NonceManager {
	return &NonceManager{
		handlingNonce: make(map[string]int64), 
		handledNonce: make(map[string]int64), 
		ExpireEpoch: expireEpoch, 
		ModifyEpoch: modifyEpoch, 
		LastModifyTime: time.Now().Unix(), 
	}
}

func (non *NonceManager) GetNonce() string {
    now := time.Now().Unix()
    if now - non.LastModifyTime >= non.ModifyEpoch {
    	non.clearExpiredNonce()
    }

    b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil { 
        return ""
    }

    nonce := hex.EncodeToString(b)
    non.handlingNonce[nonce] = time.Now().Unix() + non.ExpireEpoch

    return nonce
}

func (non *NonceManager) VerifyNonce(nonce string) bool {
	if nonce == "" {
		return false
	}

	now := time.Now().Unix()
	if now - non.LastModifyTime >= non.ModifyEpoch {
		non.clearExpiredNonce()
	}

	expireTime, ok := non.handlingNonce[nonce]
	if ok {
		delete(non.handlingNonce, nonce)
		if now < expireTime {
			return true
		}
	}

	if time.Now().Unix() - non.LastModifyTime < non.ExpireEpoch {
		expireTime, ok = non.handledNonce[nonce]
		if ok {
			delete(non.handledNonce, nonce)
			if now < expireTime {
				return true
			}
		}
	}

	return false
}

func (non *NonceManager) clearExpiredNonce() {
	now := time.Now().Unix()
	non.modifyMutex.Lock()
	defer non.modifyMutex.Unlock()
	if now - non.LastModifyTime >= non.ModifyEpoch {
		non.handledNonce = non.handlingNonce
		non.handlingNonce = make(map[string]int64)
		non.LastModifyTime = now
	}
}