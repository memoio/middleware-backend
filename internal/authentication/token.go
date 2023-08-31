package auth

import (
	"encoding/hex"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/logs"
	"golang.org/x/xerrors"
)

var (
	ErrNullToken      = logs.AuthenticationFailed{Message: "Token is Null, not found in `Authorization: Bearer ` header"}
	ErrValidToken     = logs.AuthenticationFailed{Message: "Invalid token"}
	ErrValidTokenType = logs.AuthenticationFailed{Message: "InValid token type"}

	JWTKey []byte
	Domain string

	DidToken     = 0
	AccessToken  = 1
	RefreshToken = 2
)

type Claims struct {
	Type int    `json:"type,omitempty"`
	DID  string `josn:"chainid,omitempty"`
	// Nonce string `json:"nonce,omitempty"`
	jwt.StandardClaims
}

func initJWTConfig() {
	config, err := config.ReadFile()
	if err != nil {
		panic(err)
	}

	JWTKey, err = hex.DecodeString(config.SecurityKey)
	if err != nil {
		JWTKey = []byte("memo.io")
	}

	Domain = config.Domain
}

func VerifyAccessToken(tokenString string) (string, error) {
	claims, err := verifyJsonWebToken(tokenString, AccessToken)
	if err != nil {
		return "", err
	}

	return claims.Subject, nil
}

func VerifyRefreshToken(tokenString string) (string, error) {
	claims, err := verifyJsonWebToken(tokenString, RefreshToken)
	if err != nil {
		return "", err
	}

	return genAccessToken(claims.Subject)
}

func genAccessToken(did string) (string, error) {
	return genJsonWebToken(did, AccessToken)
}

func genRefreshToken(subject string) (string, error) {
	return genJsonWebToken(subject, RefreshToken)
}

func verifyJsonWebToken(tokenString string, jwtType int) (*Claims, error) {
	parts := strings.SplitN(tokenString, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, ErrNullToken
	}

	claims := &Claims{}
	_, _, err := new(jwt.Parser).ParseUnverified(parts[1], claims)
	if err != nil {
		return nil, ErrValidToken
	}

	// check Audience
	if claims.Audience != Domain || claims.Issuer != Domain {
		return nil, ErrValidToken
	}

	// check token type
	if claims.Type != jwtType {
		return nil, ErrValidTokenType
	}

	// check signature, Expires time and Issued time
	token, err := parseToken(parts[1])
	if err != nil || !token.Valid {
		return nil, ErrValidToken
	}

	return claims, nil
}

func genJsonWebToken(did string, jwtType int) (string, error) {
	var expireTime int64
	if jwtType == AccessToken {
		expireTime = time.Now().Add(2 * time.Hour).Unix()
	} else if jwtType == RefreshToken {
		expireTime = time.Now().Add(7 * 24 * time.Hour).Unix()
	} else {
		return "", xerrors.Errorf("unsupported json web token type")
	}

	claims := &Claims{
		Type: jwtType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime,
			IssuedAt:  time.Now().Unix(),
			Audience:  Domain,
			Issuer:    Domain,
			Subject:   did,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTKey)
}

// func ParseDidToken(tokenString string, did string) (*jwt.Token, error) {
//     return jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
//     	parts := strings.Split(did, ":")
//     	if len(parts) != 3 || parts[0] != "did" || parts[1] != "eth" {
//     		return nil, ErrValidToken
//     	}

//     	pubKeyBytes, err := hex.DecodeString(parts[2])
//     	if err != nil {
//     		return nil, err
//     	}

//         return crypto.UnmarshalPubkey(pubKeyBytes)
//     })
// }

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return JWTKey, nil
	})
}
