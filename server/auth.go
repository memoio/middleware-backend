package server

import(
	"time"
	"context"

	"github.com/shurcooL/graphql"
	"github.com/dgrijalva/jwt-go"
	"github.com/memoio/backend/gateway"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Claims struct {
	Type  int    `json:"type,omitempty"`
	// Nonce string `json:"nonce,omitempty"`
	jwt.StandardClaims
}

type LoginRequest struct {
	Address   string `json:"address,omitempty"`
	Nonce     string `json:"nonce,omitempty"`
	Domain    string `json:"domain,omitempty"`
	Signature string `json:"signature,omitempty"`
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
	ErrNullToken = gateway.AuthenticationFailed{"Token is Null"}
	ErrValidToken = gateway.AuthenticationFailed{"Invalid token"}
	ErrValidType = gateway.AuthenticationFailed{"InValid token type"}
	ErrValidPayload = gateway.AuthenticationFailed{"Invaliad token payload"}

	jwtkey = []byte("memo.io")

	DidToken = 0
	AccessToken = 1
	FreshToken = 2

	LensAccount = 0x10
	EthAccount = 0x11 
)

func Login(nonceManager *NonceManager, request LoginRequest) (string, string, error) {
	return LoginWithMethod(nonceManager, request, LensAccount)
}

func LoginWithMethod(nonceManager *NonceManager, request LoginRequest, method int) (string, string, error) {
	switch method {
	case LensAccount:
		return loginWithLens(nonceManager, request)
	case EthAccount:
		return loginWithEth(nonceManager, request)
	}
	return "", "", gateway.NotImplemented{""}
}

func loginWithLens(nonceManager *NonceManager, request LoginRequest) (string, string, error) {
	if err := isLensAccount(request.Address); err != nil {
		return "", "", err
	}

	return loginWithEth(nonceManager, request)
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

	pubKey, err := crypto.Ecrecover(hash, sig)
    if err != nil {
        return "", "", gateway.AuthenticationFailed{err.Error()}
    }

    if address != common.BytesToAddress(crypto.Keccak256(pubKey[1:])[12:]).Hex() {
    	return "", "", gateway.AuthenticationFailed{"Got wrong address"}
    }

	if !crypto.VerifySignature(pubKey, hash, sig[:len(sig)-1]) {
		return "", "", gateway.AuthenticationFailed{"Got wrong signature"}
	}

	accessToken, err := GenAccessToken(address)
	if err != nil {
		return "", "", err
	}

	freshToken, err := GenFreshToken(address)

	return accessToken, freshToken, err
}

func isLensAccount(address string) error {
	var query profile
	var client = graphql.NewClient("https://api-mumbai.lens.dev", nil)
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
// 		return "", "", "", ErrValidPayload
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

func VerifyAccessToken(tokenString string) (string, error) {
	if tokenString == "" {
		return "", ErrNullToken
	}

	claims := &Claims{}
	_, _, err := new(jwt.Parser).ParseUnverified(tokenString, claims)
	if err != nil {
		return "", ErrValidPayload
	}

	// check Audience
	if claims.Audience != "memo.io" || claims.Issuer != "memo.io" {
		return "", ErrValidToken
	}

	// check token type
	if claims.Type != AccessToken {
		return "", ErrValidType
	}

	// check signature, Expires time and Issued time
	token, err := ParseToken(tokenString)
	if err != nil || !token.Valid {
		return "", ErrValidToken
	}

	return claims.Subject, nil
}

func VerifyFreshToken(tokenString string) (string, error) {
	if tokenString == "" {
		return "", ErrNullToken
	}

	claims := &Claims{}
	_, _, err := new(jwt.Parser).ParseUnverified(tokenString, claims)
	if err != nil {
		return "", ErrValidPayload
	}

	// check Audience
	if claims.Audience != "memo.io" || claims.Issuer != "memo.io" {
		return "", ErrValidToken
	}

	// check token type
	if claims.Type != FreshToken {
		return "", ErrValidType
	}

	token, err := ParseToken(tokenString)
	if err != nil || !token.Valid {
		return "", ErrValidToken
	} 

	return GenAccessToken(claims.Subject)
}

func GenAccessToken(did string) (string, error) {
	expireTime := time.Now().Add(15 * time.Minute)
    claims := &Claims{
        Type: AccessToken,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expireTime.Unix(), 
            IssuedAt:  time.Now().Unix(), 
            Audience:  "memo.io", 
            Issuer:    "memo.io", 
            Subject:   did, 
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtkey)
}

func GenFreshToken(did string) (string, error) {
	expireTime := time.Now().Add(7 * 24 * time.Hour)
    claims := &Claims{
        Type: FreshToken,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expireTime.Unix(), 
            IssuedAt:  time.Now().Unix(), 
            Audience:  "memo.io", 
            Issuer:    "memo.io", 
            Subject:   did, 
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtkey)
}

// func ParseDidToken(tokenString string, did string) (*jwt.Token, error) {
//     return jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
//     	parts := strings.Split(did, ":")
//     	if len(parts) != 3 || parts[0] != "did" || parts[1] != "eth" {
//     		return nil, ErrValidPayload
//     	}

//     	pubKeyBytes, err := hex.DecodeString(parts[2])
//     	if err != nil {
//     		return nil, err
//     	}

//         return crypto.UnmarshalPubkey(pubKeyBytes)
//     })
// }

func ParseToken(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
    	return jwtkey, nil
    })
}