package auth

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/siwb"
	"github.com/memoio/backend/internal/siws"
	"github.com/mr-tron/base58"
	"github.com/shurcooL/graphql"
	"github.com/spruceid/siwe-go"
)

var purposeStatement = "The message is only used for login"

type EIP4361Request struct {
	EIP191Message string `json:"message,omitempty"`
	Signature     string `json:"signature,omitempty"`

	// used for registe
	Recommender string `json:"recommender,omitempty"`
	Source      string `json:"source,omitempty"`
	// used for choose user
	UserID int `json:"userID,omitempty"`
}

type BTCSignedMessage struct {
	Message   string `json:"message,omitempty"`
	Signature string `json:"signature,omitempty"`

	// used for registe
	Recommender string `json:"recommender,omitempty"`
	Source      string `json:"source,omitempty"`
	// used for choose user
	UserID int `json:"userID,omitempty"`
}

type SOLSignedMessage struct {
	Message   string `json:"message,omitempty"`
	Signature string `json:"signature,omitempty"`

	// used for registe
	Recommender string `json:"recommender,omitempty"`
	Source      string `json:"source,omitempty"`
	// used for choose user
	UserID int `json:"userID,omitempty"`
}

type profile struct {
	DefaultProfile struct {
		ID   string
		Name string
	} `graphql:"defaultProfile(request: $request)"`
}

type DefaultProfileRequest struct {
	EthereumAddress string `json:"ethereumAddress"`
}

var (
	ErrNullToken      = logs.AuthenticationFailed{Message: "Token is Null, not found in `Authorization: Bearer ` header"}
	ErrValidToken     = logs.AuthenticationFailed{Message: "Invalid token"}
	ErrValidTokenType = logs.AuthenticationFailed{Message: "InValid token type"}

	// ChainID = 985
	Version = 1

	JWTKey  []byte
	Domain  string
	LensAPI string

	DidToken     = 0
	AccessToken  = 1
	RefreshToken = 2

	LensMod = 0x10
	EthMod  = 0x11
)

func InitAuthConfig(jwtKey string, domain string, url string) {
	var err error
	JWTKey, err = hex.DecodeString(jwtKey)
	if err != nil {
		JWTKey = []byte("memo.io")
	}

	Domain = domain
	LensAPI = url
}

func Challenge(domain, address, uri, nonce string, chainID int) (string, error) {
	var opt = map[string]interface{}{
		"chainId":   chainID,
		"statement": purposeStatement,
	}
	msg, err := siwe.InitMessage(domain, address, uri, nonce, opt)
	if err != nil {
		return "", err
	}
	return msg.String(), nil
}

func ChallengeWithBTC(domain, address, uri, nonce string) (string, error) {
	msg, err := siwb.InitMessage(domain, address, uri, nonce, map[string]interface{}{})
	if err != nil {
		return "", err
	}
	return msg.String(), nil
}

func ChallengeWithSOL(domain, address, uri, nonce, chainID string) (string, error) {
	var opt = map[string]interface{}{
		"chainId":   chainID,
		"statement": purposeStatement,
	}
	msg, err := siws.InitMessage(domain, address, uri, nonce, opt)
	if err != nil {
		return "", err
	}
	return msg.String(), nil
}

func Login(nonceManager *NonceManager, request interface{}) (string, string, string, error) {
	req, ok := request.(EIP4361Request)
	if !ok {
		return "", "", "", fmt.Errorf("")
	}
	return loginWithEth(nonceManager, req)
}

