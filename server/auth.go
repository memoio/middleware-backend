package server

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/memoio/backend/gateway"
	"github.com/shurcooL/graphql"
	"github.com/spruceid/siwe-go"
)

type LoginRequest struct {
	Address   string `json:"address,omitempty"`
	Nonce     string `json:"nonce,omitempty"`
	Domain    string `json:"domain,omitempty"`
	Signature string `json:"signature,omitempty"`
}

type EIP4361Request struct {
	EIP191Message string `json:"message,omitempty"`
	Signature     string `json:"signature,omitempty"`
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
	ErrNullToken      = gateway.AuthenticationFailed{Message: "Token is Null, not found in `Authorization: Bearer ` header"}
	ErrValidToken     = gateway.AuthenticationFailed{Message: "Invalid token"}
	ErrValidTokenType = gateway.AuthenticationFailed{Message: "InValid token type"}

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

func Login(nonceManager *NonceManager, request interface{}) (string, string, error) {
	req, ok := request.(LoginRequest)
	if !ok {
		return "", "", fmt.Errorf("")
	}
	return loginWithEth(nonceManager, req)
}

// func LoginWithMethod(nonceManager *NonceManager, request interface{}, method int, checkRegistered bool) (string, string, error) {
// 	switch method {
// 	case LensMod:
// 		req, ok := request.(EIP4361Request)
// 		if !ok {
// 			return "", "", xerrors.Errorf("")
// 		}
// 		return loginWithLens(req, checkRegistered)
// 	case EthMod:
// 		req, ok := request.(LoginRequest)
// 		if !ok {
// 			return "", "", xerrors.Errorf("")
// 		}
// 		return loginWithEth(nonceManager, req)
// 	}
// 	return "", "", gateway.NotImplemented{Message: ""}
// }

func LoginWithLens(request EIP4361Request, required bool) (string, string, bool, error) {
	message, err := parseLensMessage(request.EIP191Message)
	if err != nil {
		return "", "", false, err
	}

	isRegistered, err := isLensAccount(message.GetAddress().Hex(), required)
	if err != nil {
		return "", "", false, err
	}

	if message.GetDomain() != Domain {
		return "", "", false, gateway.AuthenticationFailed{Message: "Got wrong domain"}
	}

	if message.GetChainID() != 137 {
		return "", "", false, gateway.AuthenticationFailed{Message: "Got wrong chain id"}
	}

	hash := crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(request.EIP191Message), request.EIP191Message)))
	sig, err := hexutil.Decode(request.Signature)
	if err != nil {
		return "", "", false, gateway.AuthenticationFailed{Message: err.Error()}
	}

	sig[len(sig)-1] %= 27
	pubKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return "", "", false, gateway.AuthenticationFailed{Message: err.Error()}
	}

	if message.GetAddress().Hex() != crypto.PubkeyToAddress(*pubKey).Hex() {
		return "", "", false, gateway.AuthenticationFailed{Message: "Got wrong address/signature"}
	}

	accessToken, err := genAccessTokenWithFlag(message.GetAddress().Hex(), isRegistered)
	if err != nil {
		return "", "", false, err
	}

	refreshToken, err := genRefreshTokenWithFlag(message.GetAddress().Hex(), isRegistered)

	return accessToken, refreshToken, isRegistered, err
}

func loginWithEth(nonceManager *NonceManager, request LoginRequest) (string, string, error) {
	var address = request.Address
	var nonce = request.Nonce
	var domain = request.Domain
	var signature = request.Signature

	if address == "" || nonce == "" || domain == "" || signature == "" {
		return "", "", gateway.AuthenticationFailed{Message: "There is an empty parameter"}
	}

	if domain != Domain {
		return "", "", gateway.AuthenticationFailed{Message: "Got wrong domain"}
	}

	if !nonceManager.VerifyNonce(nonce) {
		return "", "", gateway.AuthenticationFailed{Message: "Got wrong nonce"}
	}

	hash := crypto.Keccak256([]byte(address), []byte(nonce), []byte(domain))
	sig, err := hexutil.Decode(signature)
	if err != nil {
		return "", "", gateway.AuthenticationFailed{Message: err.Error()}
	}

	pubKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return "", "", gateway.AuthenticationFailed{Message: err.Error()}
	}

	if address != crypto.PubkeyToAddress(*pubKey).Hex() {
		return "", "", gateway.AuthenticationFailed{Message: "Got wrong address/signature"}
	}

	accessToken, err := genAccessToken(address)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := genRefreshToken(address)

	return accessToken, refreshToken, err
}

func parseLensMessage(message string) (*siwe.Message, error) {
	message = strings.TrimPrefix(message, "\n")
	message = strings.TrimPrefix(message, "https://")
	message = strings.TrimPrefix(message, "http://")
	message = strings.TrimSuffix(message, "\n ")

	return siwe.ParseMessage(message)
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

// Verify token's type, audience, nonce, expires time and signatrue
// Then, return access token, refresh token and usr id
// The format of usr id is did:eth:{usr's publickey key || usr's address}
// func VerifyDidToken(nonceManager *NonceManager, tokenString string) (string, string, string, error) {
// 	if tokenString == "" {
// 		return "", "", "", ErrNullToken
// 	}

// 	claims := &Claims{}
// 	_, _, err := new(jwt.Parser).ParseUnverified(tokenString, claims)
// 	if err != nil {
// 		return "", "", "", ErrValidToken
// 	}

// 	// check Audience
// 	if claims.Audience != Domain {
// 		return "", "", "", ErrValidToken
// 	}

// 	// check token type
// 	if claims.Type != DidToken {
// 		return "", "", "", ErrValidType
// 	}

// 	// check token nonce
// 	if nonceManager.VerifyNonce(claims.Nonce) == false {
// 		return "", "", "", ErrValidToken
// 	}

// 	// check signature, expires time and issued time
// 	token, err := ParseDidToken(tokenString, claims.Subject)
// 	if err != nil || !token.Valid {
// 		return "", "", "", ErrValidToken
// 	}

// 	accessToken, err := GenAccessToken(claims.Subject)
// 	if err != nil {
// 		return "", "", "", err
// 	}

// 	refreshToken, err := GenRefreshToken(claims.Subject)

// 	return accessToken, refreshToken, claims.Subject, err
// }
