package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var globalDID = "did:memo:d687daa192ffa26373395872191e8502cc41fbfbf27dc07d3da3a35de57c2d96"
var globalKey1 = "f9729aef404b8c13d06cf888376b04fd17581b9c308f9b4b16c020736ae89cd4"
var globalKey2 = ""

var token = "520"

func TestAuth(t *testing.T) {
	router := gin.Default()
	LoadAuthRouter(router.Group("/"))

	w := httptest.NewRecorder()
	req, _ := GetLoginRequest(globalDID, globalKey1)
	router.ServeHTTP(w, req)

	t.Log(w.Body.String())

	assert.Equal(t, 200, w.Code)
	// assert.Equal(t, "pong", w.Body.String())

	for index := 1; index < 10; index++ {
		w := httptest.NewRecorder()
		req, _ = GetVerifyIdentityRequest(globalDID, globalKey1, int64(index))
		router.ServeHTTP(w, req)

		t.Log(w.Body.String())

		assert.Equal(t, 200, w.Code)
	}
}

func GetLoginRequest(did string, privateKeyHex string) (*http.Request, error) {
	var timestamp = time.Now().Unix()

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	hash := crypto.Keccak256([]byte(did), []byte(token), int64ToBytes(timestamp))
	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return nil, err
	}

	var payload = make(map[string]interface{})
	payload["did"] = did
	payload["token"] = token
	payload["timestamp"] = timestamp
	payload["signature"] = hexutil.Encode(signature)

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "/login", bytes.NewReader(b))

	return req, err
}

func GetVerifyIdentityRequest(did string, privateKeyHex string, requestID int64) (*http.Request, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	hash := crypto.Keccak256([]byte(did), []byte(token), int64ToBytes(requestID))
	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return nil, err
	}

	var payload = make(map[string]interface{})
	payload["did"] = did
	payload["token"] = token
	payload["requestID"] = requestID
	payload["signature"] = hexutil.Encode(signature)

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", "/identity", bytes.NewReader(b))

	return req, err
}