func LoginWithLens(request EIP4361Request, required bool) (string, string, string, bool, error) {
	message, err := parseLensMessage(request.EIP191Message)
	if err != nil {
		return "", "", "", false, err
	}

	isRegistered, err := isLensAccount(message.GetAddress().Hex(), required)
	if err != nil {
		return "", "", "", false, err
	}

	if message.GetDomain() != Domain {
		return "", "", "", false, logs.AuthenticationFailed{Message: "Got wrong domain"}
	}

	if message.GetChainID() != 137 {
		return "", "", "", false, logs.AuthenticationFailed{Message: "Got wrong chain id"}
	}

	hash := crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(request.EIP191Message), request.EIP191Message)))
	sig, err := hexutil.Decode(request.Signature)
	if err != nil {
		return "", "", "", false, logs.AuthenticationFailed{Message: err.Error()}
	}

	sig[len(sig)-1] %= 27
	pubKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return "", "", "", false, logs.AuthenticationFailed{Message: err.Error()}
	}

	if message.GetAddress().Hex() != crypto.PubkeyToAddress(*pubKey).Hex() {
		return "", "", "", false, logs.AuthenticationFailed{Message: "Got wrong address/signature"}
	}

	accessToken, err := genAccessTokenWithFlag(message.GetAddress().Hex(), message.GetChainID(), request.UserID, isRegistered)
	if err != nil {
		return "", "", "", false, err
	}

	refreshToken, err := genRefreshTokenWithFlag(message.GetAddress().Hex(), message.GetChainID(), request.UserID, isRegistered)

	return accessToken, refreshToken, message.GetAddress().Hex(), isRegistered, err
}

func loginWithEth(nonceManager *NonceManager, request EIP4361Request) (string, string, string, error) {
	message, err := parseLensMessage(request.EIP191Message)
	if err != nil {
		return "", "", "", err
	}

	// if message.GetChainID() != ChainID {
	// 	return "", "", "", logs.AuthenticationFailed{Message: "Got wrong chain id"}
	// }

	if !nonceManager.VerifyNonce(message.GetNonce()) {
		return "", "", "", logs.AuthenticationFailed{Message: "Got wrong nonce"}
	}

	hash := crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(request.EIP191Message), request.EIP191Message)))
	sig, err := hexutil.Decode(request.Signature)
	if err != nil {
		return "", "", "", logs.AuthenticationFailed{Message: err.Error()}
	}

	sig[len(sig)-1] %= 27
	pubKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return "", "", "", logs.AuthenticationFailed{Message: err.Error()}
	}

	if message.GetAddress().Hex() != crypto.PubkeyToAddress(*pubKey).Hex() {
		return "", "", "", logs.AuthenticationFailed{Message: "Got wrong address/signature"}
	}

	accessToken, err := genAccessToken(message.GetAddress().Hex(), message.GetChainID(), request.UserID)
	if err != nil {
		return "", "", "", err
	}

	refreshToken, err := genRefreshToken(message.GetAddress().Hex(), message.GetChainID(), request.UserID)

	return accessToken, refreshToken, message.GetAddress().Hex(), err
}

func LoginWithBTC(nonceManager *NonceManager, request BTCSignedMessage) (string, string, string, error) {
	message, err := siwb.ParseMessage(request.Message)
	if err != nil {
		return "", "", "", err
	}

	hash, err := message.MessageHash()
	if err != nil {
		return "", "", "", err
	}

	sig, err := base64.StdEncoding.DecodeString(request.Signature)
	if err != nil {
		return "", "", "", err
	}

	pk, _, err := btcec.RecoverCompact(btcec.S256(), sig, hash)
	if err != nil {
		return "", "", "", err
	}

	switch message.GetAddress()[:1] {
	case "1":
		// get P2PKH address
		if message.GetAddress() != getP2PKHAddress(pk.SerializeCompressed()) {
			return "", "", "", logs.AuthenticationFailed{Message: "Got wrong address/signature"}
		}
	case "3":
		// get P2SH address
		if message.GetAddress() != getP2SHAddress(pk.SerializeCompressed()) {
			return "", "", "", logs.AuthenticationFailed{Message: "Got wrong address/signature"}
		}
	default:
		switch message.GetAddress()[:4] {
		case "bc1q":
			// get Native SegWit address
			if message.GetAddress() != getNativeSegWitAddress(pk.SerializeCompressed()) {
				return "", "", "", logs.AuthenticationFailed{Message: "Got wrong address/signature"}
			}
		case "bc1p":
			// TODO: get Traproot address
			return "", "", "", logs.AuthenticationFailed{Message: "Unsupported address format"}
		default:
			return "", "", "", logs.AuthenticationFailed{Message: "Invalid address format"}
		}
	}

	pubKeyEcdsa, err := crypto.UnmarshalPubkey(pk.SerializeUncompressed())
	if err != nil {
		panic(err)
	}

	ethAddress := crypto.PubkeyToAddress(*pubKeyEcdsa)

	accessToken, err := genAccessToken(ethAddress.String(), 0, request.UserID)
	if err != nil {
		return "", "", "", err
	}

	refreshToken, err := genRefreshToken(ethAddress.String(), 0, request.UserID)

	return accessToken, refreshToken, ethAddress.String(), err
}

