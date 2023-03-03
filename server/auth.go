package server

import(
	"fmt"
	"context"
	"strings"

	"golang.org/x/xerrors"

	"github.com/spruceid/siwe-go"
	"github.com/shurcooL/graphql"
	"github.com/memoio/backend/gateway"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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

var(
	ErrNullToken = gateway.AuthenticationFailed{"Token is Null, not found in `Authorization: Bearer ` header"}
	ErrValidToken = gateway.AuthenticationFailed{"Invalid token"}
	ErrValidTokenType = gateway.AuthenticationFailed{"InValid token type"}

	LensMod = 0x10
	EthMod = 0x11 
)

func Login(nonceManager *NonceManager, request interface{}) (string, string, error) {
	return LoginWithMethod(nonceManager, request, EthMod)
}

func LoginWithMethod(nonceManager *NonceManager, request interface{}, method int) (string, string, error) {
	switch method {
	case LensMod:
		req, ok := request.(EIP4361Request)
		if !ok {
			return "", "", xerrors.Errorf("")
		}
		return loginWithLens(req)
	case EthMod:
		req, ok := request.(LoginRequest)
		if !ok {
			return "", "", xerrors.Errorf("")
		}
		return loginWithEth(nonceManager, req)
	}
	return "", "", gateway.NotImplemented{""}
}

func loginWithLens(request EIP4361Request) (string, string, error) {
	message, err := parseLensMessage(request.EIP191Message)
	if err != nil {
		return "", "", err
	}

	// if err := isLensAccount(message.GetAddress()); err != nil {
	// 	return "", "", err
	// }

	if message.GetDomain() != "memo.io" {
		return "", "", gateway.AuthenticationFailed{"Got wrong domain"}
	}

	if message.GetChainID() != 137 {
		return "", "", gateway.AuthenticationFailed{"Got wrong chain id"}
	}

	hash := crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(request.EIP191Message), request.EIP191Message)))
	sig, err := hexutil.Decode(request.Signature)
	if err != nil {
		return "", "", gateway.AuthenticationFailed{err.Error()}
	}

	sig[len(sig)-1] %= 27
	pubKey, err := crypto.SigToPub(hash, sig)
    if err != nil {
        return "", "", gateway.AuthenticationFailed{err.Error()}
    }

    if message.GetAddress().Hex() != crypto.PubkeyToAddress(*pubKey).Hex() {
    	return "", "", gateway.AuthenticationFailed{"Got wrong address/signature"}
    }

	accessToken, err := genAccessToken(message.GetAddress().Hex())
	if err != nil {
		return "", "", err
	}

	freshToken, err := genFreshToken(message.GetAddress().Hex())

	return accessToken, freshToken, err
}

func loginWithEth(nonceManager *NonceManager, request LoginRequest) (string, string, error) {
	var address = request.Address
	var nonce = request.Nonce
	var domain = request.Domain
	var signature = request.Signature

	if address == "" || nonce == "" || domain == "" || signature == "" {
		return "", "", gateway.AuthenticationFailed{"There is an empty parameter"}
	}

	if domain != "memo.io" {
		return "", "", gateway.AuthenticationFailed{"Got wrong domain"}
	}

	if !nonceManager.VerifyNonce(nonce) {
		return "", "", gateway.AuthenticationFailed{"Got wrong nonce"}
	}

	hash := crypto.Keccak256([]byte(address), []byte(nonce), []byte(domain))
	sig, err := hexutil.Decode(signature)
	if err != nil {
		return "", "", gateway.AuthenticationFailed{err.Error()}
	}

	pubKey, err := crypto.SigToPub(hash, sig)
    if err != nil {
        return "", "", gateway.AuthenticationFailed{err.Error()}
    }

    if address != crypto.PubkeyToAddress(*pubKey).Hex() {
    	return "", "", gateway.AuthenticationFailed{"Got wrong address/signature"}
    }

	accessToken, err := genAccessToken(address)
	if err != nil {
		return "", "", err
	}

	freshToken, err := genFreshToken(address)

	return accessToken, freshToken, err
}

func parseLensMessage(message string) (*siwe.Message, error) {
	message = strings.TrimPrefix(message, "\n")
    message = strings.TrimPrefix(message, "https://")
    message = strings.TrimPrefix(message, "http://")
    message = strings.TrimSuffix(message, "\n ")

    return siwe.ParseMessage(message)
}

func isLensAccount(address string) error {
	var query profile
	var client = graphql.NewClient("https://api.lens.dev", nil)
    var variables = map[string]interface{}{
        "request": DefaultProfileRequest{
            EthereumAddress: address,
        }, 
    }

    err := client.Query(context.Background(), &query, variables)
    if err != nil {
        return err
    }
    if query.DefaultProfile.ID == "" {
    	return gateway.AddressError{"The address{" + address + "} is not registered on lens"}
    }

    return nil
}

// Verify token's type, audience, nonce, expires time and signatrue
// Then, return access token, fresh token and usr id
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
// 	if claims.Audience != "memo.io" {
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

// 	freshToken, err := GenFreshToken(claims.Subject)

// 	return accessToken, freshToken, claims.Subject, err
// }