func LoginWithSOL(nonceManager *NonceManager, request SOLSignedMessage) (string, string, string, error) {
	message, err := siws.ParseMessage(request.Message)
	if err != nil {
		return "", "", "", err
	}

	sig, err := base64.StdEncoding.DecodeString(request.Signature)
	if err != nil {
		return "", "", "", err
	}

	pub, err := base58.Decode(message.GetAddress())
	if err != nil {
		return "", "", "", err
	}

	if !ed25519.Verify(pub, []byte(message.String()), sig) {
		return "", "", "", logs.AuthenticationFailed{Message: "Got wrong address/signature"}
	}

	ethAddress := common.BytesToAddress(crypto.Keccak256(pub[:])[12:])

	accessToken, err := genAccessToken(ethAddress.String(), -1, request.UserID)
	if err != nil {
		return "", "", "", err
	}

	refreshToken, err := genRefreshToken(ethAddress.String(), -1, request.UserID)

	return accessToken, refreshToken, ethAddress.String(), err
}

func parseLensMessage(message string) (*siwe.Message, error) {
	message = strings.TrimPrefix(message, "\n")
	message = strings.TrimPrefix(message, "https://")
	message = strings.TrimPrefix(message, "http://")
	message = strings.TrimSuffix(message, "\n ")

	return siwe.ParseMessage(message)
}

func getP2PKHAddress(publicKey []byte) string {
	address, err := btcutil.NewAddressPubKey(publicKey, &chaincfg.MainNetParams)
	if err != nil {
		return ""
	}
	return address.EncodeAddress()
}

func getP2SHAddress(publicKey []byte) string {
	witnessProg := btcutil.Hash160(publicKey)
	addressWitnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, &chaincfg.MainNetParams)
	if err != nil {
		return ""
	}

	serializedScript, err := txscript.PayToAddrScript(addressWitnessPubKeyHash)
	if err != nil {
		return ""
	}
	addressScriptHash, err := btcutil.NewAddressScriptHash(serializedScript, &chaincfg.MainNetParams)
	if err != nil {
		return ""
	}

	return addressScriptHash.EncodeAddress()
}

func getNativeSegWitAddress(publicKey []byte) string {
	witnessProg := btcutil.Hash160(publicKey)
	addressWitnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, &chaincfg.MainNetParams)
	if err != nil {
		return ""
	}
	return addressWitnessPubKeyHash.EncodeAddress()
}

func isLensAccount(address string, required bool) (bool, error) {
	if required {
		var query profile
		var client = graphql.NewClient(LensAPI, nil)
		var variables = map[string]interface{}{
			"request": DefaultProfileRequest{
				EthereumAddress: address,
			},
		}

		err := client.Query(context.Background(), &query, variables)
		if err != nil {
			return false, err
		}
		if query.DefaultProfile.ID == "" {
			return false, nil
			// return false, gateway.AddressError{Message: "The address{" + address + "} is not registered on lens"}
		}
	}

	return true, nil
}
